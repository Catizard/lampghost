package service

import (
	"fmt"
	"slices"

	"github.com/Catizard/lampghost_wails/internal/dto"
	"github.com/Catizard/lampghost_wails/internal/entity"
	"github.com/Catizard/lampghost_wails/internal/vo"
	"gorm.io/gorm"
)

const MAX_FOLDER_COUNT = 25
const BEGIN_FOLDER_INDEX = 5

type FolderService struct {
	db *gorm.DB
}

func NewFolderService(db *gorm.DB) *FolderService {
	return &FolderService{
		db: db,
	}
}

// Add a new folder
func (s *FolderService) AddFolder(folderName string) (*entity.Folder, error) {
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := checkDuplicateFolderName(tx, folderName); err != nil {
		tx.Rollback()
		return nil, err
	}

	prevFolders, prevFoldersLen, err := findFolderList(tx, nil)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	if prevFoldersLen > MAX_FOLDER_COUNT {
		tx.Rollback()
		return nil, fmt.Errorf("cannot add new folder: folder count exceeds to %d", MAX_FOLDER_COUNT)
	}

	nextBitIndex := findFolderBitIndex(prevFolders)
	newFolder := entity.NewFolder(folderName, nextBitIndex)
	if err := tx.Create(&newFolder).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	return &newFolder, nil
}

// Delete a folder, and its contents
func (s *FolderService) DelFolder(ID uint) error {
	if err := s.db.Transaction(func(tx *gorm.DB) error {
		var candidate entity.Folder
		if err := tx.First(&candidate, ID).Error; err != nil {
			return err
		}
		if err := tx.Unscoped().Where("folder_id = ?", candidate.ID).Delete(&entity.FolderContent{}).Error; err != nil {
			return err
		}
		if err := tx.Delete(&entity.Folder{}, candidate.ID).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

// Delete specific folder content
func (s *FolderService) DelFolderContent(contentID uint) error {
	if err := s.db.Unscoped().Where("id = ?", contentID).Delete(&entity.FolderContent{}).Error; err != nil {
		return err
	}
	return nil
}

// Bind one song to multiple folders
// Requirements:
// 1) `songDataID` & `folderIDs` must be existed.
// 2) `folderIDs` defines the final binding, for example, if the song "AIR" is currently binds to folder "1" & "2"
// After calling BindSongToFolder("AIR", ["2", "3"]), "AIR" is now binds to "2" & "3"
// 3) This function must keep old data. In previous example, the content which represents "AIR" and "2" won't be modified
func (s *FolderService) BindSongToFolder(songDataID uint, folderIDs []uint) error {
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Check existence
	songData, err := findRivalSongDataByID(tx, songDataID)
	if err != nil {
		tx.Rollback()
		return err
	}

	rawFolders, rlen, err := findFolderList(tx, &vo.FolderVo{IDs: folderIDs})
	if err != nil {
		return err
	}
	if rlen != len(folderIDs) {
		return fmt.Errorf("bind song to folder failed: folder id array is cracked")
	}

	folderIDMapsToSelf := make(map[uint]*entity.Folder)
	for _, folder := range rawFolders {
		folderIDMapsToSelf[folder.ID] = folder
	}

	// Fetch previous bindings
	previous, _, err := findFolderContentList(tx, &vo.FolderContentVo{FolderIDs: folderIDs})
	if err != nil {
		return err
	}

	// Delete unused bindings
	unused := make([]uint, 0)
	for _, prev := range previous {
		if slices.Contains(folderIDs, prev.ID) {
			unused = append(unused, prev.ID)
		}
	}
	if err := tx.Unscoped().Where(scopeInIDs(unused)).Delete(&entity.FolderContent{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	// Insert new bindings
	candidate := make([]*entity.FolderContent, 0)
	for _, folderID := range folderIDs {
		if slices.ContainsFunc(previous, func(prev *entity.FolderContent) bool {
			return prev.ID == folderID
		}) {
			candidate = append(candidate, entity.FromSongDataToFolderContent(folderIDMapsToSelf[folderID], songData))
		}
	}
	if err := tx.Create(&candidate).Error; err != nil {
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}
	return nil
}

func (s *FolderService) FindFolderTree() ([]dto.FolderDto, int, error) {
	rawFolders, _, err := findFolderList(s.db, nil)
	if err != nil {
		return nil, 0, err
	}
	if len(rawFolders) == 0 {
		return nil, 0, err
	}

	folderIDs := make([]uint, len(rawFolders))
	for i, folder := range rawFolders {
		folderIDs[i] = folder.ID
	}

	rawContents, _, err := findFolderContentList(s.db, &vo.FolderContentVo{FolderIDs: folderIDs})
	if err != nil {
		return nil, 0, err
	}
	folderIDMapsToContents := make(map[uint][]*entity.FolderContent)
	for _, content := range rawContents {
		if _, ok := folderIDMapsToContents[content.FolderID]; !ok {
			folderIDMapsToContents[content.FolderID] = make([]*entity.FolderContent, 0)
		}
		folderIDMapsToContents[content.FolderID] = append(folderIDMapsToContents[content.FolderID], content)
	}

	folders := make([]dto.FolderDto, len(rawFolders))
	for i, rawFolder := range rawFolders {
		contents := make([]dto.FolderContentDto, 0)
		if rawContents, ok := folderIDMapsToContents[rawFolder.ID]; ok {
			for _, rawContent := range rawContents {
				contents = append(contents, *dto.NewFolderContentDto(rawContent))
			}
		}
		folders[i] = *dto.NewFolderDto(rawFolder, contents)
	}
	return folders, len(folders), nil
}

func (s *FolderService) FindFolderList() ([]*entity.Folder, int, error) {
	return findFolderList(s.db, nil)
}

// Generate folder definition json with all folders
func (s *FolderService) GenerateJson() ([]dto.FolderDefinitionDto, int, error) {
	rawFolders, _, err := findFolderList(s.db, nil)
	if err != nil {
		return nil, 0, err
	}
	folderDefinitions := make([]dto.FolderDefinitionDto, len(rawFolders))
	for i, rawFolder := range rawFolders {
		folderDefinitions[i] = *dto.NewFolderDefinitionDto(rawFolder)
	}
	return folderDefinitions, len(folderDefinitions), nil
}

func checkDuplicateFolderName(tx *gorm.DB, folderName string) error {
	var dupCount int64
	if err := tx.Model(&entity.Folder{}).Where(&entity.Folder{FolderName: folderName}).Count(&dupCount).Error; err != nil {
		return err
	}
	if dupCount > 0 {
		return fmt.Errorf("folder name: %s is duplicated", folderName)
	}
	return nil
}

func findFolderList(tx *gorm.DB, filter *vo.FolderVo) ([]*entity.Folder, int, error) {
	if filter == nil {
		var folders []*entity.Folder
		if err := tx.Find(&folders).Error; err != nil {
			return nil, 0, err
		}
		return folders, len(folders), nil
	}

	rawFilter := filter.Entity()
	var folders []*entity.Folder
	basicFiltered := tx.Where(rawFilter)
	if len(filter.IDs) > 0 {
		basicFiltered.Scopes(scopeInIDs(filter.IDs))
	}
	if err := basicFiltered.Find(&folders).Error; err != nil {
		return nil, 0, err
	}
	return folders, len(folders), nil
}

// Query if there exists a folder that satisfies the condition
func queryFolderExistence(tx *gorm.DB, filter *vo.FolderVo) (bool, error) {
	if filter == nil {
		var dupCount int64
		if err := tx.Model(&entity.Folder{}).Count(&dupCount).Error; err != nil {
			return false, err
		}
		return dupCount > 0, nil
	}

	rawFilter := filter.Entity()
	var dupCount int64
	basicFiltered := tx.Where(rawFilter)
	if len(filter.IDs) > 0 {
		basicFiltered.Scopes(scopeInIDs(filter.IDs))
	}
	if err := basicFiltered.Count(&dupCount).Error; err != nil {
		return false, nil
	}
	return dupCount > 0, nil
}

func findFolderBitIndex(folders []*entity.Folder) int {
	mex := BEGIN_FOLDER_INDEX
	for {
		noProgess := true
		for _, folder := range folders {
			if folder.BitIndex == mex {
				mex++
				noProgess = false
				break
			}
		}
		if noProgess {
			break
		}
	}
	return mex
}
