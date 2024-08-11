package clearType

const (
	NO_PLAY = iota
	Failed
	AssistEasy
	LightAssistEasy
	Easy
	Normal
	Hard
	ExHard
	FullCombo
	Perfect
	Max
)

const (
	STR_NO_PLAY = "NO_PLAY"
	STR_Failed = "FAILED"
	STR_AssistEasy = "Assist Easy"
	STR_LightAssistEasy = "Light Assist Easy"
	STR_Easy = "Easy"
	STR_Normal = "Normal"
	STR_Hard = "Hard"
	STR_ExHard = "EX Hard"
	STR_FullCombo = "FULL COMBO"
	STR_Perfect = "PERFECT"
	STR_Max = "MAX"
)

// Generated code, do not edit
func ConvOraja(v int32) string {
	if v == NO_PLAY {
		return STR_NO_PLAY
	}
	if v == Failed {
		return STR_Failed
	}
	if v == AssistEasy {
		return STR_AssistEasy
	}
	if v == LightAssistEasy {
		return STR_LightAssistEasy
	}
	if v == Easy {
		return STR_Easy
	}
	if v == Normal {
		return STR_Normal
	}
	if v == Hard {
		return STR_Hard
	}
	if v == ExHard {
		return STR_ExHard
	}
	if v == FullCombo {
		return STR_FullCombo
	}
	if v == Perfect {
		return STR_Perfect
	}
	if v == Max {
		return STR_Max
	}
	panic("unexpected clear type")
}

// hack: LR2 doesn't have assist lamp
// So except no play(=0) and fail(=1), clear should += 2
func ConvLR2(v int32) string {
	if v > 1 {
		v += 2
	}
	return ConvOraja(v)
}