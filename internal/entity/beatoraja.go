package entity

// Beatoraja related database entity definition
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
	TimeStamp int64 `gorm:"column:date"`
}

func (ScoreLog) TableName() string {
	return "scorelog"
}

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

func (SongData) TableName() string {
	return "song"
}

type ScoreDataLog struct {
	Sha256     string
	Mode       string
	Clear      int32
	Epg        int32
	Lpg        int32
	Egr        int32
	Lgr        int32
	Egd        int32
	Lgd        int32
	Ebd        int32
	Lbd        int32
	Epr        int32
	Lpr        int32
	Ems        int32
	Lms        int32
	Notes      int32
	Combo      int32
	Minbp      int32
	PlayCount  int32
	ClearCount int32
	Option     int32
	Seed       int32
	Random     int32
	TimeStamp  int64 `gorm:"column:date"`
	State      int32
}

func (ScoreDataLog) TableName() string {
	return "scoredatalog"
}
