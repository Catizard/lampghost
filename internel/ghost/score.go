package ghost

type ScoreLog struct {
	Sha256   string
	Mode     string
	Clear    int32
	OldClear int32
	Score    int32
	OldScore int32
	Combo    int32
	OldCombo int32
	Minbp    int32
	OldMinbp int32
	Date     int64
}

func ReadScoreLogFromSqlite(filePath string) []ScoreLog {

}