package score

import (
	"fmt"
	"strings"

	"github.com/Catizard/lampghost/internal/common/source"
	"github.com/charmbracelet/log"
	"github.com/guregu/null/v5"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type ScoreLog struct {
	Sha256    string
	Mode      string
	Clear     int32
	OldClear  int32
	Score     int32
	OldScore  int32
	Combo     int32
	OldCombo  int32
	Minbp     int32
	OldMinbp  int32
	TimeStamp int64 `db:"date"`
}

type ScoreLogFilter struct {
	BeginTime null.Int `db:"beginTime"`
	EndTime   null.Int `db:"endTime"`
}

func (f *ScoreLogFilter) GenerateWhereClause() string {
	where := []string{"1=1"}
	if v := f.BeginTime; v.Valid {
		where = append(where, "date>=:beginTime")
	}
	if v := f.EndTime; v.Valid {
		where = append(where, "date<=:endTime")
	}
	return strings.Join(where, " AND ")
}

func ReadScoreLogFromSqlite(filePath string, filter ScoreLogFilter) ([]ScoreLog, error) {
	if !strings.HasSuffix(filePath, ".db") {
		return nil, fmt.Errorf("try reading from sqlite database file while path doesn't contains a .db suffix")
	}
	db, err := sqlx.Open("sqlite3", filePath)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.NamedQuery("select * from scorelog where "+filter.GenerateWhereClause(), filter)
	if err != nil {
		return nil, err
	}
	scoreLogArray := make([]ScoreLog, 0)
	for rows.Next() {
		var log ScoreLog
		err = rows.StructScan(&log)
		if err != nil {
			return nil, err
		}
		scoreLogArray = append(scoreLogArray, log)
	}
	return scoreLogArray, nil
}

type LR2ScoreLog struct {
	Md5        string      `db:"hash"`
	Clear      int         `db:"clear"`
	Perfect    int         `db:"perfect"`
	Great      int         `db:"great"`
	Good       int         `db:"good"`
	Bad        int         `db:"bad"`
	Poor       int         `db:"poor"`
	TotalNotes int         `db:"totalnotes"`
	MaxCombo   int         `db:"maxcombo"`
	MinBP      int         `db:"minbp"`
	PlayCount  int         `db:"playcount"`
	ClearCount int         `db:"clearcount"`
	FailCount  int         `db:"failcount"`
	Rank       int         `db:"rank"`
	Rate       int         `db:"rate"`
	ClearDB    int         `db:"clear_db"`
	OPHistory  int         `db:"op_history"`
	ScoreHash  string      `db:"scorehash"`
	Ghost      null.String `db:"ghost"`
	ClearSD    int         `db:"clear_sd"`
	ClearEX    int         `db:"clear_ex"`
	OPBest     int         `db:"op_best"`
	RSeed      int         `db:"rseed"`
	Complete   int         `db:"complete"`
	RowId      int         `db:"rowid"`
}

type CommonScoreLog struct {
	Sha256    null.String
	Md5       null.String
	Clear     int32
	TimeStamp null.Int
	LogType   string
}

func NewCommonScoreLogFromOraja(log *ScoreLog) *CommonScoreLog {
	return &CommonScoreLog{
		Sha256:    null.StringFrom(log.Sha256),
		Clear:     log.Clear,
		TimeStamp: null.IntFrom(log.TimeStamp),
	}
}

func NewCommonScoreLogFromLR2(scoreLog *LR2ScoreLog) *CommonScoreLog {
	// Hack: Take row_id as timestamp
	ret := &CommonScoreLog{
		Md5:       null.StringFrom(scoreLog.Md5),
		Clear:     int32(scoreLog.Clear),
		TimeStamp: null.IntFrom(int64(scoreLog.RowId)),
	}
	// Note: Workaround for LR2 log, because LR2 doesn't have assist lamp
	// So except no play(=0) and fail(=1), clear should += 2
	if ret.Clear > 1 {
		ret.Clear += 2
	}
	// Hack: If len(MD5) > 32, erase the first 32 numbers
	// Then join the last every 32 numbers with ','
	hash := ret.Md5.ValueOrZero()
	// If len cannot be divided by 32, report it
	if len(hash) % 32 != 0 { 
		log.Warnf("Corrupted data: len(hash) cannot divide 32, report this to the author.")
	}
	if len(hash) > 32 {
		hash = hash[32:]
		joined := ""
		for i := 0; i < len(hash); i += 32 {
			if i != 0 {
				joined += ","
			}
			joined += hash[i:i+32]
		}
		ret.Md5 = null.StringFrom(joined)
		ret.LogType = source.Course
	} else {
		ret.LogType = source.Song
	}
	return ret
}

// Oraja songdata.db entity
type SongData struct {
	Md5        string
	Sha256     string
	Title      string
	SubTitle   string
	Genre      string
	Artist     string
	SubArtist  string
	Tag        string
	Path       string
	Folder     string
	StageFile  string
	Banner     string
	BackBmp    string
	Preview    string
	Parent     string
	Level      int32
	Difficulty int32
	MaxBpm     int32
	MinBpm     int32
	Length     int32
	Mode       int32
	Judge      int32
	Feature    int32
	Content    int32
	Date       int64
	Favorite   int32
	AddDate    int64
	Notes      int32
	ChartHash  string
}
