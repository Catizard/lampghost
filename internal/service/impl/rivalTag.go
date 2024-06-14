package impl

import (
	"fmt"
	"sort"

	"github.com/Catizard/lampghost/internal/common/clearType"
	"github.com/Catizard/lampghost/internal/data/difftable"
	"github.com/Catizard/lampghost/internal/data/rival"
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
	if len(courseArr) == 0 {
		return nil
	}
	if err := r.LoadData(nil); err != nil {
		panic(err)
	}

	// TODO: Structure of this procdure should be changed
	// Maps songdata by md5
	md5MapsToSongData := make(map[string]score.SongData)
	for _, v := range r.SongData {
		md5MapsToSongData[v.Md5] = v
	}

	// TODO: genocide 2018 course contains some mutations like no speed/no good. They contain the exactly same md5 so sha256 would be the same
	// This causes tag generation actually works incorectly, I think delete them directly from course would be the easist way to handle this problem
	// At current development stage, the specific implementation is difftable/saveCourseInfoFromTableHeader (Actaully hasn't implemented now)
	interestSha256 := make(map[string]struct{}, 0)
	for i, course := range courseArr {
		var sha256 string
		valid := true
		// Iteration on plain array should be sequential
		for _, md5 := range course.Md5 {
			if songData, ok := md5MapsToSongData[md5]; ok {
				sha256 += songData.Sha256
			} else {
				valid = false
			}
		}
		if !valid {
			log.Warnf("Course %s builds up failed due to lack of data, songdata path=%s", course.Name, r.SongDataPath.ValueOrZero())
			continue
		}
		log.Debug("course name=%s, sha256=%s\n", course.Name, sha256)
		courseArr[i].Sha256s = sha256
		interestSha256[sha256] = struct{}{}
	}

	// Maps scorelog by sha256
	sha256MapsToScoreLog := make(map[string][]score.ScoreLog)
	// Before doing iteration, make sure scorelog is sorted by time
	sort.Slice(r.ScoreLog, func(i, j int) bool {
		return r.ScoreLog[i].TimeStamp < r.ScoreLog[j].TimeStamp
	})
	for _, v := range r.ScoreLog {
		// Skip
		if _, ok := interestSha256[v.Sha256]; !ok {
			continue
		}
		if _, ok := sha256MapsToScoreLog[v.Sha256]; !ok {
			sha256MapsToScoreLog[v.Sha256] = make([]score.ScoreLog, 0)
		}
		sha256MapsToScoreLog[v.Sha256] = append(sha256MapsToScoreLog[v.Sha256], v)
	}

	// For now, only "first clear" and "first hard clear" tags are generated

	// TODO: extract this part out
	tags := make([]*rival.RivalTag, 0)
	// First Clear Tag
	for _, course := range courseArr {
		if logs, ok := sha256MapsToScoreLog[course.Sha256s]; !ok {
			continue // No record, continue
		} else {
			for _, log := range logs {
				if log.Clear >= clearType.Normal {
					fct := rival.RivalTag{
						TagName:   course.Name + " First Clear",
						Generated: true,
						TimeStamp: log.TimeStamp,
					}
					tags = append(tags, &fct)
					break
				}
			}
		}
	}
	// First Hard Clear Tag
	for _, course := range courseArr {
		if logs, ok := sha256MapsToScoreLog[course.Sha256s]; !ok {
			continue // No record, continue
		} else {
			for _, log := range logs {
				if log.Clear >= clearType.Hard {
					fct := rival.RivalTag{
						TagName:   course.Name + " First Hard Clear",
						Generated: true,
						TimeStamp: log.TimeStamp,
					}
					tags = append(tags, &fct)
					break
				}
			}
		}
	}
	// Add rival's id all together
	for i := range tags {
		tags[i].RivalId = r.Id
	}
	// Sync rival's generated tag
	if err := syncGeneratedTags(tx, r, tags); err != nil {
		panic(err)
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
	_, err := tx.NamedExec(`INSERT INTO rival_tag(rival_id,tag_name,generated,timestamp) VALUES(:rival_id,:tag_name,:generated,:timestamp)`, rivalTag)
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
