package service

import (
	"fmt"

	"github.com/Catizard/lampghost_wails/internal/dto"
	"github.com/Catizard/lampghost_wails/internal/entity"
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

	prevFolders, prevFoldersLen, err := findFolderList(tx)
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

func (s *FolderService) FindFolderTree() ([]dto.FolderDto, int, error) {
	rawFolders, _, err := findFolderList(s.db)
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

	rawContents, _, err := findFolderContentByFolderIDs(s.db, folderIDs)
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

// Generate folder definition json with all folders
func (s *FolderService) GenerateJson() ([]dto.FolderDefinitionDto, int, error) {
	rawFolders, _, err := findFolderList(s.db)
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
	if err := tx.Model(&entity.Folder{}).Count(&dupCount).Error; err != nil {
		return err
	}
	if dupCount > 0 {
		return fmt.Errorf("folder name: %s is duplicated", folderName)
	}
	return nil
}

func findFolderList(tx *gorm.DB) ([]*entity.Folder, int, error) {
	var folders []*entity.Folder
	if err := tx.Find(&folders).Error; err != nil {
		return nil, 0, err
	}
	return folders, len(folders), nil
}

// Find specific folder contents by folder id
func findFolderContentByFolderID(tx *gorm.DB, folderID uint) ([]*entity.FolderContent, int, error) {
	var contents []*entity.FolderContent
	if err := tx.Where("folder_id = ?", folderID).Find(&contents).Error; err != nil {
		return nil, 0, err
	}
	return contents, len(contents), nil
}

// Extends to findFolderContentByFolderID, which could query multiple folders
func findFolderContentByFolderIDs(tx *gorm.DB, folderIDs []uint) ([]*entity.FolderContent, int, error) {
	var contents []*entity.FolderContent
	if err := tx.Where("folder_id in ?", folderIDs).Find(&contents).Error; err != nil {
		return nil, 0, err
	}
	return contents, len(contents), nil
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
