package rival

import (
	"sort"

	"github.com/Catizard/lampghost/internel/common"
	"github.com/Catizard/lampghost/internel/common/clearType"
	"github.com/Catizard/lampghost/internel/difftable"
	"github.com/Catizard/lampghost/internel/score"
	"github.com/charmbracelet/log"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type RivalTag struct {
	Id        int    `db:"id"`
	RivalId   int    `db:"rival_id"`
	TagName   string `db:"tag_name"`
	Generated bool   `db:"generated"`
	TimeStamp int64  `db:"timestamp"`
}

func InitRivalTagTable() error {
	db, err := sqlx.Open("sqlite3", common.DBFileName)
	if err != nil {
		return err
	}
	defer db.Close()
	_, err = db.Exec("DROP TABLE IF EXISTS 'rival_tag';CREATE TABLE rival_tag (id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, rival_id INTEGER NOT NULL, tag_name TEXT(255) NOT NULL, 'generated' INTEGER DEFAULT (0) NOT NULL, 'timestamp' TEXT NOT NULL)")
	return err
}

// Build tags for one rival based on passed course data
// Note: this function can be seen as re-build all generated tags
func (r *RivalInfo) BuildTags(courseArr []difftable.CourseInfo) error {
	if len(courseArr) == 0 {
		return nil
	}
	if err := r.LoadDataIfNil(); err != nil {
		panic(err)
	}

	// TODO: only calculate sha256 once
	// Maps songdata by md5
	md5MapsToSongData := make(map[string]score.SongData)
	for _, v := range r.SongData {
		md5MapsToSongData[v.Md5] = v
	}

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
			log.Warnf("Course %s builds up failed due to lack of data, songdata path=%s", course.Name, r.SongDataPath)
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
	tags := make([]RivalTag, 0)
	// First Clear Tag
	for _, course := range courseArr {
		if logs, ok := sha256MapsToScoreLog[course.Sha256s]; !ok {
			continue // No record, continue
		} else {
			for _, log := range logs {
				if log.Clear >= clearType.Normal {
					fct := RivalTag{
						TagName:   course.Name + " First Clear",
						Generated: true,
						TimeStamp: log.TimeStamp,
					}
					tags = append(tags, fct)
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
					fct := RivalTag{
						TagName:   course.Name + " First Hard Clear",
						Generated: true,
						TimeStamp: log.TimeStamp,
					}
					tags = append(tags, fct)
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
	if err := r.syncGeneratedTags(tags); err != nil {
		panic(err)
	}
	return nil
}

// Sync one rival's generated tag
// Protected by transaction
func (r *RivalInfo) syncGeneratedTags(tags []RivalTag) error {
	if len(tags) == 0 {
		log.Warn("No tags to sync")
		// OK, nothing to do
		return nil
	}
	db := common.OpenDB()
	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}

	// 1) Delete any generated tags
	if _, err := db.NamedExec("DELETE FROM 'rival_tag' WHERE generated=true AND rival_id=:id", r); err != nil {
		tx.Rollback()
		return err
	}
	// 2) Insert tags
	// TODO: I'm too lazy to generate a batch insert call...
	for _, tag := range tags {
		if _, err := db.NamedExec("INSERT INTO rival_tag(rival_id,tag_name,generated,timestamp) VALUES(:rival_id,:tag_name,:generated,:timestamp)", &tag); err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()
	return nil
}
