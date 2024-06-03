package choose

import (
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/log"
)

// Open a tui application for choices
// Returns index of user's choice
// Fast return if there is no need to choose
func OpenChooseTuiSkippable(choices []string, msg string) int {
	return OpenChooseTui(choices, msg, true)
}

func OpenChooseTui(choices []string, msg string, skip bool) int {
	if len(choices) == 0 {
		panic("OpenChooseTui: len(choices) == 0")
	}
	if len(choices) == 1 && skip {
		return 0 // fast fail
	}
	options := make([]huh.Option[int], 0)
	for i, v := range choices {
		options = append(options, huh.Option[int]{Key: v, Value: i})
	}
	var ret int
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[int]().
				Title(msg).
				Options(options...).
				Value(&ret),
		),
	)
	if err := form.Run(); err != nil {
		log.Fatal(err)
	}
	return ret
}

// Simple wrapper of OpenChooseTui
// With two choices: yes and no
// Returns false if choose no, true if choose yes
func OpenYesOrNoChooseTui(msg string) bool {
	return huh.NewConfirm().
		Title(msg).
		Affirmative("Yes").
		Negative("No").GetValue().(bool)
}
