package service

import (
	"fmt"
	"slices"

	"github.com/Catizard/lampghost_wails/internal/dto"
	"github.com/Catizard/lampghost_wails/internal/entity"
	"github.com/Catizard/lampghost_wails/internal/vo"
	"github.com/charmbracelet/log"
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
//
// Requirements:
// 1) `Sha256` & `Title` & `folderIDs` must be existed.
// 2) if len(folderIDs) == 0, nothing would be done
//
// This function implements `incremental update`, which means the `folderIDs` doesn't define the final
// bindings for the `content` but only the updates
func bindSongToFolder(tx *gorm.DB, content entity.FolderContent, folderIDs []uint) error {
	if len(folderIDs) == 0 {
		return nil
	}
	// Fetch previous bindings
	previous, _, err := findFolderContentList(tx, &vo.FolderContentVo{
		FolderIDs: folderIDs,
		Sha256:    content.Sha256,
	})
	if err != nil {
		return err
	}
	// If we already binds it, do nothing
	newFolderIDs := make([]uint, 0)
	for _, folderID := range folderIDs {
		if !slices.ContainsFunc(previous, func(p *entity.FolderContent) bool {
			return p.FolderID == folderID
		}) {
			newFolderIDs = append(newFolderIDs, folderID)
		}
	}
	if len(newFolderIDs) == 0 {
		return nil // Okay dokey
	}

	folders, _, err := findFolderList(tx, &vo.FolderVo{
		IDs: newFolderIDs,
	})
	if err != nil {
		return err
	}
	folderIDMapsToSelf := make(map[uint]*entity.Folder)
	for _, folder := range folders {
		folderIDMapsToSelf[folder.ID] = folder
	}

	// Insert new bindings
	candidate := make([]*entity.FolderContent, 0)
	for _, folderID := range newFolderIDs {
		newContent := content
		newContent.FolderID = folderID
		newContent.FolderName = folderIDMapsToSelf[folderID].FolderName
		candidate = append(candidate, &newContent)
	}
	return tx.Create(&candidate).Error
}

func (s *FolderService) FindFolderTree(filter *vo.FolderVo) ([]*dto.FolderDto, int, error) {
	return findFolderTree(s.db, filter)
}

func (s *FolderService) FindFolderList(filter *vo.FolderVo) ([]*dto.FolderDto, int, error) {
	rawFolders, _, err := findFolderList(s.db, filter)
	if err != nil {
		return nil, 0, err
	}
	queryBoundSha256 := filter.IgnoreSha256
	boundFolderIDs := make(map[uint]interface{})
	if filter != nil && filter.IgnoreRivalSongDataID != nil {
		songData, err := findRivalSongDataByID(s.db, *filter.IgnoreRivalSongDataID)
		if err != nil {
			return nil, 0, err
		}
		queryBoundSha256 = &songData.Sha256
	}
	if queryBoundSha256 != nil && *queryBoundSha256 != "" {
		contents, _, err := findFolderContentList(s.db, &vo.FolderContentVo{Sha256: *queryBoundSha256})
		if err != nil {
			return nil, 0, err
		}
		for _, content := range contents {
			boundFolderIDs[content.FolderID] = new(interface{})
		}
	}
	log.Debugf("boundFolderIDs: %v", boundFolderIDs)
	ret := make([]*dto.FolderDto, 0)
	for i := range rawFolders {
		if _, ok := boundFolderIDs[rawFolders[i].ID]; ok {
			continue
		}
		ret = append(ret, dto.NewFolderDto(rawFolders[i], nil))
	}
	return ret, len(ret), nil
}

func (s *FolderService) FindFolderContentList(filter *vo.FolderContentVo) ([]*entity.FolderContent, int, error) {
	return findFolderContentList(s.db, filter)
}

func (s *FolderService) GenerateFolderDefinition() ([]dto.FolderDefinitionDto, int, error) {
	return generateFolderDefinition(s.db)
}

// Generate folder definition for each folder
func generateFolderDefinition(tx *gorm.DB) ([]dto.FolderDefinitionDto, int, error) {
	rawFolders, _, err := findFolderList(tx, nil)
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

func findFolderTree(tx *gorm.DB, filter *vo.FolderVo) ([]*dto.FolderDto, int, error) {
	rawFolders, _, err := findFolderList(tx, filter)
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

	rawContents, _, err := findFolderContentList(tx, &vo.FolderContentVo{FolderIDs: folderIDs})
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

	// Hack: If no specific rival, skip for merging score log
	scoreLogSha256Map := make(map[string][]*dto.RivalScoreLogDto)
	if filter != nil && filter.RivalID != 0 {
		sha256s := make([]string, 0)
		for _, rawContent := range rawContents {
			sha256s = append(sha256s, rawContent.Sha256)
		}
		scoreLogSha256Map, err = findRivalMaximumClearScoreLogSha256Map(tx, &vo.RivalScoreLogVo{
			RivalId: filter.RivalID,
			Sha256s: sha256s,
		})
		if err != nil {
			return nil, 0, err
		}
	}

	folders := make([]*dto.FolderDto, len(rawFolders))
	for i, rawFolder := range rawFolders {
		contents := make([]*dto.FolderContentDto, 0)
		if rawContents, ok := folderIDMapsToContents[rawFolder.ID]; ok {
			for _, rawContent := range rawContents {
				content := dto.NewFolderContentDto(rawContent)
				if _, ok := scoreLogSha256Map[content.Sha256]; ok {
					content.Lamp = int(scoreLogSha256Map[content.Sha256][0].Clear)
				} else {
					content.Lamp = 0
				}
				contents = append(contents, content)
			}
		}
		folders[i] = dto.NewFolderDto(rawFolder, contents)
	}
	return folders, len(folders), nil
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
