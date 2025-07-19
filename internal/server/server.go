package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Catizard/lampghost_wails/internal/dto"
	"github.com/Catizard/lampghost_wails/internal/entity"
	"github.com/Catizard/lampghost_wails/internal/service"
	"github.com/Catizard/lampghost_wails/internal/vo"
	"github.com/charmbracelet/log"
	. "github.com/samber/lo"
	"gorm.io/gorm"
)

// TODO: make port configurable
// However make it configurable might be broken in the future, since a ir connect jar
// cannot set port dynamically
const port = 7391

type InternalServer struct {
	customDiffTableService *service.CustomDiffTableService
	customCourseService    *service.CustomCourseService
	folderService          *service.FolderService
	rivalInfoService       *service.RivalInfoService
	rivalScoreLogService   *service.RivalScoreLogService
	rivalTagService        *service.RivalTagService
	rivalSongDataService   *service.RivalSongDataService
}

func NewInternalServer(
	customDiffTableService *service.CustomDiffTableService,
	customCourseService *service.CustomCourseService,
	folderService *service.FolderService,
	rivalInfoService *service.RivalInfoService,
	rivalScoreLogService *service.RivalScoreLogService,
	rivalTagService *service.RivalTagService,
	rivalSongDataService *service.RivalSongDataService,
) *InternalServer {
	return &InternalServer{
		customDiffTableService: customDiffTableService,
		customCourseService:    customCourseService,
		folderService:          folderService,
		rivalInfoService:       rivalInfoService,
		rivalScoreLogService:   rivalScoreLogService,
		rivalTagService:        rivalTagService,
		rivalSongDataService:   rivalSongDataService,
	}
}

// Setup an internal server for mocking network resource (difficult table or IR)
//
// Mock a difficult table distribution server:
//  1. /table/[???].json: return a difficult table metadata which name is '???' and could be imported by beatoraja
//  2. /content/[???].json: return a difficult table's contents data which custom_table_id = '???'
//
// Mock an IR server for providing version lockable rivals import mechanism
// 1. /ir/...: route dispatch, see 'irHandler' for details
func (s *InternalServer) RunServer() error {
	http.HandleFunc("/table/", s.tableHandler)
	http.HandleFunc("/content/", s.tableContentHandler)
	http.HandleFunc("/ir/", s.irHandler)
	go http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	return nil
}

func (s *InternalServer) tableHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	suffix := strings.TrimPrefix(path, "/table/")
	log.Debugf("suffix: %s", suffix)
	if !strings.HasSuffix(suffix, ".json") {
		http.Error(w, "assert: table must be suffixed with .json", 500)
		return
	}
	tableName := suffix[:strings.Index(suffix, ".json")]
	log.Debugf("table name: %s", tableName)
	customTables, _, err := s.customDiffTableService.FindCustomDiffTableList(&vo.CustomDiffTableVo{
		Name: tableName,
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("faile to export table[%s]: %v", tableName, err), 500)
		return
	}
	if len(customTables) != 1 {
		http.Error(w, fmt.Sprintf("cannot determinate table[%s]: %v", tableName, err), 500)
		return
	}
	customTable := customTables[0]
	folderList, _, err := s.folderService.FindFolderList(&vo.FolderVo{
		CustomTableID: customTable.ID,
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("faile to export table[%s]: %v", tableName, err), 500)
		return
	}
	export := &dto.DiffTableHeaderExportDto{
		Name:      customTable.Name,
		Symbol:    customTable.Symbol,
		HeaderUrl: fmt.Sprintf("http://localhost:%d/table/%s.json", port, tableName),
		DataUrl:   fmt.Sprintf("http://localhost:%d/content/%d.json", port, customTable.ID),
		LevelOrder: Map(folderList, func(folder *dto.FolderDto, _ int) string {
			return folder.FolderName
		}),
		Courses: make([][]dto.DiffTableCourseExportDto, 0),
	}
	courseList, n, err := s.customCourseService.FindCustomCourseList(&vo.CustomCourseVo{
		CustomTableID: customTable.ID,
	})
	courseDataList, _, err := s.customCourseService.FindCustomCourseDataList(&vo.CustomCourseDataVo{
		CustomCourseIDs: Map(courseList, func(course *entity.CustomCourse, _ int) uint {
			return course.ID
		}),
	})

	if n > 0 {
		export.Courses = append(export.Courses, make([]dto.DiffTableCourseExportDto, 0))
		for _, rawCourse := range courseList {
			// I'm too lazy to convert it into a map...whatever
			md5s := make([]string, 0)
			for _, courseData := range courseDataList {
				if courseData.CustomCourseID != rawCourse.ID {
					continue
				}
				md5s = append(md5s, courseData.Md5)
			}
			if len(md5s) == 0 {
				continue
			}
			export.Courses[0] = append(export.Courses[0], dto.DiffTableCourseExportDto{
				Name:       rawCourse.Name,
				Constraint: strings.Split(rawCourse.Constraints, ","),
				Md5:        md5s,
			})
		}
	}
	data, err := json.Marshal(export)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to marshal: %s", err), 500)
		return
	}
	fmt.Fprintf(w, "%s", data)
}

func (s *InternalServer) tableContentHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	suffix := strings.TrimPrefix(path, "/content/")
	if !strings.HasSuffix(suffix, ".json") {
		http.Error(w, "assert: table must be suffixed with .json", 500)
		return
	}
	tableIDStr := suffix[:strings.Index(suffix, ".json")]
	tableID, err := strconv.Atoi(tableIDStr)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to export table[%d]: %s", tableID, err), 500)
		return
	}
	folders, _, err := s.folderService.FindFolderTree(&vo.FolderVo{
		CustomTableID: uint(tableID),
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("query folder tree: %v", err), 500)
		return
	}
	contents := make([]*dto.DiffTableDataExportDto, 0)
	for _, folder := range folders {
		for _, folderContent := range folder.Contents {
			def := &dto.DiffTableDataExportDto{
				Level:   folder.FolderName,
				Sha256:  folderContent.Sha256,
				Md5:     folderContent.Md5,
				Title:   folderContent.Title,
				Comment: folderContent.Comment,
			}
			contents = append(contents, def)
		}
	}
	data, err := json.Marshal(contents)
	if err != nil {
		log.Errorf("failed to marshal json: %v", err)
		http.Error(w, "Failed to marshal json", 500)
		return
	}
	fmt.Fprintf(w, "%s", data)
}

// Dispatch IR requests.
//  1. ir/rivals: return registered rivals back
//  2. ir/scores/[rival id]: return one rival's score data back
func (s *InternalServer) irHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	route := strings.TrimPrefix(path, "/ir/")
	log.Debugf("[InternalServer | IR] route: %s", route)
	if strings.HasPrefix(route, "rivals") {
		s.rivalsHandler(w, r)
	} else if strings.HasPrefix(route, "scores") {
		s.scoresHandler(w, r)
	} else {
		http.Error(w, "Not Found", 404)
		return
	}
}

func (s *InternalServer) rivalsHandler(w http.ResponseWriter, r *http.Request) {
	rivals, n, err := s.rivalInfoService.FindRivalInfoList(&vo.RivalInfoVo{
		IgnoreMainUser: true,
		ReverseImport:  1,
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to query rival: %s", err), 500)
		return
	}
	if n == 0 {
		fmt.Fprintf(w, "[]")
	}
	rivalTagIDs := make([]uint, 0)
	for _, rival := range rivals {
		if rival.LockTagID != 0 {
			rivalTagIDs = append(rivalTagIDs, rival.LockTagID)
		}
	}
	lockTags, _, err := s.rivalTagService.FindRivalTagList(&vo.RivalTagVo{IDs: rivalTagIDs})
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to query rival tags: %s", err), 500)
		return
	}

	ret := Map(rivals, func(rival *dto.RivalInfoDto, _ int) *entity.IRPlayer {
		playerData := &entity.IRPlayer{
			ID:   rival.ID,
			Name: rival.Name,
			Rank: fmt.Sprintf("%d", rival.LockTagID),
		}
		for _, lockTag := range lockTags {
			if lockTag.ID == rival.LockTagID && lockTag.Symbol != "" {
				playerData.Name = fmt.Sprintf("%s(%s)", rival.Name, lockTag.Symbol)
			}
		}
		return playerData
	})
	data, err := json.Marshal(ret)
	if err != nil {
		http.Error(w, fmt.Sprintf("marshal: %s", err), 500)
		return
	}
	fmt.Fprintf(w, "%s", data)
}

func (s *InternalServer) scoresHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	strRivalID := strings.TrimPrefix(path, "/ir/scores/")
	log.Debugf("[InternalServer | IR] /ir/scores/: %s", strRivalID)
	rivalID, err := strconv.Atoi(strRivalID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse rival's id: %s", err), 500)
		return
	}

	rivals, n, err := s.rivalInfoService.FindRivalInfoList(&vo.RivalInfoVo{
		IgnoreMainUser: true,
		ReverseImport:  1,
		Model: gorm.Model{
			ID: uint(rivalID),
		},
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to query rival: %s", err), 500)
		return
	}
	if n == 0 {
		http.Error(w, "No Data", 404)
		return
	}

	rival := rivals[0]

	tagTime := time.Time{}
	if rival.LockTagID != 0 {
		tag, err := s.rivalTagService.FindRivalTagByID(rival.LockTagID)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to query tag: %s", err), 500)
			return
		}
		tagTime = tag.RecordTime
	}

	lampData, _, err := s.rivalScoreLogService.QueryReverseImportScoreData(&vo.RivalScoreLogVo{
		RivalId:       rival.ID,
		EndRecordTime: tagTime,
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to query score log: %s", err), 500)
		return
	}
	data, err := json.Marshal(lampData)
	if err != nil {
		http.Error(w, fmt.Sprintf("marshal: %s", err), 500)
		return
	}
	fmt.Fprintf(w, "%s", data)
}
