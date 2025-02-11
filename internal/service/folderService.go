package service

import (
	"fmt"
	"slices"

	"github.com/Catizard/lampghost_wails/internal/dto"
	"github.com/Catizard/lampghost_wails/internal/entity"
	"github.com/Catizard/lampghost_wails/internal/vo"
	"github.com/charmbracelet/log"
	"gorm.io/driver/sqlite"
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

// Sync current folder definition to main user's songdata.db file
//
//	WARNING:
//	This function would modify songdata.db directly and it's unrecoverable
func (s *FolderService) SyncSongData() error {
	// 1) connect to songdata.db
	mainUser, err := queryMainUser(s.db)
	if err != nil {
		return err
	}
	songDataDB, err := gorm.Open(sqlite.Open(*mainUser.SongDataPath))
	if err != nil {
		return err
	}
	songDataService := NewSongDataService(songDataDB)
	// 2) generate contents definition
	contentDefinition, _, err := generateContentDefinition(s.db)
	if err != nil {
		return err
	}
	// 3) refersh data
	return songDataService.SyncFolderContentDefinition(contentDefinition)
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
// 1) `diffTableDataID` & `folderIDs` must be existed.
// 2) `folderIDs` defines the final binding, for example, if the song "AIR" is currently binds to folder "1" & "2"
// After calling BindSongToFolder("AIR", ["2", "3"]), "AIR" is now binds to "2" & "3"
// 3) This function must keep old data. In previous example, the content which represents "AIR" and "2" won't be modified
func (s *FolderService) BindSongToFolder(diffTableDataID uint, folderIDs []uint) error {
	log.Debugf("[FolderService] Calling BindSongToFolder with diffTableDataID=%v, folderIDs=%v", diffTableDataID, folderIDs)
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Check existence
	songData, err := findDiffTableDataByID(tx, diffTableDataID)
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
	if len(unused) > 0 {
		if err := tx.Unscoped().Where("id in ?", unused).Delete(&entity.FolderContent{}).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	// Insert new bindings
	candidate := make([]*entity.FolderContent, 0)
	for _, folderID := range folderIDs {
		if !slices.ContainsFunc(previous, func(prev *entity.FolderContent) bool {
			log.Debugf("prev.ID=%v, folderID=%v, equals?=%v", prev.ID, folderID, prev.ID == folderID)
			return prev.ID == folderID
		}) {
			candidate = append(candidate, entity.FromDiffTableDataToFolderContent(folderIDMapsToSelf[folderID], songData))
		}
	}
	if len(candidate) > 0 {
		if err := tx.Create(&candidate).Error; err != nil {
			return err
		}
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

func (s *FolderService) FindFolderList(filter *vo.FolderVo) ([]*entity.Folder, int, error) {
	return findFolderList(s.db, filter)
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

// Generate content definition for each contents within folder
func generateContentDefinition(tx *gorm.DB) ([]dto.FolderContentDefinitionDto, int, error) {
	rawFolders, _, err := findFolderList(tx, nil)
	if err != nil {
		return nil, 0, err
	}
	folderIDMapsToBit := make(map[uint]int)
	for _, rawFolder := range rawFolders {
		folderIDMapsToBit[rawFolder.ID] = rawFolder.BitIndex
	}
	rawContents, _, err := findFolderContentList(tx, nil)
	if err != nil {
		return nil, 0, err
	}
	sha256MapsToMask := make(map[string]int)
	for _, rawContent := range rawContents {
		if _, ok := sha256MapsToMask[rawContent.Sha256]; !ok {
			sha256MapsToMask[rawContent.Sha256] = 0
		}
		sha256MapsToMask[rawContent.Sha256] |= 1 << folderIDMapsToBit[rawContent.FolderID]
	}

	definition := make([]dto.FolderContentDefinitionDto, 0)
	for sha256, mask := range sha256MapsToMask {
		definition = append(definition, dto.FolderContentDefinitionDto{
			Sha256: sha256,
			Mask:   mask,
		})
	}
	return definition, len(definition), nil
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
