package main

import (
	"log"
	"strings"

	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

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

type model struct {
	searchInput    textinput.Model
	activeTabIndex int
	winWidth       int
	isTyping       bool
	list           list.Model
	focusInput     bool
	searching      bool
}

func initModel() model {
	si := textinput.New()
	si.Placeholder = "Search for a card..."
	si.SetWidth(30)

	l := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	l.SetShowTitle(false)
	l.SetShowHelp(false)

	return model{
		list:           l,
		searchInput:    si,
		activeTabIndex: 0,
		winWidth:       200,
		isTyping:       false,
		focusInput:     false,
		searching:      true,
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

	case tea.KeyPressMsg:
		switch m.activeTabIndex {
		case 1:
			switch msg.String() {
			case "enter":
				if !m.focusInput {
					m.focusInput = true
					m.searchInput.Focus()
					return m, nil
				} else if m.focusInput {
					m.focusInput = false
					m.searchInput.Blur()
					return m, nil
				}
			case "c":
				if !m.focusInput {
					m.searchInput.Reset()
				}
			}

			if m.focusInput {
				m.searchInput, cmd = m.searchInput.Update(msg)
			}
		}

		if !m.focusInput {
			switch msg.String() {
			case "l", "right":
				m.activeTabIndex++
				if m.activeTabIndex >= len(tuiTabs) {
					m.activeTabIndex = 0
				}
			case "h", "left":
				m.activeTabIndex--
				if m.activeTabIndex < 0 {
					m.activeTabIndex = len(tuiTabs) - 1
				}
			}
		}

		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, cmd
}

func (m model) View() tea.View {

	tabsRow := tabs(m)

	title := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#fff")).
		Render(tuiTabs[m.activeTabIndex].title)

	helpText := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Render("Navigate with <-/H and L/-> between tabs")

	selected := ""
	switch m.activeTabIndex {
	case 1:
		selected = searchView(m)
	default:
		selected = stdView(m)
	}

	v := tea.NewView(lipgloss.JoinVertical(lipgloss.Center,
		tabsRow,
		title,
		"\n"+selected,
		strings.Repeat("\n", 5),
		helpText))

	v.AltScreen = true
	v.WindowTitle = "Battlesphere"

	return v
}

func stdView(m model) string {
	return "NOT INCLUDED YET: " + tuiTabs[m.activeTabIndex].title
}

func searchView(m model) string {
	if len(m.list.Items()) == 0 {
		return m.searchInput.View()
	}
	return m.searchInput.View() + "\n" + m.list.View()
}

func tabs(m model) string {
	var s []string
	for i, tab := range tuiTabs {
		if i == m.activeTabIndex {
			s = append(s, tabStyle.Border(activeTabBorder).BorderForeground(lipgloss.Color("183")).Render(tab.title))
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

func main() {
	if _, err := tea.NewProgram(initModel()).Run(); err != nil {
		log.Fatalf("Error running program: %v", err)
	}
}
