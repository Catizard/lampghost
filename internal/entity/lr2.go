package entity

type LR2Log struct {
	MD5        string `gorm:"column:hash"`
	Clear      int    `gorm:"column:clear"`
	Perfect    int    `gorm:"column:perfect"`
	Great      int    `gorm:"column:great"`
	Good       int    `gorm:"column:good"`
	Bad        int    `gorm:"column:bad"`
	Poor       int    `gorm:"column:Poor"`
	TotalNotes int    `gorm:"column:totalnotes"`
	MaxCombo   int    `gorm:"column:maxcombo"`
	Minbp      int    `gorm:"column:minbp"`
	PlayCount  int    `gorm:"column:playcount"`
	ClearCount int    `gorm:"column:clearcount"`
	FailCount  int    `gorm:"column:failcount"`
	Rank       int    `gorm:"column:rank"`
	Rate       int    `gorm:"column:rate"`
	ClearDB    int    `gorm:"column:clear_db"`
	OpHistory  int    `gorm:"column:op_history"`
	ScoreHash  string `gorm:"column:scorehash"`
	Ghost      string `gorm:"column:ghost"`
	ClearSD    int    `gorm:"column:clear_sd"`
	ClearEX    int    `gorm:"column:clear_ex"`
	OpBest     int    `gorm:"column:op_best"`
	RSeed      int    `gorm:"column:rseed"`
	Complete   int    `gorm:"column:complete"`
	RowID      int    `gorm:"column:row_id"`
}

func (LR2Log) TableName() string {
	return "score"
}

func ConvLR2Clear(x int) int {
	if x == 5 {
		return FullCombo
	} else if x > 1 {
		return x + 2
	}
	return x
}
