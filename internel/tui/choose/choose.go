package choose

import (
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	choices  []string
	cursor   int
	selected int
	msg      string
}

func newModel(choices []string, msg string) model {
	return model{
		choices:  choices,
		selected: 0,
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
			if m.cursor > 0 {
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
	m := newModel(choices, msg)
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		log.Fatalf("OpenChooseTui returns a error: %v\n", err)
	}
	return m.selected
}
