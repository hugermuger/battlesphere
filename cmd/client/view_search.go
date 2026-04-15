package main

import (
	"fmt"

	"charm.land/lipgloss/v2"
)

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
