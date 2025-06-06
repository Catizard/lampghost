package service

import (
	"github.com/Catizard/lampghost_wails/internal/dto"
	"github.com/Catizard/lampghost_wails/internal/entity"
	"github.com/Catizard/lampghost_wails/internal/vo"
	"github.com/rotisserie/eris"
	"gorm.io/gorm"
)

type CustomDiffTableService struct {
	db *gorm.DB
}

func NewCustomDiffTableService(db *gorm.DB) *CustomDiffTableService {
	return &CustomDiffTableService{
		db: db,
	}
}

func (s *CustomDiffTableService) AddCustomDiffTable(param *vo.CustomDiffTableVo) error {
	if param == nil {
		return eris.New("AddCustomDiffTable: param cannot be nil")
	}
	if param.Name == "" {
		return eris.New("AddCustomDiffTable: name cannot be empty")
	}
	return s.db.Transaction(func(tx *gorm.DB) error {
		count, err := selectCustomDiffTableCount(tx, &vo.CustomDiffTableVo{
			Name: param.Name,
		})
		if err != nil {
			return eris.Wrap(err, "query custom_diff_table")
		}
		if count != 0 {
			return eris.New("AddCustomDiffTable: name %s is duplicated")
		}
		return addCustomDiffTable(tx, param.Entity())
	})
}

func (s *CustomDiffTableService) DeleteCustomDiffTable(ID uint) error {
	if ID == 0 {
		return eris.New("DeleteCustomDiffTable: ID cannot be 0")
	}
	if ID == 1 {
		return eris.New("DeleteCustomDiffTable: cannot delete default table")
	}

	if err := s.db.Transaction(func(tx *gorm.DB) error {
		// CustomDiffTable
		if err := tx.Delete(&entity.CustomDiffTable{}, ID).Error; err != nil {
			return eris.Wrap(err, "delete custom_diff_table")
		}
		// Currently, deleting a custom difficult table wouldn't delete
		// all related data, this is for the possible data recorver and
		// storing these data wouldn't be a problem
		return nil
	}); err != nil {
		return eris.Wrap(err, "transaction")
	}
	return nil
}

// Update one custom difficult table
//
// NOTE: param must contains all fields or unmentioned fields would be overwritten by zero value
func (s *CustomDiffTableService) UpdateCustomDiffTable(param *vo.CustomDiffTableVo) error {
	if param == nil {
		return eris.New("UpdateCustomDiffTable: param cannot be nil")
	}
	if param.ID <= 0 {
		return eris.New("UpdateCustomDiffTable: ID cannot be 0")
	}

	if err := s.db.Transaction(func(tx *gorm.DB) error {
		return tx.Debug().Select("*").Where("ID = ?", param.ID).Updates(param.Entity()).Error
	}); err != nil {
		return eris.Wrap(err, "transaction")
	}
	return nil
}

func (s *CustomDiffTableService) FindCustomDiffTableList(filter *vo.CustomDiffTableVo) ([]*dto.CustomDiffTableDto, int, error) {
	var out []*dto.CustomDiffTableDto
	var err error
	if err := s.db.Transaction(func(tx *gorm.DB) error {
		out, _, err = findCustomDiffTableList(tx, filter)
		if err != nil {
			return eris.Wrap(err, "query custom_diff_table")
		}
		return nil
	}); err != nil {
		return make([]*dto.CustomDiffTableDto, 0), 0, eris.Wrap(err, "transaction")
	}
	return out, len(out), nil
}

func (s *CustomDiffTableService) FindCustomDiffTableByID(ID uint) (*dto.CustomDiffTableDto, error) {
	raw := entity.CustomDiffTable{}
	if err := s.db.First(&raw, ID).Error; err != nil {
		return nil, eris.Wrap(err, "query custom_diff_table")
	}
	return dto.NewCustomDiffTableDto(&raw), nil
}

func addCustomDiffTable(tx *gorm.DB, param *entity.CustomDiffTable) error {
	if err := tx.Create(param).Error; err != nil {
		return eris.Wrap(err, "create custom_diff_table")
	}
	return nil
}

func selectCustomDiffTableCount(tx *gorm.DB, filter *vo.CustomDiffTableVo) (int64, error) {
	moved := tx.Model(&entity.CustomDiffTable{})
	var count int64
	if err := moved.Scopes(scopeCustomDiffTableFilter(filter)).Count(&count).Error; err != nil {
		return 0, eris.Wrap(err, "query custom_diff_table")
	}
	return count, nil
}

func findCustomDiffTableList(tx *gorm.DB, filter *vo.CustomDiffTableVo) ([]*dto.CustomDiffTableDto, int, error) {
	moved := tx.Model(&entity.CustomDiffTable{})
	if filter != nil {
		moved = moved.Scopes(
			scopeCustomDiffTableFilter(filter),
			pagination(filter.Pagination),
		)
	}
	var out []*dto.CustomDiffTableDto
	if err := moved.Find(&out).Error; err != nil {
		return nil, 0, eris.Wrap(err, "query custom_diff_table")
	}
	// pagination
	if filter != nil && filter.Pagination != nil {
		count, err := selectCustomDiffTableCount(tx, filter)
		if err != nil {
			return nil, 0, eris.Wrap(err, "query custom_diff_table")
		}
		filter.Pagination.PageCount = calcPageCount(count, filter.Pagination.PageSize)
	}
	return out, len(out), nil
}

// Common query scope for vo.CustomDiffTableVo
func scopeCustomDiffTableFilter(filter *vo.CustomDiffTableVo) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if filter == nil {
			return db
		}
		moved := db.Where(filter.Entity())
		// Add extra filter here
		if filter.IgnoreDefaultTable {
			moved = moved.Where("ID != 1")
		}
		return moved
	}
}
