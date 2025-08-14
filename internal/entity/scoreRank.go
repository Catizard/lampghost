package entity

var ScoreRanks []ScoreRank = []ScoreRank{
	{"MAX", 100.00, 9},
	{"AAA", 88.88, 8},
	{"AA", 77.77, 7},
	{"A", 66.66, 6},
	{"B", 55.55, 5},
	{"C", 44.44, 4},
	{"D", 33.33, 3},
	{"E", 22.22, 2},
	{"F", 0, 1},
	{"NO_PLAY", 0, 0},
}

type ScoreRank struct {
	Name  string
	Low   float32
	Value int
}

func GetScoreRank(acc float32) ScoreRank {
	for _, scoreRank := range ScoreRanks {
		if acc >= scoreRank.Low {
			return scoreRank
		}
	}
	return ScoreRanks[len(ScoreRanks)-1]
}
