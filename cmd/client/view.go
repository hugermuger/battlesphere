package main

import (
	"fmt"
	"reflect"
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
	if m.focusList {
		m.listSearchByName.SetDelegate(m.selectedDelegate)
		m.listSearchByOracle.SetDelegate(m.selectedDelegate)
	} else {
		m.listSearchByName.SetDelegate(m.unselectedDelegate)
		m.listSearchByOracle.SetDelegate(m.unselectedDelegate)
	}

	rulings := "\nRulings"
	if m.focusViewport {
		rulings = (lipgloss.NewStyle().
			Foreground(lipgloss.Color("117")).
			Render(rulings))
	} else {
		rulings = (lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			Render(rulings))
	}

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
		buffer := ""
		if m.oracleCard.ManaCost != nil {
			buffer += *m.oracleCard.ManaCost + "\n"
		}
		buffer += m.oracleCard.Name + "\n" + m.oracleCard.TypeLine
		if m.oracleCard.Power != nil && m.oracleCard.Toughness != nil {
			buffer += "\n" + *m.oracleCard.Power + "/" + *m.oracleCard.Toughness
		}
		if m.oracleCard.Defense != nil {
			buffer += "\n" + *m.oracleCard.Defense
		}
		if m.oracleCard.Loyalty != nil {
			buffer += "\n" + *m.oracleCard.Loyalty
		}
		if m.oracleCard.OracleText != nil {
			buffer += "\n\n" + *m.oracleCard.OracleText
		}
		if m.oracleCard.PrintedText != nil {
			buffer += "\n\n" + *m.oracleCard.PrintedText
		}
		buffer = fmt.Sprintf("%v\n\nCmc: %v", buffer, m.oracleCard.Cmc)
		if len(m.oracleCard.ColorIdentity) > 0 {
			buffer += ", Color Identity: "
			for i := len(m.oracleCard.ColorIdentity) - 1; i >= 0; i-- {
				buffer += fmt.Sprintf("{%v}", m.oracleCard.ColorIdentity[i])
			}
		}
		if m.oracleCard.EdhrecRank != nil {
			buffer = fmt.Sprintf("%v\n\nEDHREC Rank: %v", buffer, *m.oracleCard.EdhrecRank)
		}
		if m.oracleCard.GameChanger != nil {
			if *m.oracleCard.GameChanger {
				buffer += "\n\n" + lipgloss.NewStyle().Foreground(lipgloss.Color("220")).Render("GAMECHANGER")
			}
		}
		buffer = lipgloss.Wrap(buffer, m.winWidth-15, "\n")
		buffer = lipgloss.PlaceHorizontal(m.winWidth, lipgloss.Center, buffer)
		card := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Width(m.winWidth).Render(buffer)
		legalities := renderLegalities(m)
		m.rulingViewport.SetWidth(m.winWidth / 2)
		m.searchViewport.SetWidth(m.winWidth / 2)
		if m.oracleCard.Multifaced {
			faces := renderMultifacesOracle(m)
			height := m.listSearchByName.Height() - (lipgloss.Height(buffer) + lipgloss.Height(faces) + lipgloss.Height(legalities))
			m.rulingViewport.SetHeight(height - 5)
			m.listSearchByOracle.SetHeight(height - 3)
			m.searchViewport.SetContent(m.listSearchByOracle.View())
			m.searchViewport.SetHeight(height - 2)
			listView := lipgloss.JoinHorizontal(lipgloss.Top, m.searchViewport.View(), rulings+"\n\n"+m.rulingViewport.View())
			v = card + "\n" + legalities + "\n" + faces + "\n" + listView + "\n\n" + lipgloss.PlaceHorizontal(m.winWidth, lipgloss.Center, la)
		} else {
			height := m.listSearchByName.Height() - (lipgloss.Height(buffer) + lipgloss.Height(legalities))
			m.rulingViewport.SetHeight(height - 5)
			m.listSearchByOracle.SetHeight(height - 3)
			m.searchViewport.SetContent(m.listSearchByOracle.View())
			m.searchViewport.SetHeight(height - 2)
			listView := lipgloss.JoinHorizontal(lipgloss.Top, m.searchViewport.View(), rulings+"\n\n"+m.rulingViewport.View())
			v = card + "\n" + legalities + "\n" + listView + "\n\n" + lipgloss.PlaceHorizontal(m.winWidth, lipgloss.Center, la)
		}
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

func renderLegalities(m model) string {
	legalities := ""
	l := m.oracleCard.Legalities
	v := reflect.ValueOf(l)
	t := reflect.TypeOf(l)

	legalStyle := lipgloss.NewStyle().Foreground(lipgloss.Green)
	notLegalStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	restrictedStyle := lipgloss.NewStyle().Foreground(lipgloss.Yellow)

	for i := 0; i < v.NumField(); i++ {
		fieldValue := v.Field(i)
		fieldName := t.Field(i).Name

		if fieldValue.IsNil() {
			continue
		}

		status := *fieldValue.Interface().(*string)

		if status == "legal" {
			legalities += legalStyle.Render("[" + fieldName + "] ")
		} else if status == "restricted" {
			legalities += restrictedStyle.Render("[" + fieldName + "] ")
		} else {
			legalities += notLegalStyle.Render("[" + fieldName + "] ")
		}
	}

	legalities = lipgloss.Wrap(legalities, m.winWidth-15, "]")
	legalities = lipgloss.PlaceHorizontal(m.winWidth, lipgloss.Center, legalities)
	legalities = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Width(m.winWidth).Render(legalities)

	return legalities
}

func renderRulings(m model) string {
	rulings := ""
	for i := len(m.oracleCard.Rulings) - 1; i >= 0; i-- {
		ruleText := lipgloss.Wrap(m.oracleCard.Rulings[i].Comment, (m.winWidth/2)-5, "")
		rulings = fmt.Sprintf("%v%v\n%v\n\n", rulings, m.oracleCard.Rulings[i].PublishedAt.Format("02.01.2006"), ruleText)
	}
	return rulings
}

func renderMultifacesOracle(m model) string {
	bufferStruct := make([]string, len(*m.oracleCard.CardFaces))
	for i, f := range *m.oracleCard.CardFaces {
		width := m.winWidth / len(*m.oracleCard.CardFaces)
		bufferStruct[i] += f.ManaCost + "\n" + f.Name
		if f.TypeLine != nil {
			bufferStruct[i] += "\n" + *f.TypeLine
		}
		if f.Power != nil && f.Toughness != nil {
			bufferStruct[i] += "\n" + *f.Power + "/" + *f.Toughness
		} else if f.Defense != nil {
			bufferStruct[i] += "\n" + *f.Defense
		} else if f.Loyalty != nil {
			bufferStruct[i] += "\n" + *f.Loyalty
		} else {
			bufferStruct[i] += "\n"
		}
		if f.OracleText != nil {
			bufferStruct[i] += "\n" + "\n" + *f.OracleText
		}
		if f.PrintedText != nil {
			bufferStruct[i] += "\n" + "\n" + *f.PrintedText
		}
		if f.Cmc != nil {
			bufferStruct[i] += fmt.Sprintf("Cmc: %v", *f.Cmc)
		}
		bufferStruct[i] = lipgloss.Wrap(bufferStruct[i], width-10, "\n")
		bufferStruct[i] = lipgloss.PlaceHorizontal(width, lipgloss.Center, bufferStruct[i])
		bufferStruct[i] = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Width(width).Render(bufferStruct[i])
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, bufferStruct...)
}
