package service_test

import (
	"testing"
	"time"

	"github.com/Catizard/lampghost_wails/internal/config"
	"github.com/Catizard/lampghost_wails/internal/database"
	"github.com/Catizard/lampghost_wails/internal/service"
	"github.com/charmbracelet/log"
)

// NOTE: This task requires config has been setted and wouldn't clean
// the downloaded files automatically. So I left this flag to skip the tests
// here entirely
const SKIP_DOWNLOAD_FLAG = true

func untilTaskFinished(service *service.DownloadTaskService, ch chan<- int) {
	// Ensure every tasks has been submitted, should be a very quick step
	time.Sleep(2 * time.Second)
	for {
		_, wait, running := service.InternalTaskCount()
		if wait == 0 && running == 0 {
			ch <- 1
			return
		}
		time.Sleep(1 * time.Second)
	}
}

func TestSubmitSingleMD5DownloadTask(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	t.Run("SmokeTest", func(t *testing.T) {
		if SKIP_DOWNLOAD_FLAG {
			t.SkipNow()
		}
		db, err := database.NewMemoryDatabase()
		if err != nil {
			t.Fatalf("db: %s", err)
		}
		config, err := config.ReadConfig()
		if err != nil {
			t.Fatalf("config: %s", err)
		}
		downloadTaskService := service.NewDownloadTaskService(db, config, nil)

		md5s := []string{
			"552467f149f79e72e783f863eebef7b3",
			"3fab6a6423490fa4ef43460368dcbaba",
			"432aab3b3f4f74f1cb226b8f41577686",
			"5987f98a8e2940bc63ae58537c60d963",
			"1b84baafb2bf4af86f926aaff45067ad",
		}
		for _, md5 := range md5s {
			go func() {
				if err := downloadTaskService.SubmitSingleMD5DownloadTask(md5, nil); err != nil {
					log.Fatalf("cannot submit: %s", err)
				}
			}()
		}

		ch := make(chan int)
		go untilTaskFinished(downloadTaskService, ch)
		<-ch
	})
}
