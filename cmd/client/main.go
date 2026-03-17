package main

import (
	"log"

	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/textinput"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/google/uuid"
	"github.com/hugermuger/battlesphere/internal/types"
)

type debounceMsg struct {
	id     int
	query  string
	langID int
}

type searchResultMsg []list.Item

type oracleSearchResultMsg struct {
	list       []list.Item
	oracleCard types.ResponseByOracleID
	legalities types.Legalities
}

var tabBorder = lipgloss.Border{
	Top:         "─",
	Bottom:      "─",
	Left:        "│",
	Right:       "│",
	TopLeft:     "╭",
	TopRight:    "╮",
	BottomLeft:  "┴",
	BottomRight: "┴",
}

var activeTabBorder = lipgloss.Border{
	Top:         "─",
	Bottom:      " ",
	Left:        "│",
	Right:       "│",
	TopLeft:     "╭",
	TopRight:    "╮",
	BottomLeft:  "┘",
	BottomRight: "└",
}

var tabStyle = lipgloss.NewStyle().
	Border(tabBorder).
	PaddingLeft(1).
	PaddingRight(1)

type TUITab struct {
	title string
}

var tuiTabs = []TUITab{
	{title: "Welcome"},
	{title: "Card Search"},
	{title: "Collection"},
}

var menuSearch = []string{
	"name",
	"oracle",
}

var lang = []string{
	"en",
	"de",
	"fr",
	"es",
	"it",
	"pt",
	"ja",
	"ko",
	"ru",
	"zhs",
	"zht",
}

type model struct {
	searchInput            textinput.Model
	activeTabIndex         int
	winWidth               int
	isTyping               bool
	listSearchByName       list.Model
	listSearchByOracle     list.Model
	focusInput             bool
	focusList              bool
	focusViewport          bool
	focusTabs              bool
	searching              bool
	err                    error
	searchID               int
	searchQuery            string
	searchQueryLang        int
	selectedLang           int
	menuSearchID           int
	oracleCardID           uuid.UUID
	oracleCard             types.ResponseByOracleID
	rulingViewport         viewport.Model
	searchViewport         viewport.Model
	selectedListDelegate   list.DefaultDelegate
	unselectedListDelegate list.DefaultDelegate
}

func initModel() model {
	si := textinput.New()
	si.Placeholder = "Search for a card..."
	si.SetWidth(30)

	d := list.NewDefaultDelegate()
	d.Styles.SelectedTitle = d.Styles.SelectedTitle.
		Foreground(lipgloss.Color("117")).
		BorderForeground(lipgloss.Color("208"))
	d.Styles.SelectedDesc = d.Styles.SelectedDesc.
		Foreground(lipgloss.Color("152")).
		BorderForeground(lipgloss.Color("208"))

	ud := list.NewDefaultDelegate()
	col := ud.Styles.NormalDesc.GetForeground()
	ud.Styles.SelectedTitle = ud.Styles.SelectedTitle.
		Foreground(lipgloss.Color("252")).
		BorderForeground(col)
	ud.Styles.SelectedDesc = ud.Styles.SelectedDesc.
		Foreground(col).
		BorderForeground(col)

	l := list.New([]list.Item{}, ud, 0, 0)
	l.SetShowTitle(false)
	l.SetShowHelp(false)
	l.DisableQuitKeybindings()
	l.SetStatusBarItemName("card", "cards")

	vpr := viewport.New()
	vpr.SetWidth(30)

	vps := viewport.New()
	vps.SetWidth(30)

	listOracle := list.New([]list.Item{}, ud, 0, 0)
	listOracle.SetShowTitle(false)
	listOracle.SetShowHelp(false)
	listOracle.DisableQuitKeybindings()
	listOracle.SetStatusBarItemName("print", "prints")

	return model{
		listSearchByName:       l,
		listSearchByOracle:     listOracle,
		searchInput:            si,
		rulingViewport:         vpr,
		searchViewport:         vps,
		activeTabIndex:         0,
		winWidth:               200,
		isTyping:               false,
		focusInput:             false,
		focusList:              false,
		focusViewport:          false,
		focusTabs:              true,
		searching:              true,
		searchQuery:            "",
		searchQueryLang:        len(lang),
		selectedLang:           0,
		menuSearchID:           0,
		selectedListDelegate:   d,
		unselectedListDelegate: ud,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func main() {
	if _, err := tea.NewProgram(initModel()).Run(); err != nil {
		log.Fatalf("Error running program: %v", err)
	}
}
