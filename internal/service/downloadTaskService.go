package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/Catizard/lampghost_wails/internal/config"
	"github.com/Catizard/lampghost_wails/internal/config/download"
	"github.com/Catizard/lampghost_wails/internal/dto"
	"github.com/Catizard/lampghost_wails/internal/entity"
	"github.com/charmbracelet/log"
	"github.com/imroc/req/v3"
	"github.com/rotisserie/eris"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"gorm.io/gorm"
)

type DownloadTaskService struct {
	db                    *gorm.DB
	ctx                   context.Context
	mutex                 sync.Mutex
	tasks                 []*entity.DownloadTask
	waitTasks             []*entity.DownloadTask
	runningTasks          map[uint]*entity.DownloadTask
	updMsgReceiver        chan taskUpdMsg
	subscribeConfigChange <-chan any
	config                *config.ApplicationConfig
	// Exeprimental, only for test case
	taskID uint
}

type taskUpdMsg struct {
	taskID        uint
	final         bool
	err           error
	downloadSize  int64
	contentLength int64
}

func NewDownloadTaskService(db *gorm.DB, config *config.ApplicationConfig, configNotify <-chan any) *DownloadTaskService {
	service := &DownloadTaskService{
		db:                    db,
		tasks:                 make([]*entity.DownloadTask, 0),
		waitTasks:             make([]*entity.DownloadTask, 0),
		runningTasks:          make(map[uint]*entity.DownloadTask),
		taskID:                1,
		updMsgReceiver:        make(chan taskUpdMsg),
		config:                config,
		subscribeConfigChange: configNotify,
	}
	go service.receive()
	// go service.debugProgress()
	go service.pushupState()
	go service.listenUpdateConfig()
	return service
}

func (s *DownloadTaskService) InjectContext(ctx context.Context) {
	s.ctx = ctx
}

func (s *DownloadTaskService) lock() {
	s.mutex.Lock()
}

func (s *DownloadTaskService) unlock() {
	s.mutex.Unlock()
}

// For internal test, return: task count, wait task count, running task count
func (s *DownloadTaskService) InternalTaskCount() (int, int, int) {
	s.lock()
	defer s.unlock()
	return len(s.tasks), len(s.waitTasks), len(s.runningTasks)
}

// DownloadTaskService's life cycle, receive other routine's message
// and update internal states
func (s *DownloadTaskService) receive() {
	for {
		select {
		case msg := <-s.updMsgReceiver:
			s.handleUpdateTask(&msg)
		default:
		}
		s.tryKickingWaitTask()
		time.Sleep(150 * time.Millisecond)
	}
}

// Push current download task states up to frontend
func (s *DownloadTaskService) pushupState() {
	for {
		tasks, n, err := s.FindDownloadTaskList()
		if err != nil {
			log.Errorf("cannot pushup download task state: %s", err)
		} else if n > 0 {
			runtime.EventsEmit(s.ctx, "DownloadTask:pushup", tasks)
		}
		time.Sleep(1 * time.Second)
	}
}

// Only for debug usage
func (s *DownloadTaskService) debugProgress() {
	for {
		s.lock()
		for _, task := range s.runningTasks {
			log.Debugf("task: %d(%s), download=%d, content=%d", task.ID, task.URL, task.DownloadSize, task.ContentLength)
		}
		s.unlock()
		time.Sleep(1 * time.Second)
	}
}

// Update the config of DownloadTaskService
//
// NOTE: This function should be called only when there is no running task,
// if not, it's not lampghost's job to handle tasks correctly
func (s *DownloadTaskService) listenUpdateConfig() {
	for {
		<-s.subscribeConfigChange
		go func() {
			s.lock()
			log.Debugf("[DownloadTaskService] updating config")
			if config, err := config.ReadConfig(); err != nil {
				// TODO: emit a message to frontend here
				log.Errorf("cannot read config: %s", err)
			} else {
				s.config = config
			}
			s.unlock()
		}()
	}
}

func (s *DownloadTaskService) handleUpdateTask(msg *taskUpdMsg) error {
	s.lock()
	defer s.unlock()
	if task, ok := s.runningTasks[msg.taskID]; ok {
		if msg.final {
			delete(s.runningTasks, msg.taskID)
			if msg.err != nil {
				log.Errorf("[DownloadTaskService] task %d fails: %s", msg.taskID, msg.err)
				*task.Status = entity.TASK_ERROR
				task.DownloadSize = 0
				task.ContentLength = 0
				task.ErrorMessage = msg.err.Error()
			} else {
				*task.Status = entity.TASK_SUCCESS
			}
		} else {
			*task.Status = entity.TASK_DOWNLOAD
			task.DownloadSize = msg.downloadSize
			task.ContentLength = msg.contentLength
		}
	} else {
		log.Warnf("[DownloadTaskService] discard updte msg: %v", *msg)
	}
	return nil
}

func (s *DownloadTaskService) submitTaskError(taskID uint, err error) {
	s.updMsgReceiver <- taskUpdMsg{
		taskID: taskID,
		final:  true,
		err:    err,
	}
}

func (s *DownloadTaskService) tryKickingWaitTask() {
	s.lock()
	defer s.unlock()
	if s.config.MaximumDownloadCount == len(s.runningTasks) || len(s.waitTasks) == 0 {
		return
	}
	next := s.waitTasks[0]
	taskID := next.ID
	log.Debugf("[DownloadTaskService] try kicking task %d(%s)", taskID, next.URL)
	s.waitTasks = s.waitTasks[1:]
	s.runningTasks[taskID] = next
	go func() {
		// Open a no timeout, cancelable client
		client := req.C().SetTimeout(0)
		// Prevent a very rare race condition?
		s.lock()
		ctx, cancel := context.WithCancel(context.Background())
		s.runningTasks[taskID].Cancel = cancel
		contextLockedReq := client.R().SetContext(ctx)
		s.unlock()
		resp, err := contextLockedReq.
			SetOutputFile(next.IntermediateFilePath).
			SetDownloadCallbackWithInterval(func(info req.DownloadInfo) {
				s.updMsgReceiver <- taskUpdMsg{
					taskID:        next.ID,
					final:         false,
					err:           nil,
					downloadSize:  info.DownloadedSize,
					contentLength: info.Response.ContentLength,
				}
			}, 1*time.Second).
			Get(next.URL)
		if err != nil {
			s.submitTaskError(next.ID, err)
			return
		}
		if !resp.IsSuccessState() {
			// NOTE: At this moment, the content should be placed at file
			// Lampghost will try read the error response data, if anything went wrong,
			// just cancel it
			f, err := os.Open(next.IntermediateFilePath)
			if err != nil {
				s.submitTaskError(next.ID, eris.New("remote server returns an unexpected error"))
				return
			}
			body, err := io.ReadAll(f)
			if err != nil {
				s.submitTaskError(next.ID, eris.New("remote server returns an unexpected error"))
				return
			}
			s.submitTaskError(next.ID, errors.New(string(body)))
			return
		}
		filename := ""
		if next.TaskName != nil {
			filename = *next.TaskName
		}
		if next.FallbackName != "" {
			filename = next.FallbackName
		}
		contentDisposition := resp.GetHeader("Content-Disposition")
		log.Debugf("Content-Disposition: %s", contentDisposition)
		if contentDisposition != "" {
			if _, params, err := mime.ParseMediaType(contentDisposition); err == nil {
				filename = params["filename"]
			} else {
				log.Warn("[DownloadTaskService] cannot parse media type from Content-Disposition")
			}
		} else {
			log.Warn("[DownloadTaskService] cannot fetch Content-Disposition from response")
		}
		// NOTE: <del>Below check & conversion was stolen from wriggle, sorry wriggle!</del>
		if filename == "" || filename == "/" || filename == "." {
			s.submitTaskError(next.ID, eris.New("cannot determine filename"))
		}

		filename = filepath.Clean(filename)
		filename = strings.ReplaceAll(filename, "/", "_")
		filename = strings.ReplaceAll(filename, "\\", "_")
		filename = strings.ReplaceAll(filename, ":", "_")
		filename = strings.ReplaceAll(filename, "*", "_")
		filename = strings.ReplaceAll(filename, "?", "_")
		filename = strings.ReplaceAll(filename, "\"", "_")
		filename = strings.ReplaceAll(filename, "<", "_")
		filename = strings.ReplaceAll(filename, ">", "_")
		filename = strings.ReplaceAll(filename, "|", "_")

		targetPath := filepath.Join(s.config.DownloadDirectory, filename)

		if err := os.MkdirAll(filepath.Dir(targetPath), 0o755); err != nil {
			s.submitTaskError(next.ID, eris.Wrapf(err, "cannot create directory for %s", targetPath))
			return
		}
		if _, err := os.Stat(targetPath); err == nil {
			log.Warn("[DownloadTaskService] target file is already existed, would be replaced with the current one")
		} else if !os.IsNotExist(err) {
			s.submitTaskError(next.ID, eris.Errorf("unexpected stat(%s) error: %s", targetPath, err))
			return
		}
		if err := os.Rename(next.IntermediateFilePath, targetPath); err != nil {
			s.submitTaskError(next.ID, eris.Wrap(err, "cannot rename"))
			return
		}

		// Everything is done
		s.updMsgReceiver <- taskUpdMsg{
			taskID: next.ID,
			final:  true,
			err:    nil,
		}
		runtime.EventsEmit(s.ctx, "global:notify", dto.NotificationDto{
			Type:    "success",
			Content: fmt.Sprintf("%s download successfully", filename),
		})
	}()
}

// Submit a new download task
func (s *DownloadTaskService) SubmitDownloadTask(url string, taskName *string) error {
	currentTaskID := s.taskID
	currentTaskID++
	intermediateFileName := fmt.Sprintf("%d.crdownload", currentTaskID)
	return s.submitSingleDownloadTask(currentTaskID, url, intermediateFileName, "", taskName)
}

func (s *DownloadTaskService) SubmitSingleMD5DownloadTask(md5 string, taskName *string) error {
	if err := s.config.EnableDownload(); err != nil {
		return eris.Wrap(err, "download config is not complete")
	}
	if md5 == "" {
		return eris.New("assert: md5 cannot be empty")
	}
	downloadSource := download.GetDownloadSource(s.config.DownloadSite)
	url, fallbackName, err := downloadSource.GetDownloadURLFromMD5(md5)
	if err != nil {
		return err
	}
	s.lock()
	for _, task := range s.tasks {
		if task.URL == url {
			// Skip duplicate download url
			s.unlock()
			return nil
		}
	}
	s.unlock()
	log.Debugf("[DownloadTaskService] build url: %s", url)
	currentTaskID := s.taskID
	s.taskID++
	intermediateFileName := fmt.Sprintf("%d.crdownload", currentTaskID)
	return s.submitSingleDownloadTask(currentTaskID, url, intermediateFileName, fallbackName, taskName)
}

func (s *DownloadTaskService) submitSingleDownloadTask(id uint, url, intermediateFileName, fallbackName string, taskName *string) error {
	if err := s.config.EnableDownload(); err != nil {
		return err
	}
	intermediateFilePath := filepath.Join(s.config.DownloadDirectory, intermediateFileName)
	s.lock()
	defer s.unlock()
	status := entity.TASK_PREPARE
	task := entity.DownloadTask{
		Model: gorm.Model{
			ID: id,
		},
		URL:                  url,
		Status:               &status,
		IntermediateFilePath: intermediateFilePath,
		FallbackName:         fallbackName,
		TaskName:             taskName,
		DownloadSize:         0,
		ContentLength:        0,
	}
	s.tasks = append(s.tasks, &task)
	s.waitTasks = append(s.waitTasks, &task)
	return nil
}

// Query a current snapshot of download tasks
func (s *DownloadTaskService) FindDownloadTaskList() ([]*entity.DownloadTask, int, error) {
	s.lock()
	defer s.unlock()
	ret := make([]*entity.DownloadTask, len(s.tasks))
	copy(ret, s.tasks)
	slices.Reverse(ret)
	return ret, len(ret), nil
}

func (s *DownloadTaskService) CancelDownloadTask(taskID uint) error {
	s.lock()
	defer s.unlock()
	if task, ok := s.runningTasks[taskID]; ok {
		if task.Cancel != nil {
			*task.Status = entity.TASK_CANCEL
			task.Cancel()
			delete(s.runningTasks, taskID)
			// NOTE: Now, the canceld task is not in wait queue nor running queue
			// It's only referenced in all task list, requring a 'Restart' command
			// to rejoin the party
		}
	}
	return nil
}

func (s *DownloadTaskService) RestartDownloadTask(taskID uint) error {
	s.lock()
	defer s.unlock()
	// Won't be a problem for now
	for _, task := range s.tasks {
		if task.ID == taskID {
			if *task.Status != entity.TASK_CANCEL && *task.Status != entity.TASK_ERROR {
				return eris.New("assert: cannot restart a task is not canceled or failed")
			}
			*task.Status = entity.TASK_PREPARE
			s.waitTasks = append(s.waitTasks, task)
			break
		}
	}
	return nil
}
