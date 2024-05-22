package ghost

import (
	"log"
	"time"

	"github.com/Catizard/lampghost/internel/difftable"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

// sessionState is used to track which model is focused
type sessionState uint

const (
	defaultTime              = time.Minute
	timerView   sessionState = iota
	spinnerView
)

var (
	// Available spinners
	spinners = []spinner.Spinner{
		spinner.Line,
		spinner.Dot,
		spinner.MiniDot,
		spinner.Jump,
		spinner.Pulse,
		spinner.Points,
		spinner.Globe,
		spinner.Moon,
		spinner.Monkey,
	}
	modelStyle = lipgloss.NewStyle().
			Align(lipgloss.Center, lipgloss.Center).
			BorderStyle(lipgloss.HiddenBorder())
	focusedModelStyle = lipgloss.NewStyle().
				Align(lipgloss.Center, lipgloss.Center).
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color("69"))
	spinnerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("69"))
	helpStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	listStyle    = lipgloss.NewStyle().Margin(1, 2)
)

type mainModel struct {
	list    list.Model
	state   sessionState
	timer   timer.Model
	spinner spinner.Model
	index   int
	choice  item
}

func newModel(timeout time.Duration, items []list.Item) mainModel {
	m := mainModel{state: timerView}
	m.timer = timer.New(timeout)
	m.spinner = spinner.New()
	m.list = list.New(items, list.NewDefaultDelegate(), 0, 0)
	m.list.Title = "Testtt"
	m.list.SetShowHelp(false)
	m.list.SetShowStatusBar(false)
	m.list.KeyMap.NextPage.Unbind()
	m.list.KeyMap.PrevPage.Unbind()
	return m
}

func (m mainModel) Init() tea.Cmd {
	// start the timer and spinner on program start
	return tea.Batch(m.timer.Init(), m.spinner.Tick, tea.EnterAltScreen)
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "tab":
			if m.state == timerView {
				m.state = spinnerView
			} else {
				m.state = timerView
			}
		case "n":
			if m.state == timerView {
				m.timer = timer.New(defaultTime)
				cmds = append(cmds, m.timer.Init())
			} else {
				m.Next()
				m.resetSpinner()
				cmds = append(cmds, m.spinner.Tick)
			}
		case "enter":
			if m.state == timerView {
				i, ok := m.list.SelectedItem().(item)
				if ok {
					m.choice = i
				}
				return m, tea.Quit
			}
		}
		switch m.state {
		// update whichever model is focused
		case spinnerView:
			m.spinner, cmd = m.spinner.Update(msg)
			cmds = append(cmds, cmd)
		default:
			m.timer, cmd = m.timer.Update(msg)
			cmds = append(cmds, cmd)
		}
	case spinner.TickMsg:
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)
	case timer.TickMsg:
		m.timer, cmd = m.timer.Update(msg)
		cmds = append(cmds, cmd)
	case tea.WindowSizeMsg:
		h, v := listStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	m.list, cmd = m.list.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m mainModel) View() string {
	var s string
	// model := m.currentFocusedModel()
	if m.state == timerView {
		s += lipgloss.JoinHorizontal(lipgloss.Top, focusedModelStyle.Render(m.list.View()), modelStyle.Render(m.spinner.View()))
	} else {
		s += lipgloss.JoinHorizontal(lipgloss.Top, modelStyle.Render(m.list.View()), focusedModelStyle.Render(m.spinner.View()))
	}
	// s += helpStyle.Render(fmt.Sprintf("\ntab: focus next • n: new %s • q: exit\n", model))
	return s
}

func (m mainModel) currentFocusedModel() string {
	if m.state == timerView {
		return "timer"
	}
	return "spinner"
}

func (m *mainModel) Next() {
	if m.index == len(spinners)-1 {
		m.index = 0
	} else {
		m.index++
	}
}

func (m *mainModel) resetSpinner() {
	m.spinner = spinner.New()
	m.spinner.Style = spinnerStyle
	m.spinner.Spinner = spinners[m.index]
}

// Open lamp ghost tui application.
// The terminal would be split into 2 pieces:
// left is the specified difficult table's levels
// right is the related song list and lamp status
func OpenGhostTui(dth *difftable.DiffTableHeader, dt []difftable.DiffTable, songData []SongData, scoreLog []ScoreLog) {
	items := []list.Item{
		item{title: "Raspberry Pi’s", desc: "I have ’em all over my house"},
		item{title: "Nutella", desc: "It's good on toast"},
		item{title: "Bitter melon", desc: "It cools you down"},
		item{title: "Nice socks", desc: "And by that I mean socks without holes"},
		item{title: "Eight hours of sleep", desc: "I had this once"},
		item{title: "Cats", desc: "Usually"},
		item{title: "Plantasia, the album", desc: "My plants love it too"},
		item{title: "Pour over coffee", desc: "It takes forever to make though"},
		item{title: "VR", desc: "Virtual reality...what is there to say?"},
		item{title: "Noguchi Lamps", desc: "Such pleasing organic forms"},
		item{title: "Linux", desc: "Pretty much the best OS"},
		item{title: "Business school", desc: "Just kidding"},
		item{title: "Pottery", desc: "Wet clay is a great feeling"},
		item{title: "Shampoo", desc: "Nothing like clean hair"},
		item{title: "Table tennis", desc: "It’s surprisingly exhausting"},
		item{title: "Milk crates", desc: "Great for packing in your extra stuff"},
		item{title: "Afternoon tea", desc: "Especially the tea sandwich part"},
		item{title: "Stickers", desc: "The thicker the vinyl the better"},
		item{title: "20° Weather", desc: "Celsius, not Fahrenheit"},
		item{title: "Warm light", desc: "Like around 2700 Kelvin"},
		item{title: "The vernal equinox", desc: "The autumnal equinox is pretty good too"},
		item{title: "Gaffer’s tape", desc: "Basically sticky fabric"},
		item{title: "Terrycloth", desc: "In other words, towel fabric"},
	}
	if _, err := tea.NewProgram(newModel(defaultTime, items)).Run(); err != nil {
		log.Fatal(err)
	}
}
