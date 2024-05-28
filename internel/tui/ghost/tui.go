package ghost

import (
	"log"
	"sort"
	"strconv"
	"time"

	"github.com/Catizard/lampghost/internel/difftable"
	"github.com/Catizard/lampghost/internel/rival"
	"github.com/Catizard/lampghost/internel/score"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type item struct {
	title, desc, level string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

// sessionState is used to track which model is focused
type sessionState uint

const (
	defaultTime              = time.Minute
	levelView   sessionState = iota
	songView
)

var (
	listStyle = lipgloss.NewStyle().Margin(1, 2)
	flBlock   = lipgloss.NewStyle().SetString(" ").Width(10).Background(lipgloss.Color("#FF0000"))
	ezBlock   = lipgloss.NewStyle().SetString(" ").Width(10).Background(lipgloss.Color("#00FF00"))
	nrBlock   = lipgloss.NewStyle().SetString(" ").Width(10).Background(lipgloss.Color("#0000FF"))
	hcBlock   = lipgloss.NewStyle().SetString(" ").Width(10).Background(lipgloss.Color("#FFFFFF"))
)

type mainModel struct {
	levelList   list.Model
	songList    list.Model
	levelData   []string // not used
	state       sessionState
	w           int
	h           int
	songDataMap map[string][]difftable.DiffTable
}

func buildLevelList(dth *difftable.DiffTableHeader, diffTable []difftable.DiffTable) ([]list.Item, string) {
	// convert diffTable to list items
	levels := make(map[string]interface{})
	for _, v := range diffTable {
		levels[v.Level] = new(interface{})
	}
	if len(levels) == 0 {
		panic("tableHeader.json file corrupted, no level found")
	}
	sortedLevels := make([]string, 0)
	for level := range levels {
		sortedLevels = append(sortedLevels, level)
	}
	sort.Slice(sortedLevels, func(i, j int) bool {
		ll := sortedLevels[i]
		rr := sortedLevels[j]
		ill, errL := strconv.Atoi(ll)
		irr, errR := strconv.Atoi(rr)
		if errL == nil && errR == nil {
			return ill < irr
		}
		return ll < rr
	})

	items := make([]list.Item, 0)
	for _, v := range sortedLevels {
		title := dth.Symbol + " " + v
		n := item{
			title: title,
			desc:  "発狂BMS難易度表" + title,
			level: v,
		}
		items = append(items, n)
	}
	return items, sortedLevels[0]
}

func buildSongList(m *mainModel, diffTable []difftable.DiffTable) {
	// Level maps to song info array
	m.songDataMap = make(map[string][]difftable.DiffTable)
	for _, v := range diffTable {
		if _, ok := m.songDataMap[v.Level]; !ok {
			m.songDataMap[v.Level] = make([]difftable.DiffTable, 0)
		}
		m.songDataMap[v.Level] = append(m.songDataMap[v.Level], v)
	}
}

func drawLamp(lamp int32) string {
	if lamp == 0 {
		return ""
	}
	if 1 <= lamp && lamp < 4 {
		return flBlock.String()
	} else if lamp < 5 {
		return ezBlock.String()
	} else if lamp < 6 {
		return nrBlock.String()
	} else if lamp < 11 {
		return hcBlock.String()
	}
	return ""
}

func (m *mainModel) transferLevel(level string) {
	rawArray := m.songDataMap[level]
	itemList := make([]list.Item, 0)
	for _, v := range rawArray {
		selfLamp := lipgloss.NewStyle().MarginRight(20).Render(drawLamp(v.Lamp))
		ghostLamp := drawLamp(v.GhostLamp)
		desc := lipgloss.JoinHorizontal(0, selfLamp, ghostLamp)
		n := item{
			title: v.Title,
			desc:  desc,
		}
		itemList = append(itemList, n)
	}
	songList := list.New(itemList, list.NewDefaultDelegate(), 0, 0)
	songList.DisableQuitKeybindings()
	// At the very begining, w & h was not set
	if m.w != 0 && m.h != 0 {
		h, v := listStyle.GetFrameSize()
		songList.SetSize(m.w-h, m.h-v)
	}
	m.songList = songList
}

func newModel(dth *difftable.DiffTableHeader, diffTable []difftable.DiffTable) mainModel {
	m := mainModel{state: levelView}
	// Build level list
	levelItems, defaultLevel := buildLevelList(dth, diffTable)
	m.levelList = list.New(levelItems, list.NewDefaultDelegate(), 0, 0)
	m.levelList.Title = "Levels"
	m.levelList.SetShowHelp(false)
	m.levelList.SetShowStatusBar(false)
	m.levelList.KeyMap.NextPage.Unbind()
	m.levelList.KeyMap.PrevPage.Unbind()
	m.levelList.DisableQuitKeybindings()
	// Build song list
	buildSongList(&m, diffTable)
	m.transferLevel(defaultLevel)
	return m
}

func (m mainModel) Init() tea.Cmd {
	return tea.Batch(tea.EnterAltScreen, tea.DisableMouse)
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.state == songView {
				m.state = levelView
			}
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			if m.state == levelView {
				i, ok := m.levelList.SelectedItem().(item)
				if ok {
					m.transferLevel(i.level)
				}
				m.state = songView
			}
		}
		switch m.state {
		// update whichever model is focused
		case levelView:
			m.levelList, cmd = m.levelList.Update(msg)
			cmds = append(cmds, cmd)
		default:
			m.songList, cmd = m.songList.Update(msg)
			cmds = append(cmds, cmd)
		}
	case tea.WindowSizeMsg:
		h, v := listStyle.GetFrameSize()
		m.w = msg.Width
		m.h = msg.Height
		m.levelList.SetSize(msg.Width-h, msg.Height-v)
		m.songList.SetSize(msg.Width-h, msg.Height-v)
	}

	return m, tea.Batch(cmds...)
}

func (m mainModel) View() string {
	var s string
	if m.state == levelView {
		s += lipgloss.JoinHorizontal(lipgloss.Top, m.levelList.View())
	} else {
		s += lipgloss.JoinHorizontal(lipgloss.Top, m.levelList.View(), m.songList.View())
	}
	return s
}

// Open lamp ghost tui application.
// The terminal would be split into 2 pieces:
// left is the specified difficult table's levels
// right is the related song list and lamp status
func OpenGhostTui(dth *difftable.DiffTableHeader, dt []difftable.DiffTable, selfInfo *rival.RivalInfo, ghostInfo *rival.RivalInfo) {
	// NOTE: merge songData -> diffTable, scoreLog -> diffTable before any operation
	// Merge self
	mergeSha256FromSongData(dt, selfInfo.SongData)
	mergeLampFromScoreLog(dt, selfInfo.ScoreLog)
	// Merge ghost
	mergeSha256FromSongData(dt, ghostInfo.SongData)
	mergeGhostLampFromScoreLog(dt, ghostInfo.ScoreLog)
	// After two merge functions, dt now contains lamp info
	if _, err := tea.NewProgram(newModel(dth, dt)).Run(); err != nil {
		log.Fatal(err)
	}
}

// Merge Sha256 field from song data
// In place function, do not return a new array
func mergeSha256FromSongData(dtArray []difftable.DiffTable, songData []score.SongData) {
	songDataMd5Map := make(map[string]score.SongData)
	for _, v := range songData {
		songDataMd5Map[v.Md5] = v
	}
	for i, dt := range dtArray {
		if songData, ok := songDataMd5Map[dt.Md5]; ok {
			dtArray[i].Sha256 = songData.Sha256
		}
	}
}

// Merge maximum lamp from scorelog
// In place function, do not return a new array
func mergeLampFromScoreLog(dtArray []difftable.DiffTable, scoreLog []score.ScoreLog) {
	dtSha256Map := make(map[string]*difftable.DiffTable)
	for i, v := range dtArray {
		dtSha256Map[v.Sha256] = &dtArray[i]
	}
	for _, v := range scoreLog {
		if t, ok := dtSha256Map[v.Sha256]; ok {
			t.Lamp = max(t.Lamp, v.Clear)
		}
	}
}

// Same with above function, the only difference is target
func mergeGhostLampFromScoreLog(dtArray []difftable.DiffTable, scoreLog []score.ScoreLog) {
	dtSha256Map := make(map[string]*difftable.DiffTable)
	for i, v := range dtArray {
		dtSha256Map[v.Sha256] = &dtArray[i]
	}
	for _, v := range scoreLog {
		if t, ok := dtSha256Map[v.Sha256]; ok {
			t.GhostLamp = max(t.GhostLamp, v.Clear)
		}
	}
}