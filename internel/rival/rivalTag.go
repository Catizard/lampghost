package rival

import (
	"github.com/Catizard/lampghost/internel/common"
	"github.com/Catizard/lampghost/internel/difftable"
	"github.com/Catizard/lampghost/internel/score"
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
	for _, course := range courseArr {
		var sha256 string
		// Iteration on plain array should be sequential
		for _, md5 := range course.Md5 {
			if songData, ok := md5MapsToSongData[md5]; ok {
				sha256 += songData.Sha256
			}
		}
		interestSha256[sha256] = struct{}{}
	}

	// Maps scorelog by sha256
	sha256MapsToScoreLog := make(map[string][]score.ScoreLog)
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

	// For now, only "first clear" and "first hard clear" is registered
	return nil
}
