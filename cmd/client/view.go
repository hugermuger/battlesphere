package main

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

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
	case "User":
		selected = loginView(m)
	case "Card Search":
		selected = searchView(m)
	case "Import":
		selected = importView(m)
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
