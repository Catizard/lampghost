package impl

import (
	"fmt"

	"github.com/Catizard/lampghost/internal/common/source"
	"github.com/Catizard/lampghost/internal/data/difftable"
	"github.com/Catizard/lampghost/internal/data/rival"
	"github.com/Catizard/lampghost/internal/data/rival/builder"
	"github.com/Catizard/lampghost/internal/data/score"
	"github.com/Catizard/lampghost/internal/sqlite"
	"github.com/Catizard/lampghost/internal/tui/choose"
	"github.com/charmbracelet/log"
	"github.com/guregu/null/v5"
)

var _ rival.RivalTagService = (*RivalTagService)(nil)

type RivalTagService struct {
	db *sqlite.DB
}

func NewRivalTagService(db *sqlite.DB) *RivalTagService {
	return &RivalTagService{db: db}
}

func (s *RivalTagService) FindRivalTagList(filter rival.RivalTagFilter) ([]*rival.RivalTag, int, error) {
	tx, err := s.db.BeginTx()
	if err != nil {
		return nil, 0, err
	}
	defer tx.Rollback()
	return findRivalTagList(tx, filter)
}

func (s *RivalTagService) FindRivalTagById(id int) (*rival.RivalTag, error) {
	tx, err := s.db.BeginTx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	return findRivalTagById(tx, id)
}

func (s *RivalTagService) InsertRivalTag(tag *rival.RivalTag) error {
	tx, err := s.db.BeginTx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := insertRivalTag(tx, tag); err != nil {
		return err
	}
	return tx.Commit()
}

func (s *RivalTagService) DeleteRivalTagById(id int) error {
	tx, err := s.db.BeginTx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := deleteRivalTagById(tx, id); err != nil {
		return err
	}
	return tx.Commit()
}

func (s *RivalTagService) DeleteRivalTag(filter rival.RivalTagFilter) error {
	tx, err := s.db.BeginTx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := deleteRivalTag(tx, filter); err != nil {
		return err
	}
	return tx.Commit()
}

func (s *RivalTagService) ChooseOneTag(msg string, filter rival.RivalTagFilter) (*rival.RivalTag, error) {
	tx, err := s.db.BeginTx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	if arr, n, err := findRivalTagList(tx, filter); err != nil {
		return nil, err
	} else if n == 0 {
		return nil, fmt.Errorf("no tag data")
	} else {
		choices := make([]string, 0)
		for _, v := range arr {
			choices = append(choices, v.String())
		}
		i := choose.OpenChooseTuiSkippable(choices, msg)
		return arr[i], nil
	}
}

func (s *RivalTagService) BuildTags(r *rival.RivalInfo, courseArr []*difftable.CourseInfo) error {
	tx, err := s.db.BeginTx()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if err := buildTags(tx, r, courseArr); err != nil {
		return err
	}
	return tx.Commit()
}

func (s *RivalTagService) SyncGeneratedTags(r *rival.RivalInfo, tags []*rival.RivalTag) error {
	tx, err := s.db.BeginTx()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if err := syncGeneratedTags(tx, r, tags); err != nil {
		return err
	}
	return tx.Commit()
}

func findRivalTagList(tx *sqlite.Tx, filter rival.RivalTagFilter) (_ []*rival.RivalTag, _ int, err error) {
	rows, err := tx.NamedQuery("SELECT * FROM rival_tag WHERE "+filter.GenerateWhereClause(), filter)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	ret := make([]*rival.RivalTag, 0)
	for rows.Next() {
		r := &rival.RivalTag{}
		if err := rows.StructScan(r); err != nil {
			return nil, 0, err
		}

		ret = append(ret, r)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return ret, len(ret), nil
}

func findRivalTagById(tx *sqlite.Tx, id int) (*rival.RivalTag, error) {
	arr, _, err := findRivalTagList(tx, rival.RivalTagFilter{Id: null.IntFrom(int64(id))})
	if err != nil {
		return nil, err
	} else if len(arr) == 0 {
		return nil, fmt.Errorf("panic: no data")
	}
	return arr[0], nil
}

func insertRivalTag(tx *sqlite.Tx, rivalTag *rival.RivalTag) error {
	_, err := tx.NamedExec(`INSERT INTO rival_tag(rival_id,tag_name,generated,timestamp,tag_source) VALUES(:rival_id,:tag_name,:generated,:timestamp,:tag_source)`, rivalTag)
	return err
}

func deleteRivalTagById(tx *sqlite.Tx, id int) error {
	if _, err := findRivalTagById(tx, id); err != nil {
		return err
	}

	_, err := tx.Exec("DELETE FROM rival_tag WHERE id=?", id)
	return err
}

func deleteRivalTag(tx *sqlite.Tx, filter rival.RivalTagFilter) error {
	_, err := tx.NamedExec("DELETE FROM rival_tag WHERE "+filter.GenerateWhereClause(), filter)
	return err
}

func syncGeneratedTags(tx *sqlite.Tx, r *rival.RivalInfo, tags []*rival.RivalTag) error {
	if len(tags) == 0 {
		log.Warn("No tags to sync")
		return nil // Okay dokey
	}
	filter := rival.RivalTagFilter{
		RivalId:   null.IntFrom(int64(r.Id)),
		Generated: null.BoolFrom(true),
	}
	if err := deleteRivalTag(tx, filter); err != nil {
		return err
	}
	// TODO: I'm too lazy to generate a batch insert call...
	for _, tag := range tags {
		if err := insertRivalTag(tx, tag); err != nil {
			return err
		}
	}
	return nil
}

func buildTags(tx *sqlite.Tx, r *rival.RivalInfo, courseArr []*difftable.CourseInfo) error {
	if len(courseArr) == 0 {
		return nil
	}

	// Preparation
	headers, _, err := findDiffTableHeaderList(tx, difftable.DiffTableHeaderFilter{})
	if err != nil {
		return err
	}
	courses, _, err := findCourseInfoList(tx, difftable.CourseInfoFilter{})
	if err != nil {
		return err
	}
	// Split logs
	songScoreLog := make([]*score.CommonScoreLog, 0)
	courseScoreLog := make([]*score.CommonScoreLog, 0)
	for _, v := range r.CommonScoreLog {
		if v.LogType == source.Course {
			courseScoreLog = append(courseScoreLog, v)
		} else {
			songScoreLog = append(songScoreLog, v)
		}
	}

	tags := builder.Build(builder.TagBuildParam{
		RivalInfo:       r,
		DiffTableHeader: headers,
		SongScoreLog:    songScoreLog,
		CourseScoreLog:  courseScoreLog,
		Courses:         courses,
	})
	log.Infof("Generated %d tags for [%s]", len(tags), r.Name)
	// Add rival's id all together
	for i := range tags {
		tags[i].RivalId = r.Id
	}
	// Sync rival's generated tag
	if err := syncGeneratedTags(tx, r, tags); err != nil {
		return err
	}
	return nil
}
