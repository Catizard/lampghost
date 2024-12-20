package entity

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
	HardPlus
)

const (
	STR_NO_PLAY         = "NO_PLAY"
	STR_Failed          = "FAILED"
	STR_AssistEasy      = "Assist Easy"
	STR_LightAssistEasy = "Light Assist Easy"
	STR_Easy            = "Easy"
	STR_Normal          = "Normal"
	STR_Hard            = "Hard"
	STR_ExHard          = "EX Hard"
	STR_FullCombo       = "FULL COMBO"
	STR_Perfect         = "PERFECT"
	STR_Max             = "MAX"
	STR_HC_PLUS         = "HC+"
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

// Generated code, do not edit
// NOTE: 转换所有等于和高于Hard Clear的评价为HC+
func ConvOrajaIgnoreHardPlus(v int32) string {
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
		return STR_HC_PLUS
	}
	if v == ExHard {
		return STR_HC_PLUS
	}
	if v == FullCombo {
		return STR_HC_PLUS
	}
	if v == Perfect {
		return STR_HC_PLUS
	}
	if v == Max {
		return STR_HC_PLUS
	}
	panic("unexpected clear type")
}
