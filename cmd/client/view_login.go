package main

import (
	"fmt"
	"strings"

	"charm.land/bubbles/v2/textinput"
	"charm.land/lipgloss/v2"
)

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
