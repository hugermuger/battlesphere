package main

import (
	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"github.com/hugermuger/battlesphere/internal/types"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.winWidth = msg.Width
		m.listSearchByName.SetSize(msg.Width, msg.Height-10)
		m.listSearchByOracle.SetSize(msg.Width, msg.Height-10)
		switch tuiTabs[m.activeTabIndex].title {
		case "Card Search":
			if len(m.oracleCard.Rulings) > 0 {
				m.rulingViewport.SetContent(renderRulings(m))
			}
		}

	case tea.KeyPressMsg:
		if m.focusTabs {
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
					} else if m.focusTabs {
						m.focusTabs = false
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
								m.rulingViewport.GotoTop()
								m.menuSearchID = 1
								m.oracleCardID = *card.OracleID
								return m, m.fetchCardsByOracleID(card.OracleID.String())
							}
						}
					}
				case 1:
					if m.focusTabs {
						m.focusTabs = false
						m.focusList = true
					}
				}

			case "right":
				switch m.menuSearchID {
				case 0:
					m.listSearchByName, cmd = m.listSearchByName.Update(msg)
					return m, nil
				case 1:
					if m.focusList {
						m.focusList = false
						m.focusViewport = true
						m.listSearchByOracle.ResetSelected()
					}
				}
			case "left":
				switch m.menuSearchID {
				case 0:
					m.listSearchByName, cmd = m.listSearchByName.Update(msg)
					return m, nil
				case 1:
					if m.focusViewport {
						m.focusList = true
						m.focusViewport = false
					}
				}

			case "down":
				if m.focusInput {
					m.focusList = true
					m.focusInput = false
					m.searchInput.Blur()
					return m, nil
				}

				if m.focusList {
					switch m.menuSearchID {
					case 0:
						m.listSearchByName, cmd = m.listSearchByName.Update(msg)
						return m, nil
					case 1:
						m.listSearchByOracle, cmd = m.listSearchByOracle.Update(msg)
						return m, nil
					}
				} else if m.focusViewport {
					m.rulingViewport, cmd = m.rulingViewport.Update(msg)
					return m, nil
				}

			case "up":
				if m.focusList {
					switch m.menuSearchID {
					case 0:
						m.listSearchByName, cmd = m.listSearchByName.Update(msg)
						return m, nil
					case 1:
						m.listSearchByOracle, cmd = m.listSearchByOracle.Update(msg)
						return m, nil
					}
				} else if m.focusViewport {
					m.rulingViewport, cmd = m.rulingViewport.Update(msg)
					return m, nil
				}

			case "ctrl+c":
				m.searchInput.Reset()
				m.searchQuery = ""
				m.listSearchByName.SetItems([]list.Item{})
				m.listSearchByOracle.SetItems([]list.Item{})
				m.menuSearchID = 0
				if m.focusList || m.focusViewport {
					m.focusList = false
					m.focusViewport = false
					m.focusInput = true
					m.searchInput.Focus()
					return m, nil
				}

			case "esc":
				if m.focusInput {
					m.focusInput = false
					m.searchInput.Blur()
				}
				m.focusTabs = true
				m.focusList = false
				m.focusViewport = false
				return m, nil

			case "ctrl+l":
				m.selectedLang++
				if m.selectedLang > len(lang)-1 {
					m.selectedLang = 0
				}
				switch m.menuSearchID {
				case 0:
					m.searchQuery = m.searchInput.Value()
					m.searchQueryLang = m.selectedLang
					m.searchID++
					return m, tea.Batch(cmd, m.debounceSearch(m.searchID, m.searchInput.Value(), m.selectedLang))
				case 1:
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
					m.focusList = true
					m.focusViewport = false
					return m, nil
				}
			}

			if m.focusInput {
				m.searchInput, cmd = m.searchInput.Update(msg)
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
			m.listSearchByName.Select(0)
			m.listSearchByName.SetItems(msg)
			return m, nil
		}

	case oracleSearchResultMsg:
		switch tuiTabs[m.activeTabIndex].title {
		case "Card Search":
			m.searching = false
			m.listSearchByOracle.Select(0)
			m.listSearchByOracle.SetItems(msg.list)
			m.oracleCard = msg.oracleCard
			if len(m.oracleCard.Rulings) > 0 {
				m.rulingViewport.SetContent(renderRulings(m))
			} else {
				m.rulingViewport.SetContent("")
			}
			return m, nil
		}

	case error:
		m.err = msg
		return m, nil
	}

	return m, cmd
}
