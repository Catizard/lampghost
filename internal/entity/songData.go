package entity

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
