package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/textinput"
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
	"cards",
}

var lang = []string{
	"en",
	"de",
	"fr",
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
	searching          bool
	err                error
	searchID           int
	searchQuery        string
	searchQueryLang    int
	selectedLang       int
	menuSearchID       int
	oracleCardID       uuid.UUID
}

func initModel() model {
	si := textinput.New()
	si.Placeholder = "Search for a card..."
	si.SetWidth(30)

	l := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	l.SetShowTitle(false)
	l.SetShowHelp(false)
	l.DisableQuitKeybindings()

	listOracle := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	listOracle.SetShowTitle(false)
	listOracle.SetShowHelp(false)
	listOracle.DisableQuitKeybindings()

	return model{
		listSearchByName:   l,
		listSearchByOracle: listOracle,
		searchInput:        si,
		activeTabIndex:     0,
		winWidth:           200,
		isTyping:           false,
		focusInput:         false,
		focusList:          false,
		searching:          true,
		searchQuery:        "",
		searchQueryLang:    len(lang),
		selectedLang:       0,
		menuSearchID:       0,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.winWidth = msg.Width
		m.listSearchByName.SetSize(msg.Width, msg.Height-10)
		m.listSearchByOracle.SetSize(msg.Width, msg.Height-8)

	case tea.KeyPressMsg:
		if !m.focusInput && !m.focusList {
			switch msg.String() {
			case "right":
				m.activeTabIndex++
				if m.activeTabIndex >= len(tuiTabs) {
					m.activeTabIndex = 0
				}
			case "left":
				m.activeTabIndex--
				if m.activeTabIndex < 0 {
					m.activeTabIndex = len(tuiTabs) - 1
				}
			}
		}

		switch msg.String() {
		case "ctrl+q":
			return m, tea.Quit
		}

		switch tuiTabs[m.activeTabIndex].title {
		case "Card Search":
			switch msg.String() {
			case "enter":
				switch m.menuSearchID {
				case 0:
					if m.focusInput {
						m.focusList = true
						m.focusInput = false
						m.searchInput.Blur()
						return m, nil
					} else if !m.focusInput && !m.focusList {
						m.focusInput = true
						m.focusList = false
						m.searchInput.Focus()
						return m, nil
					} else if m.focusList {
						selectedItem := m.listSearchByName.SelectedItem()
						if selectedItem != nil {
							card := selectedItem.(types.CardResponseSearchByName)
							if card.OracleID != nil {
								m.searching = true
								m.menuSearchID = 1
								m.oracleCardID = *card.OracleID
								return m, m.fetchCardsByOracleID(card.OracleID.String())
							}
						}
					}
				case 1:
					if !m.focusList {
						m.focusList = true
					}
				}

			case "down":
				if m.focusInput {
					m.focusList = true
					m.focusInput = false
					m.searchInput.Blur()
					return m, nil
				}

			case "ctrl+c":
				m.searchInput.Reset()
				m.searchQuery = ""
				m.listSearchByName.SetItems([]list.Item{})
				m.listSearchByOracle.SetItems([]list.Item{})
				m.menuSearchID = 0
				if m.focusList {
					m.focusList = false
					m.focusInput = true
					m.searchInput.Focus()
					return m, nil
				}

			case "esc":
				if m.focusInput {
					m.focusInput = false
					m.searchInput.Blur()
				}
				m.focusList = false
				return m, nil

			case "ctrl+l":
				m.selectedLang++
				if m.selectedLang > len(lang)-1 {
					m.selectedLang = 0
				}
				if m.menuSearchID == 1 {
					return m, tea.Batch(cmd, m.fetchCardsByOracleID(m.oracleCardID.String()))
				}

			case "backspace":
				switch m.menuSearchID {
				case 0:
					if m.focusList {
						m.focusInput = true
						m.focusList = false
						m.searchInput.Focus()
						return m, nil
					}
				case 1:
					m.menuSearchID = 0
					m.selectedLang = m.searchQueryLang
					return m, nil
				}
			}

			if m.focusInput {
				m.searchInput, cmd = m.searchInput.Update(msg)
			} else if m.focusList {
				switch m.menuSearchID {
				case 0:
					m.listSearchByName, cmd = m.listSearchByName.Update(msg)
					return m, nil
				case 1:
					m.listSearchByOracle, cmd = m.listSearchByOracle.Update(msg)
					return m, nil
				}

			}

			if m.menuSearchID == 0 {
				if m.searchInput.Value() != m.searchQuery || m.searchQueryLang != m.selectedLang {
					m.searchQuery = m.searchInput.Value()
					m.searchQueryLang = m.selectedLang
					m.searchID++
					return m, tea.Batch(cmd, m.debounceSearch(m.searchID, m.searchInput.Value(), m.selectedLang))
				}
			}

			return m, cmd
		}

	case debounceMsg:
		switch tuiTabs[m.activeTabIndex].title {
		case "Card Search":
			if msg.id == m.searchID && msg.query != "" {
				m.searching = true
				return m, m.fetchCardsByName(msg.query)
			} else if msg.query == "" {
				m.listSearchByName.SetItems([]list.Item{})
			}
			return m, nil
		}

	case searchResultMsg:
		switch tuiTabs[m.activeTabIndex].title {
		case "Card Search":
			m.searching = false
			switch m.menuSearchID {
			case 0:
				m.listSearchByName.Select(0)
				m.listSearchByName.SetItems(msg)
			case 1:
				m.listSearchByOracle.Select(0)
				m.listSearchByOracle.SetItems(msg)
			}
			return m, nil
		}

	case error:
		m.err = msg
		return m, nil
	}

	return m, cmd
}

func (m model) View() tea.View {

	tabsRow := tabs(m)

	titleText := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#fff")).
		Render(tuiTabs[m.activeTabIndex].title)
	title := lipgloss.PlaceHorizontal(m.winWidth, lipgloss.Center, titleText)

	helpText := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Render("[ctrl+l] select language, [ctrl+c] reset search, [ctrl+q] exit")
	help := lipgloss.PlaceHorizontal(m.winWidth, lipgloss.Center, helpText)

	selected := ""
	switch tuiTabs[m.activeTabIndex].title {
	case "Card Search":
		selected = searchView(m)
	default:
		selected = stdView(m)
	}

	v := tea.NewView(lipgloss.JoinVertical(lipgloss.Left,
		tabsRow,
		title,
		"\n"+selected,
		"\n"+help))

	v.AltScreen = true
	v.WindowTitle = "Battlesphere"

	return v
}

func stdView(m model) string {
	return lipgloss.PlaceHorizontal(m.winWidth, lipgloss.Center, "NOT INCLUDED YET: "+tuiTabs[m.activeTabIndex].title)
}

func searchView(m model) string {
	la := ""
	for i, l := range lang {
		if i == m.selectedLang {
			la = fmt.Sprintf("%v [%v]", la, l)
		} else {
			la += " " + l
		}
	}

	la = lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Render(la)

	v := ""

	switch m.menuSearchID {
	case 0:
		if len(m.listSearchByName.Items()) == 0 {
			v = lipgloss.PlaceHorizontal(m.winWidth, lipgloss.Center, m.searchInput.View()) + "\n" + "\n" + lipgloss.PlaceHorizontal(m.winWidth, lipgloss.Center, la)
		} else {
			v = lipgloss.PlaceHorizontal(m.winWidth, lipgloss.Center, m.searchInput.View()) + "\n" + "\n" + m.listSearchByName.View() + "\n" + lipgloss.PlaceHorizontal(m.winWidth, lipgloss.Center, la)
		}
	case 1:
		v = m.listSearchByOracle.View() + "\n" + lipgloss.PlaceHorizontal(m.winWidth, lipgloss.Center, la)
	}
	return v
}

func tabs(m model) string {
	var s []string
	for i, tab := range tuiTabs {
		if i == m.activeTabIndex {
			s = append(s, tabStyle.Border(activeTabBorder).Render(tab.title))
		} else {
			s = append(s, tabStyle.Render(tab.title))
		}
	}

	tabs := lipgloss.JoinHorizontal(lipgloss.Bottom, s...)
	spaceLeft := m.winWidth - lipgloss.Width(tabs) - 2

	spacer := lipgloss.NewStyle().
		Border(lipgloss.Border{Bottom: "─"}, false, false, true, false).
		Render(strings.Repeat(" ", spaceLeft))

	return lipgloss.JoinHorizontal(lipgloss.Bottom, tabs, spacer)
}

func (m model) debounceSearch(id int, query string, langID int) tea.Cmd {
	return func() tea.Msg {
		time.Sleep(300 * time.Millisecond)
		return debounceMsg{id: id, query: query, langID: langID}
	}
}

func (m model) fetchCardsByName(query string) tea.Cmd {
	return func() tea.Msg {
		endpoint := "http://localhost:8080/cards/search?name=" + url.QueryEscape(query) + "&limit=600" + "&lang=" + lang[m.selectedLang]
		resp, err := http.Get(endpoint)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		var data struct {
			Results []types.CardResponseSearchByName `json:"results"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return err
		}

		var items []list.Item
		for _, v := range data.Results {
			items = append(items, v)
		}

		return searchResultMsg(items)
	}
}

func (m model) fetchCardsByOracleID(oracleID string) tea.Cmd {
	return func() tea.Msg {
		endpoint := "http://localhost:8080/cards/oracle/" + oracleID + "?&limit=600" + "&lang=" + lang[m.selectedLang]
		resp, err := http.Get(endpoint)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		var data struct {
			Results []types.CardResponseSearchByOracleID `json:"results"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return err
		}

		var items []list.Item
		for _, v := range data.Results {
			items = append(items, v)
		}
		return searchResultMsg(items)
	}
}

func main() {
	if _, err := tea.NewProgram(initModel()).Run(); err != nil {
		log.Fatalf("Error running program: %v", err)
	}
}
