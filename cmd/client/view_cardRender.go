package main

import (
	"fmt"
	"reflect"

	"charm.land/lipgloss/v2"
)

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
