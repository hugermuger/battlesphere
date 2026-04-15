package main

import (
	"strings"
)

func importView(m model) string {
	var s strings.Builder
	if m.focusTabs {
		s.WriteString("\n" + "Click enter to import")
	} else {
		s.WriteString("\n" + m.fp.model.View() + "\n")
	}
	return s.String()
}
