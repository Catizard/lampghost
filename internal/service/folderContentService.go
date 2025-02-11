package service

import (
	"github.com/Catizard/lampghost_wails/internal/entity"
	"github.com/Catizard/lampghost_wails/internal/vo"
	"gorm.io/gorm"
)

func findFolderContentList(tx *gorm.DB, filter *vo.FolderContentVo) ([]*entity.FolderContent, int, error) {
	if filter == nil {
		var contents []*entity.FolderContent
		if err := tx.Find(&contents).Error; err != nil {
			return nil, 0, err
		}
		return contents, len(contents), nil
	}

	var contents []*entity.FolderContent
	if err := tx.Debug().Where(filter.Entity()).Scopes(
		scopeInIDs(filter.IDs),
		scopeInFolderIDs(filter.FolderIDs),
	).Find(&contents).Error; err != nil {
		return nil, 0, err
	}
	return contents, len(contents), nil
}
