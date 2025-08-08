package controller

import (
	"github.com/Catizard/lampghost_wails/internal/entity"
	"github.com/Catizard/lampghost_wails/internal/result"
	"github.com/Catizard/lampghost_wails/internal/service"
	"github.com/charmbracelet/log"
	"github.com/rotisserie/eris"
)

type SongDirectoryController struct {
	songDirectoryService *service.SongDirectoryService
}

func NewSongDirectoryController(songDirectoryService *service.SongDirectoryService) *SongDirectoryController {
	return &SongDirectoryController{
		songDirectoryService: songDirectoryService,
	}
}

func (ctl *SongDirectoryController) SaveSongDirectories(songDirecotires []string) result.RtnMessage {
	log.Info("[Controller] calling SongDirectoryController.SaveSongDirectories")
	if err := ctl.songDirectoryService.SaveSongDirectories(songDirecotires); err != nil {
		log.Errorf("[SongDirectoryController] returning err: %s", eris.ToString(err, true))
		return result.NewErrorMessage(err)
	}
	return result.SUCCESS
}

func (ctl *SongDirectoryController) FindSongDirectories() result.RtnDataList {
	log.Info("[Controller] calling SongDirectoryController.FindSongDirectories")
	rows, _, err := ctl.songDirectoryService.FindSongDirectories()
	if err != nil {
		log.Errorf("[SongDirectoryController] returning err: %s", eris.ToString(err, true))
		return result.NewErrorDataList(err)
	}
	return result.NewRtnDataList(rows)
}

func (ctl *SongDirectoryController) GENERATOR_SONG_DIRECTORY() *entity.SongDirectory { return nil }
