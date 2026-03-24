package main

import (
	"fmt"
	"reflect"
	"strings"

	"charm.land/bubbles/v2/textinput"
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

func loginView(m model) string {
	var b strings.Builder

	if m.login.loggedIn {
		fmt.Fprintf(&b, "%s\n\n", lipgloss.PlaceHorizontal(m.winWidth, lipgloss.Center, fmt.Sprintf("Welcome %v, you are logged in!", m.username)))
		if m.focusTabs {
			fmt.Fprintf(&b, "%s\n\n", lipgloss.PlaceHorizontal(m.winWidth, lipgloss.Center, m.login.blurredStyle.Render("[ Logout ]")+"        "+m.login.blurredStyle.Render("[ Delete Account ]")))
		} else if !m.login.logoutView {
			fmt.Fprintf(&b, "%s\n\n", lipgloss.PlaceHorizontal(m.winWidth, lipgloss.Center, m.login.blurredStyle.Render("[ Logout ]")+"        "+m.login.focusedStyle.Render("[ Delete Account ]")))
		} else {
			fmt.Fprintf(&b, "%s\n\n", lipgloss.PlaceHorizontal(m.winWidth, lipgloss.Center, m.login.focusedStyle.Render("[ Logout ]")+"        "+m.login.blurredStyle.Render("[ Delete Account ]")))
		}
	} else {
		inputs := []textinput.Model{}

		if m.login.registerView {
			inputs = m.login.registerInput
		} else {
			inputs = m.login.loginInput
		}

		if m.focusTabs {
			fmt.Fprintf(&b, "%s\n\n", lipgloss.PlaceHorizontal(m.winWidth, lipgloss.Center, m.login.blurredStyle.Render("[ Login ]")+"        "+m.login.blurredStyle.Render("[ Register ]")))
		} else if m.login.registerView {
			fmt.Fprintf(&b, "%s\n\n", lipgloss.PlaceHorizontal(m.winWidth, lipgloss.Center, m.login.blurredStyle.Render("[ Login ]")+"        "+m.login.focusedStyle.Render("[ Register ]")))
		} else {
			fmt.Fprintf(&b, "%s\n\n", lipgloss.PlaceHorizontal(m.winWidth, lipgloss.Center, m.login.focusedStyle.Render("[ Login ]")+"        "+m.login.blurredStyle.Render("[ Register ]")))
		}

		for i, in := range inputs {
			b.WriteString(lipgloss.PlaceHorizontal(m.winWidth, lipgloss.Center, in.View()))
			if i < len(inputs)-1 {
				b.WriteRune('\n')
			}
		}

		button := m.login.blurredStyle.Render("[ Submit ]")
		if m.login.focusIndex == len(inputs) {
			button = m.login.focusedStyle.Render("[ Submit ]")
		}
		fmt.Fprintf(&b, "\n\n%s\n\n", lipgloss.PlaceHorizontal(m.winWidth, lipgloss.Center, button))

		if m.login.err != "" {
			fmt.Fprintf(&b, "%s\n\n", lipgloss.PlaceHorizontal(m.winWidth, lipgloss.Center, lipgloss.NewStyle().Foreground(lipgloss.Red).Render(m.login.err)))
		}

		if m.login.registerSucces {
			fmt.Fprintf(&b, "%s\n\n", lipgloss.PlaceHorizontal(m.winWidth, lipgloss.Center, lipgloss.NewStyle().Foreground(lipgloss.Green).Render("Your registration was successfull!")))
		}

	}
	return lipgloss.PlaceHorizontal(m.winWidth, lipgloss.Center, b.String())
}

func searchView(m model) string {
	if m.search.focusList {
		m.search.listSearchByName.SetDelegate(m.listStyles.selectedDelegate)
		m.search.listSearchByOracle.SetDelegate(m.listStyles.selectedDelegate)
	} else {
		m.search.listSearchByName.SetDelegate(m.listStyles.unselectedDelegate)
		m.search.listSearchByOracle.SetDelegate(m.listStyles.unselectedDelegate)
	}

	rulings := "\nRulings"
	if m.search.focusViewport {
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
		if i == m.search.selectedLang {
			la = fmt.Sprintf("%v [%v]", la, l)
		} else {
			la += " " + l
		}
	}

	la = lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Render(la)

	v := ""

	switch m.search.menuSearchIndex {
	case 0:
		if len(m.search.listSearchByName.Items()) == 0 {
			v = lipgloss.PlaceHorizontal(m.winWidth, lipgloss.Center, m.search.searchInput.View()) + "\n" + "\n" + lipgloss.PlaceHorizontal(m.winWidth, lipgloss.Center, la)
		} else {
			v = lipgloss.PlaceHorizontal(m.winWidth, lipgloss.Center, m.search.searchInput.View()) + "\n" + "\n" + m.search.listSearchByName.View() + "\n" + lipgloss.PlaceHorizontal(m.winWidth, lipgloss.Center, la)
		}
	case 1:
		buffer := ""
		if m.search.oracleCard.ManaCost != nil {
			buffer += *m.search.oracleCard.ManaCost + "\n"
		}
		buffer += m.search.oracleCard.Name + "\n" + m.search.oracleCard.TypeLine
		if m.search.oracleCard.Power != nil && m.search.oracleCard.Toughness != nil {
			buffer += "\n" + *m.search.oracleCard.Power + "/" + *m.search.oracleCard.Toughness
		}
		if m.search.oracleCard.Defense != nil {
			buffer += "\n" + *m.search.oracleCard.Defense
		}
		if m.search.oracleCard.Loyalty != nil {
			buffer += "\n" + *m.search.oracleCard.Loyalty
		}
		if m.search.oracleCard.OracleText != nil {
			buffer += "\n\n" + *m.search.oracleCard.OracleText
		}
		if m.search.oracleCard.PrintedText != nil {
			buffer += "\n\n" + *m.search.oracleCard.PrintedText
		}
		buffer = fmt.Sprintf("%v\n\nCmc: %v", buffer, m.search.oracleCard.Cmc)
		if len(m.search.oracleCard.ColorIdentity) > 0 {
			buffer += ", Color Identity: "
			for i := len(m.search.oracleCard.ColorIdentity) - 1; i >= 0; i-- {
				buffer += fmt.Sprintf("{%v}", m.search.oracleCard.ColorIdentity[i])
			}
		}
		if m.search.oracleCard.EdhrecRank != nil {
			buffer = fmt.Sprintf("%v\n\nEDHREC Rank: %v", buffer, *m.search.oracleCard.EdhrecRank)
		}
		if m.search.oracleCard.GameChanger != nil {
			if *m.search.oracleCard.GameChanger {
				buffer += "\n\n" + lipgloss.NewStyle().Foreground(lipgloss.Color("220")).Render("GAMECHANGER")
			}
		}
		buffer = lipgloss.Wrap(buffer, m.winWidth-15, "\n")
		buffer = lipgloss.PlaceHorizontal(m.winWidth, lipgloss.Center, buffer)
		card := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Width(m.winWidth).Render(buffer)
		legalities := renderLegalities(m)
		m.search.rulingViewport.SetWidth(m.winWidth / 2)
		m.search.searchViewport.SetWidth(m.winWidth / 2)
		if m.search.oracleCard.Multifaced {
			faces := renderMultifacesOracle(m)
			height := m.search.listSearchByName.Height() - (lipgloss.Height(buffer) + lipgloss.Height(faces) + lipgloss.Height(legalities))
			m.search.rulingViewport.SetHeight(height - 5)
			m.search.listSearchByOracle.SetHeight(height - 3)
			m.search.searchViewport.SetContent(m.search.listSearchByOracle.View())
			m.search.searchViewport.SetHeight(height - 2)
			listView := lipgloss.JoinHorizontal(lipgloss.Top, m.search.searchViewport.View(), rulings+"\n\n"+m.search.rulingViewport.View())
			v = card + "\n" + legalities + "\n" + faces + "\n" + listView + "\n\n" + lipgloss.PlaceHorizontal(m.winWidth, lipgloss.Center, la)
		} else {
			height := m.search.listSearchByName.Height() - (lipgloss.Height(buffer) + lipgloss.Height(legalities))
			m.search.rulingViewport.SetHeight(height - 5)
			m.search.listSearchByOracle.SetHeight(height - 3)
			m.search.searchViewport.SetContent(m.search.listSearchByOracle.View())
			m.search.searchViewport.SetHeight(height - 2)
			listView := lipgloss.JoinHorizontal(lipgloss.Top, m.search.searchViewport.View(), rulings+"\n\n"+m.search.rulingViewport.View())
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
	l := m.search.oracleCard.Legalities
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
	for i := len(m.search.oracleCard.Rulings) - 1; i >= 0; i-- {
		ruleText := lipgloss.Wrap(m.search.oracleCard.Rulings[i].Comment, (m.winWidth/2)-5, "")
		rulings = fmt.Sprintf("%v%v\n%v\n\n", rulings, m.search.oracleCard.Rulings[i].PublishedAt.Format("02.01.2006"), ruleText)
	}
	return rulings
}

func renderMultifacesOracle(m model) string {
	bufferStruct := make([]string, len(*m.search.oracleCard.CardFaces))
	for i, f := range *m.search.oracleCard.CardFaces {
		width := m.winWidth / len(*m.search.oracleCard.CardFaces)
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
