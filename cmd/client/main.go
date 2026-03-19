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

const website = "http://localhost:8080"

type debounceMsg struct {
	id     int
	query  string
	langID int
}

type searchResultMsg []list.Item

type oracleSearchResultMsg struct {
	list       []list.Item
	oracleCard types.OracleDetail
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
	searchInput        textinput.Model
	activeTabIndex     int
	winWidth           int
	isTyping           bool
	listSearchByName   list.Model
	listSearchByOracle list.Model
	focusInput         bool
	focusList          bool
	focusViewport      bool
	focusTabs          bool
	searching          bool
	err                error
	searchID           int
	searchQuery        string
	searchQueryLang    int
	selectedLang       int
	menuSearchID       int
	oracleCardID       uuid.UUID
	oracleCard         types.OracleDetail
	rulingViewport     viewport.Model
	searchViewport     viewport.Model
	selectedDelegate   list.DefaultDelegate
	unselectedDelegate list.DefaultDelegate
}

func initModel() model {
	searchInput := textinput.New()
	searchInput.Placeholder = "Search for a card..."
	searchInput.SetWidth(30)

	selectedDelegate := list.NewDefaultDelegate()
	selectedDelegate.Styles.SelectedTitle = selectedDelegate.Styles.SelectedTitle.
		Foreground(lipgloss.Color("117")).
		BorderForeground(lipgloss.Color("208"))
	selectedDelegate.Styles.SelectedDesc = selectedDelegate.Styles.SelectedDesc.
		Foreground(lipgloss.Color("152")).
		BorderForeground(lipgloss.Color("208"))

	unselectedDelegate := list.NewDefaultDelegate()
	col := unselectedDelegate.Styles.NormalDesc.GetForeground()
	unselectedDelegate.Styles.SelectedTitle = unselectedDelegate.Styles.SelectedTitle.
		Foreground(lipgloss.Color("252")).
		BorderForeground(col)
	unselectedDelegate.Styles.SelectedDesc = unselectedDelegate.Styles.SelectedDesc.
		Foreground(col).
		BorderForeground(col)

	listSearchByName := list.New([]list.Item{}, unselectedDelegate, 0, 0)
	listSearchByName.SetShowTitle(false)
	listSearchByName.SetShowHelp(false)
	listSearchByName.DisableQuitKeybindings()
	listSearchByName.SetStatusBarItemName("card", "cards")

	rulingViewport := viewport.New()
	rulingViewport.SetWidth(30)

	searchViewport := viewport.New()
	searchViewport.SetWidth(30)

	listSearchByOracle := list.New([]list.Item{}, unselectedDelegate, 0, 0)
	listSearchByOracle.SetShowTitle(false)
	listSearchByOracle.SetShowHelp(false)
	listSearchByOracle.DisableQuitKeybindings()
	listSearchByOracle.SetStatusBarItemName("print", "prints")

	return model{
		listSearchByName:   listSearchByName,
		listSearchByOracle: listSearchByOracle,
		searchInput:        searchInput,
		rulingViewport:     rulingViewport,
		searchViewport:     searchViewport,
		activeTabIndex:     0,
		winWidth:           200,
		isTyping:           false,
		focusInput:         false,
		focusList:          false,
		focusViewport:      false,
		focusTabs:          true,
		searching:          true,
		searchQuery:        "",
		searchQueryLang:    len(lang),
		selectedLang:       0,
		menuSearchID:       0,
		selectedDelegate:   selectedDelegate,
		unselectedDelegate: unselectedDelegate,
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
