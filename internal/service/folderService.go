package service

import (
	"slices"

	"github.com/Catizard/lampghost_wails/internal/dto"
	"github.com/Catizard/lampghost_wails/internal/entity"
	"github.com/Catizard/lampghost_wails/internal/vo"
	"github.com/rotisserie/eris"
	"gorm.io/gorm"

	. "github.com/samber/lo"
)

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
func (s *FolderService) AddFolder(param *vo.FolderVo) (*entity.Folder, error) {
	if param == nil {
		return nil, eris.New("AddFolder: param cannot be nil")
	}
	if param.FolderName == "" {
		return nil, eris.New("AddFolder: FolderName cannot be empty")
	}
	if param.CustomTableID == 0 {
		return nil, eris.New("AddFolder: CustomTableID cannot be 0")
	}
	var ret *entity.Folder
	if err := s.db.Transaction(func(tx *gorm.DB) error {
		if count, err := selectFolderCount(tx, &vo.FolderVo{
			FolderName:    param.FolderName,
			CustomTableID: param.CustomTableID,
		}); err != nil {
			return eris.Wrap(err, "query folder")
		} else if count > 0 {
			return eris.Errorf("AddFolder: folder name %s is duplicated", param.FolderName)
		}
		ret = param.Entity()
		if err := tx.Create(ret).Error; err != nil {
			return eris.Wrap(err, "create folder")
		}
		return nil
	}); err != nil {
		return nil, eris.Wrap(err, "transaction")
	}

	return ret, nil
}

func (s *FolderService) BindSongToFolder(param *vo.BindToFolderVo) error {
	if param == nil {
		return eris.New("BindSongToFolder: param cannot be nil")
	}
	if len(param.FolderIDs) == 0 {
		return eris.New("BindSongToFolder: bind song to nothing")
	}
	if param.Md5 == "" && param.Sha256 == "" {
		return eris.New("BindSongToFolder: md5 and sha256 are both not provided")
	}

	cache, err := queryDefaultSongHashCache(s.db)
	if err != nil {
		return eris.Wrap(err, "query cache")
	}
	if param.Md5 == "" {
		md5, ok := cache.GetMD5(param.Sha256)
		if !ok {
			return eris.New("BindSongToFolder: cannot bind an unexist song to folder")
		}
		param.Md5 = md5
	}
	if param.Sha256 == "" {
		sha256, ok := cache.GetSHA256(param.Md5)
		if !ok {
			return eris.New("BindSongToFolder: cannot bind an unexist song to folder")
		}
		param.Sha256 = sha256
	}

	content := entity.FolderContent{
		Sha256: param.Sha256,
		Md5:    param.Md5,
		Title:  param.Title,
	}

	if err := bindSongToFolder(s.db, content, param.FolderIDs); err != nil {
		return eris.Wrap(err, "bind")
	}
	return nil
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

// Similar function to QueryDiffTableDataWithRival, query one folder contents
// with player related fields (e.g lamp, play count) and difficult table tags
//
// Also the difficult table tags are built in memory not sql, since
// one content could be related to multiple difficult tables, which is
// obviously conflict with pagination
//
// Requirements:
//  1. FolderID & RivalID should not be empty
func (s *FolderService) QueryFolderContentWithRival(filter *vo.FolderContentVo) ([]*dto.FolderContentDto, int, error) {
	if filter.FolderID <= 0 {
		return nil, 0, eris.New("QueryFolderContentWithRival: FolderID should > 0")
	}
	if filter.RivalID <= 0 {
		return nil, 0, eris.New("QueryFolderContentWithRival: RivalID should > 0")
	}

	contents, _, err := findFolderContentListWithRival(s.db, filter)
	if err != nil {
		return nil, 0, eris.Wrap(err, "query folder_content")
	}
	tableTags, _, err := queryDiffTableTag(s.db, &vo.DiffTableDataVo{
		Md5s: Map(contents, func(content *dto.FolderContentDto, _ int) string {
			return content.Md5
		}),
	})
	if err != nil {
		return nil, 0, err
	}
	ForEach(contents, func(content *dto.FolderContentDto, _ int) {
		content.TableTags = make([]*dto.DiffTableTagDto, 0)
		ForEach(tableTags, func(tag *dto.DiffTableTagDto, _ int) {
			if tag.Md5 == content.Md5 {
				content.TableTags = append(content.TableTags, tag)
			}
		})
	})
	return contents, len(contents), nil
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
	ret := make([]*dto.FolderDto, len(rawFolders))
	for i := range rawFolders {
		ret[i] = dto.NewFolderDto(rawFolders[i], nil)
	}
	return ret, len(ret), nil
}

func (s *FolderService) FindFolderContentList(filter *vo.FolderContentVo) ([]*entity.FolderContent, int, error) {
	return findFolderContentList(s.db, filter)
}

func findFolderList(tx *gorm.DB, filter *vo.FolderVo) ([]*entity.Folder, int, error) {
	var folders []*entity.Folder
	if err := tx.Debug().Scopes(scopeFolderFilter(filter)).Find(&folders).Error; err != nil {
		return nil, 0, err
	}
	return folders, len(folders), nil
}

func selectFolderCount(tx *gorm.DB, filter *vo.FolderVo) (int64, error) {
	var count int64
	if err := tx.Model(&entity.Folder{}).Scopes(scopeFolderFilter(filter)).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
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

func scopeFolderFilter(filter *vo.FolderVo) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if filter == nil {
			return db
		}
		moved := db.Where(filter.Entity())
		// Extra filters here
		moved = moved.Scopes(scopeInIDs(filter.IDs))
		if filter.IgnoreSha256 != nil {
			moved = moved.Where("folder.id not in (select distinct folder_id from folder_content fc where fc.sha256 = ?)", *filter.IgnoreSha256)
		}
		return moved
	}
}
