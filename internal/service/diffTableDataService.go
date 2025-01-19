package service

import (
	"github.com/Catizard/lampghost_wails/internal/entity"
	"gorm.io/gorm"
)

func findDiffTableDataByID(tx *gorm.DB, ID uint) (*entity.DiffTableData, error) {
	var data *entity.DiffTableData
	if err := tx.Find(&data, ID).Error; err != nil {
		return nil, err
	}
	return data, nil
}
