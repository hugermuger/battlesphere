package main

import (
	"log"
	"os"

	"charm.land/bubbles/v2/filepicker"
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
	{title: "Card Search"},
	{title: "Collection"},
	{title: "Decks"},
	{title: "User"},
	{title: "Import"},
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
	jwtToken       string
	refreshToken   string
	username       string
	activeTabIndex int
	winWidth       int
	isTyping       bool
	focusTabs      bool
	searching      bool
	err            error
	search         searchMenu
	login          loginMenu
	listStyles     listStyle
	filepicker     filepickerModel
}

type filepickerModel struct {
	model                filepicker.Model
	selectedFile         string
	focusIndex           int
	collectionFilepicker bool
	collectionImport     bool
	err                  string
	status               string
	cue                  string
	progress             int
	missing              [][]string
}

type listStyle struct {
	selectedDelegate   list.DefaultDelegate
	unselectedDelegate list.DefaultDelegate
}

type searchMenu struct {
	searchID           int
	searchQuery        string
	searchQueryLang    int
	selectedLang       int
	menuSearchIndex    int
	oracleCardID       uuid.UUID
	oracleCard         types.OracleDetail
	rulingViewport     viewport.Model
	searchViewport     viewport.Model
	listSearchByName   list.Model
	listSearchByOracle list.Model
	focusInput         bool
	focusList          bool
	focusViewport      bool
	searchInput        textinput.Model
}

type loginMenu struct {
	focusedStyle   lipgloss.Style
	blurredStyle   lipgloss.Style
	registerView   bool
	logoutView     bool
	registerInput  []textinput.Model
	loginInput     []textinput.Model
	focusIndex     int
	loggedIn       bool
	message        string
	err            string
	registerSucces bool
}

func initModel() model {
	// Login Menu
	focusedStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("208"))
	blurredStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("241"))

	registerInput := make([]textinput.Model, 4)

	var t textinput.Model
	for i := range registerInput {
		t = textinput.New()
		t.CharLimit = 64

		s := t.Styles()
		s.Cursor.Color = lipgloss.Color("208")
		s.Focused.Prompt = focusedStyle
		s.Focused.Text = focusedStyle
		s.Blurred.Prompt = blurredStyle
		s.Focused.Text = focusedStyle
		t.SetStyles(s)
		t.SetWidth(64)

		switch i {
		case 0:
			t.Placeholder = "Username"
		case 1:
			t.Placeholder = "Email"
		case 2:
			t.Placeholder = "Password"
			t.EchoMode = textinput.EchoPassword
			t.EchoCharacter = '*'
		case 3:
			t.Placeholder = "Repeat Password"
			t.EchoMode = textinput.EchoPassword
			t.EchoCharacter = '*'

		}

		registerInput[i] = t
	}

	loginInput := make([]textinput.Model, 2)

	for i := range loginInput {
		t = textinput.New()
		t.CharLimit = 64

		s := t.Styles()
		s.Cursor.Color = lipgloss.Color("208")
		s.Focused.Prompt = focusedStyle
		s.Focused.Text = focusedStyle
		s.Blurred.Prompt = blurredStyle
		s.Focused.Text = focusedStyle
		t.SetStyles(s)
		t.SetWidth(64)

		switch i {
		case 0:
			t.Placeholder = "Username"
		case 1:
			t.Placeholder = "Password"
			t.EchoMode = textinput.EchoPassword
			t.EchoCharacter = '*'
		}

		loginInput[i] = t
	}

	// Search Menu
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

		activeTabIndex: 0,
		winWidth:       200,
		isTyping:       false,
		focusTabs:      true,
		searching:      true,
		login: loginMenu{
			blurredStyle:   blurredStyle,
			focusedStyle:   focusedStyle,
			registerInput:  registerInput,
			loginInput:     loginInput,
			focusIndex:     0,
			registerView:   false,
			logoutView:     true,
			registerSucces: false,
		},
		search: searchMenu{
			listSearchByName:   listSearchByName,
			listSearchByOracle: listSearchByOracle,
			searchInput:        searchInput,
			rulingViewport:     rulingViewport,
			searchViewport:     searchViewport,
			searchQuery:        "",
			searchQueryLang:    len(lang),
			selectedLang:       0,
			menuSearchIndex:    0,
			focusInput:         false,
			focusList:          false,
			focusViewport:      false,
		},
		listStyles: listStyle{
			selectedDelegate:   selectedDelegate,
			unselectedDelegate: unselectedDelegate,
		},
		filepicker: filepickerModel{
			focusIndex:           0,
			collectionFilepicker: false,
			collectionImport:     false,
		},
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, m.filepicker.model.Init())
}

func main() {
	fp := filepicker.New()
	fp.AllowedTypes = []string{".csv"}
	fp.AutoHeight = false
	home, err := os.UserHomeDir()
	if err != nil {
		fp.CurrentDirectory = "."
	} else {
		fp.CurrentDirectory = home
	}
	m := initModel()
	m.filepicker.model = fp

	config, err := loadUserConfig()
	if err == nil && config.LastToken != "" {
		m.refreshToken = config.LastToken
		handlerRefresh(config.LastToken, &m)
	}
	if _, err := tea.NewProgram(m).Run(); err != nil {
		log.Fatalf("Error running program: %v", err)
	}
}
