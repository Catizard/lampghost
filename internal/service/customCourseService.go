package service

import (
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
		for _, rawSong := range rawSongs {
			md5s = append(md5s, rawSong.Md5)
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
