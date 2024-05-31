package choose

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
)

var ynChoices []string = []string{"no", "yes"}

type model struct {
	choices  []string
	cursor   int
	selected int
	msg      string
}

func newModel(choices []string, msg string) model {
	return model{
		choices:  choices,
		selected: -1,
		msg:      msg,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor+1 < len(m.choices) {
				m.cursor++
			}
		case "enter", " ":
			m.selected = m.cursor
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
	var s string
	if len(m.msg) > 0 {
		s += m.msg + "\n"
	}

	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		s += fmt.Sprintf("%s %s\n", cursor, choice)
	}

	return s
}

// Open a tui application for choices
// Returns index of user's choice
// Fast return if there is no need to choose
func OpenChooseTui(choices []string, msg string) int {
	if len(choices) == 0 {
		panic("OpenChooseTui: len(choices) == 0")
	}
	if len(choices) == 1 {
		return 0 // fast fail
	}
	m := newModel(choices, msg)
	p := tea.NewProgram(m)

	// Result of p.Run() is tea.Model other than model
	// So here we have to declare a new variable and perform cast on it
	rm, err := p.Run()
	if err != nil {
		log.Fatalf("OpenChooseTui returns a error: %v\n", err)
	}
	if rm, ok := rm.(model); ok && rm.selected == -1 {
		log.Fatalf("You didn't make any choices") // -----+
	} else { //                                           |
		return rm.selected //                             |
	} //                                                  |
	// Cannot be reached   // <---------------------------+
	return -1
}

// Simple wrapper of OpenChooseTui
// With two choices: yes and no
// Returns false if choose no, true if choose yes
func OpenYesOrNoChooseTui(msg string) bool {
	i := OpenChooseTui(ynChoices, msg)
	return i != 0
}
