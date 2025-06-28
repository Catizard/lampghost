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
	SubTitle   string `gorm:"column:subtitle"`
	Genre      string
	Artist     string
	SubArtist  string `gorm:"column:subartist"`
	Tag        string
	Path       string
	Folder     string
	StageFile  string `gorm:"column:stagefile"`
	Banner     string
	BackBmp    string `gorm:"column:backbmp"`
	Preview    string
	Parent     string
	Level      int32
	Difficulty int32
	MaxBpm     int32 `gorm:"column:maxbpm"`
	MinBpm     int32 `gorm:"column:minbpm"`
	Length     int32
	Mode       int32
	Judge      int32
	Feature    int32
	Content    int32
	Date       int64
	Favorite   int32
	AddDate    int64 `gorm:"column:adddate"`
	Notes      int32
	ChartHash  string `gorm:"column:charthash"`
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
	Seed       int64
	Random     int32
	TimeStamp  int64 `gorm:"column:date"`
	State      int32
}

func (ScoreDataLog) TableName() string {
	return "scoredatalog"
}

type BeatorajaSongMode int

const (
	BeatorajaBeat5K  = 5
	BeatorajaBeat7K  = 7
	BeatorajaBeat10K = 10
	BeatorajaBeat14K = 14
	// BeatorajaPopn5K = 9 We cannot distinguish a mode is popn5k or popn9k
	BeatorajaPopn9K            = 9
	BeatorajaKeyboard24K       = 25
	BeatorajaKeyboard24KDouble = 50
)

var beatorajaSongModeName = map[BeatorajaSongMode]string{
	BeatorajaBeat5K:  "beat-5k",
	BeatorajaBeat7K:  "beat-7k",
	BeatorajaBeat10K: "beat-10k",
	BeatorajaBeat14K: "beat-14k",
	// BeatorajaPopn5K: "popn-5k",
	BeatorajaPopn9K:            "popn-9k",
	BeatorajaKeyboard24K:       "keyboard-24k",
	BeatorajaKeyboard24KDouble: "keyboard-24k-double",
}

func (s BeatorajaSongMode) String() string {
	return beatorajaSongModeName[s]
}

type IRPlayer struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Rank string `json:"rank"`
}

// Part of the ScoreData, the missing fields are filled by Lampghost-IR
// The reason do such a simplify is Beatoraja's save file is incomplete,
// which makes Lampghost cannot build completely historical play records
// but only lamp data
type IRLampData struct {
	Clear  int    `json:"clear"`
	Sha256 string `json:"sha256"`
	Mode   int    `json:"mode"`
}
