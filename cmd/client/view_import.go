package main

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
)

func importView(m model) string {
	var s strings.Builder
	if !m.filepicker.collectionFilepicker {
		if m.filepicker.collectionImport {
			fmt.Fprintf(&s, "%s\n\n", lipgloss.PlaceHorizontal(m.winWidth, lipgloss.Center, "Import started!"))
			if m.filepicker.err != "" {
				fmt.Fprintf(&s, "Error: %s\n\n", m.filepicker.err)
			}
			if m.filepicker.status != "" {
				fmt.Fprintf(&s, "Status: %s\n", m.filepicker.status)
				fmt.Fprintf(&s, "Progress: %d cards imported\n\n", m.filepicker.progress)
			}
		} else {
			fmt.Fprintf(&s, "%s\n\n", lipgloss.PlaceHorizontal(m.winWidth, lipgloss.Center, "What do you want to import?"))
			if m.focusTabs {
				fmt.Fprintf(&s, "%s\n\n", lipgloss.PlaceHorizontal(m.winWidth, lipgloss.Center, m.login.blurredStyle.Render("[ Collection ]")+"        "+m.login.blurredStyle.Render("[ Deck ]")))
			} else if m.filepicker.focusIndex != 0 {
				fmt.Fprintf(&s, "%s\n\n", lipgloss.PlaceHorizontal(m.winWidth, lipgloss.Center, m.login.blurredStyle.Render("[ Collection ]")+"        "+m.login.focusedStyle.Render("[ Deck ]")))
			} else {
				fmt.Fprintf(&s, "%s\n\n", lipgloss.PlaceHorizontal(m.winWidth, lipgloss.Center, m.login.focusedStyle.Render("[ Collection ]")+"        "+m.login.blurredStyle.Render("[ Deck ]")))
			}
			if m.filepicker.err != "" {
				fmt.Fprintf(&s, "Error: %s\n\n", m.filepicker.err)
			}
			if m.filepicker.status != "" {
				fmt.Fprintf(&s, "Status: %s\n", m.filepicker.status)
				fmt.Fprintf(&s, "Progress: %d cards imported\n\n", m.filepicker.progress)
			}
		}
	} else {
		fmt.Fprintf(&s, "%s\n\n", lipgloss.PlaceHorizontal(m.winWidth, lipgloss.Center, "Choose file:"))
		s.WriteString("\n\n" + m.filepicker.model.View() + "\n")
	}
	return s.String()
}
