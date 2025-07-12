package service

import (
	"sort"

	"github.com/Catizard/lampghost_wails/internal/dto"
	"github.com/Catizard/lampghost_wails/internal/entity"
	"github.com/Catizard/lampghost_wails/internal/vo"
	"github.com/rotisserie/eris"
	"gorm.io/gorm"
)

type CustomCourseService struct {
	db *gorm.DB
}

func NewCustomCourseService(db *gorm.DB) *CustomCourseService {
	return &CustomCourseService{
		db: db,
	}
}

func (s *CustomCourseService) AddCustomCourse(param *vo.CustomCourseVo) error {
	if param == nil {
		return eris.Errorf("AddCustomCourse: param cannot be nil")
	}
	if param.Name == "" {
		return eris.Errorf("AddCustomCourse: name cannot be empty")
	}
	return s.db.Transaction(func(tx *gorm.DB) error {
		return addCustomCourse(tx, param)
	})
}

func (s *CustomCourseService) FindCustomCourseList(filter *vo.CustomCourseVo) (out []*entity.CustomCourse, n int, err error) {
	err = s.db.Transaction(func(tx *gorm.DB) error {
		out, n, err = findCustomCourseList(tx, filter)
		return err
	})
	return
}

func (s *CustomCourseService) QueryCustomCourseSongListWithRival(filter *vo.CustomCourseVo) (out []*dto.RivalSongDataDto, n int, err error) {
	if filter == nil {
		err = eris.Errorf("QueryCustomCourseSongListWithRival: filter cannot be nil")
		return
	}
	if filter.ID == 0 {
		err = eris.Errorf("QueryCustomCourseSongListWithRival: ID cannot be 0")
		return
	}
	if filter.RivalID == 0 {
		err = eris.Errorf("QueryCustomCourseSongListWithRival: RivalID cannot be 0")
		return
	}
	err = s.db.Transaction(func(tx *gorm.DB) error {
		rawSongs, _, err := findCustomCourseDataListByID(tx, filter.ID)
		if err != nil {
			return err
		}
		md5s := make([]string, 0)
		md5ToSong := make(map[string]*entity.CustomCourseData)
		for _, rawSong := range rawSongs {
			md5s = append(md5s, rawSong.Md5)
			md5ToSong[rawSong.Md5] = rawSong
		}
		if len(md5s) == 0 {
			out = make([]*dto.RivalSongDataDto, 0)
			return nil
		}
		queryParam := &vo.RivalSongDataVo{
			RemoveDuplicate: true,
			RivalId:         filter.RivalID,
			Md5s:            md5s,
		}
		out, n, err = findRivalSongDataList(tx, queryParam)
		// NOTE: Inherit some fields from rawSong
		for _, r := range out {
			r.ID = md5ToSong[r.Md5].ID
			r.OrderNumber = md5ToSong[r.Md5].OrderNumber
		}
		sort.Slice(out, func(i, j int) bool {
			return out[i].OrderNumber < out[j].OrderNumber
		})
		return err
	})
	return
}

func (s *CustomCourseService) AddCustomCourseData(param *entity.CustomCourseData) error {
	if param == nil {
		return eris.Errorf("AddCustomCourseData: param cannot be nil")
	}
	if param.CustomCourseID == 0 {
		return eris.Errorf("AddCustomCourseData: custom course's id cannot 0")
	}
	if param.Sha256 == "" {
		return eris.Errorf("AddCustomCourseData: sha256 cannot be empty")
	}
	if param.Md5 == "" {
		return eris.Errorf("AddCustomCourseData: md5 cannot be empty")
	}
	return s.db.Transaction(func(tx *gorm.DB) error {
		siblings, _, err := findCustomCourseDataListByID(tx, param.CustomCourseID)
		if err != nil {
			return err
		}
		for _, sib := range siblings {
			if sib.Md5 == param.Md5 || sib.Sha256 == param.Sha256 {
				return eris.Errorf("Cannot add duplicated song into one course")
			}
		}
		return addCustomCourseData(tx, param)
	})
}

func (s *CustomCourseService) BindSongToCustomCourse(sha256, md5 string, customCourseID uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		return addCustomCourseData(tx, &entity.CustomCourseData{
			CustomCourseID: customCourseID,
			Sha256:         sha256,
			Md5:            md5,
		})
	})
}

func (s *CustomCourseService) UpdateCustomCourseOrder(courseIDs []uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		return updateCustomCourseOrder(tx, courseIDs)
	})
}

func (s *CustomCourseService) UpdateCustomCourseDataOrder(courseDataIDs []uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		return updateCustomCourseDataOrder(tx, courseDataIDs)
	})
}

func (s *CustomCourseService) DeleteCustomCourse(courseID uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		candidate, err := findCustomCourseByID(tx, courseID)
		if err != nil {
			return err
		}
		if err := tx.Unscoped().Where("custom_course_id = ?", candidate.ID).Delete(&entity.CustomCourseData{}).Error; err != nil {
			return err
		}
		if err := tx.Unscoped().Delete(&entity.CustomCourse{}, candidate.ID).Error; err != nil {
			return err
		}
		return nil
	})
}

func (s *CustomCourseService) DeleteCustomCourseData(courseDataID uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		return tx.Unscoped().Delete(&entity.CustomCourseData{}, courseDataID).Error
	})
}

func findCustomCourseList(tx *gorm.DB, filter *vo.CustomCourseVo) (out []*entity.CustomCourse, n int, err error) {
	err = tx.Model(&entity.CustomCourse{}).Scopes(scopeCustomCourseFilter(filter)).Order("order_number").Find(&out).Error
	n = len(out)
	return
}

func findCustomCourseDataListByID(tx *gorm.DB, courseID uint) (out []*entity.CustomCourseData, n int, err error) {
	err = tx.Model(&entity.CustomCourseData{}).Where("custom_course_id = ?", courseID).Find(&out).Error
	n = len(out)
	return
}

func findCustomCourseByID(tx *gorm.DB, id uint) (course *entity.CustomCourse, err error) {
	err = tx.First(&course, id).Error
	return
}

func addCustomCourse(tx *gorm.DB, param *vo.CustomCourseVo) error {
	return tx.Create(param.Entity()).Error
}

func addCustomCourseData(tx *gorm.DB, param *entity.CustomCourseData) error {
	return tx.Create(param).Error
}

func updateCustomCourseOrder(tx *gorm.DB, courseIDs []uint) error {
	if len(courseIDs) == 0 {
		return nil
	}
	for i, courseID := range courseIDs {
		entity := &entity.CustomCourse{}
		entity.ID = courseID
		if err := tx.Model(entity).Update("order_number", i).Error; err != nil {
			return err
		}
	}
	return nil
}

func updateCustomCourseDataOrder(tx *gorm.DB, courseDataIDs []uint) error {
	if len(courseDataIDs) == 0 {
		return nil
	}
	for i, courseID := range courseDataIDs {
		entity := &entity.CustomCourseData{}
		entity.ID = courseID
		if err := tx.Debug().Model(entity).Update("order_number", i).Error; err != nil {
			return err
		}
	}
	return nil
}

func scopeCustomCourseFilter(filter *vo.CustomCourseVo) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if filter == nil {
			return db
		}
		moved := db.Where(filter.Entity())
		// Extra filters here
		return moved
	}
}
