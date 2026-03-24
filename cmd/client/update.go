package main

import (
	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"github.com/hugermuger/battlesphere/internal/types"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.winWidth = msg.Width
		m.search.listSearchByName.SetSize(msg.Width, msg.Height-10)
		m.search.listSearchByOracle.SetSize(msg.Width, msg.Height-10)
		switch tuiTabs[m.activeTabIndex].title {
		case "Card Search":
			if len(m.search.oracleCard.Rulings) > 0 {
				m.search.rulingViewport.SetContent(renderRulings(m))
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
		case "User":
			if !m.login.loggedIn {
				switch msg.String() {
				case "esc":
					m.focusTabs = true
					if m.login.registerView {
						for i, _ := range m.login.registerInput {
							m.login.registerInput[i].Blur()
						}
					} else {
						for i, _ := range m.login.loginInput {
							m.login.loginInput[i].Blur()
						}
					}

				case "tab":
					if m.focusTabs {
						return m, nil
					}
					if m.login.registerView {
						m.login.registerView = false
						m.login.loginInput[0].Focus()
					} else {
						m.login.registerView = true
						m.login.registerInput[0].Focus()
					}
					m.login.focusIndex = 0
					return m, nil

				case "enter", "down":
					if m.focusTabs {
						m.login.focusIndex = 0
						m.focusTabs = false
						if m.login.registerView {
							m.login.registerInput[0].Focus()
						} else {
							m.login.loginInput[0].Focus()
						}
						return m, nil
					}
					inputList := []textinput.Model{}
					if m.login.registerView {
						inputList = m.login.registerInput
					} else {
						inputList = m.login.loginInput
					}

					if msg.String() == "enter" && m.login.focusIndex == len(inputList) {
						if m.login.registerView {
							handlerCreateUser(inputList, &m)
						} else {
							handlerLogin(m.login.loginInput[1].Value(), m.login.loginInput[0].Value(), &m)
						}
						return m, nil
					}
					m.login.focusIndex++
					if m.login.focusIndex > len(inputList) {
						m.login.focusIndex = 0
					}
					if m.login.focusIndex == 0 {
						inputList[m.login.focusIndex].Focus()
					} else if m.login.focusIndex < len(inputList) {
						inputList[m.login.focusIndex-1].Blur()
						inputList[m.login.focusIndex].Focus()
					} else if m.login.focusIndex == len(inputList) {
						inputList[m.login.focusIndex-1].Blur()
					}
					return m, nil

				case "up":
					m.login.focusIndex--
					if m.login.registerView {
						if m.login.focusIndex < 0 {
							m.login.focusIndex = len(m.login.registerInput)
						}
						if m.login.focusIndex == len(m.login.registerInput) {
							m.login.registerInput[0].Blur()
						} else if m.login.focusIndex < len(m.login.registerInput)-1 {
							m.login.registerInput[m.login.focusIndex+1].Blur()
							m.login.registerInput[m.login.focusIndex].Focus()
						} else if m.login.focusIndex == len(m.login.registerInput)-1 {
							m.login.registerInput[m.login.focusIndex].Focus()
						}
					} else {
						if m.login.focusIndex < 0 {
							m.login.focusIndex = len(m.login.loginInput)
						}
						if m.login.focusIndex == len(m.login.loginInput) {
							m.login.loginInput[0].Blur()
						} else if m.login.focusIndex < len(m.login.loginInput)-1 {
							m.login.loginInput[m.login.focusIndex+1].Blur()
							m.login.loginInput[m.login.focusIndex].Focus()
						} else if m.login.focusIndex == len(m.login.loginInput)-1 {
							m.login.loginInput[m.login.focusIndex].Focus()
						}
					}
					return m, nil
				}

				inputList := []textinput.Model{}
				if m.login.registerView {
					inputList = m.login.registerInput
				} else {
					inputList = m.login.loginInput
				}

				if m.login.focusIndex < len(inputList) {
					if m.login.registerView {
						m.login.registerInput[m.login.focusIndex], cmd = m.login.registerInput[m.login.focusIndex].Update(msg)
					} else {
						m.login.loginInput[m.login.focusIndex], cmd = m.login.loginInput[m.login.focusIndex].Update(msg)
					}
					return m, cmd
				}
				return m, nil
			} else {
				switch msg.String() {
				case "enter":
					if m.focusTabs {
						m.focusTabs = false
						return m, nil
					}
					if m.login.logoutView {
						handlerLogout(m.refreshToken, &m)
					}
				case "esc":
					m.focusTabs = true
				case "tab":
					if m.focusTabs {
						return m, nil
					}
					if m.login.logoutView {
						m.login.logoutView = false
					} else {
						m.login.logoutView = true
					}
				}
				return m, nil
			}

		case "Card Search":
			switch msg.String() {
			case "enter":
				switch m.search.menuSearchIndex {
				case 0:
					if m.search.focusInput {
						m.search.focusList = true
						m.search.focusInput = false
						m.search.searchInput.Blur()
						return m, nil
					} else if m.focusTabs {
						m.focusTabs = false
						m.search.focusInput = true
						m.search.focusList = false
						m.search.searchInput.Focus()
						return m, nil
					} else if m.search.focusList {
						selectedItem := m.search.listSearchByName.SelectedItem()
						if selectedItem != nil {
							card := selectedItem.(types.CardSearchItem)
							if card.OracleID != nil {
								m.searching = true
								m.search.rulingViewport.GotoTop()
								m.search.menuSearchIndex = 1
								m.search.oracleCardID = *card.OracleID
								return m, m.fetchCardsByOracleID(card.OracleID.String())
							}
						}
					}
				case 1:
					if m.focusTabs {
						m.focusTabs = false
						m.search.focusList = true
					}
				}

			case "right":
				switch m.search.menuSearchIndex {
				case 0:
					m.search.listSearchByName, cmd = m.search.listSearchByName.Update(msg)
					return m, nil
				case 1:
					if m.search.focusList {
						m.search.focusList = false
						m.search.focusViewport = true
						m.search.listSearchByOracle.ResetSelected()
					}
				}
			case "left":
				switch m.search.menuSearchIndex {
				case 0:
					m.search.listSearchByName, cmd = m.search.listSearchByName.Update(msg)
					return m, nil
				case 1:
					if m.search.focusViewport {
						m.search.focusList = true
						m.search.focusViewport = false
					}
				}

			case "down":
				if m.search.focusInput {
					m.search.focusList = true
					m.search.focusInput = false
					m.search.searchInput.Blur()
					return m, nil
				}

				if m.search.focusList {
					switch m.search.menuSearchIndex {
					case 0:
						m.search.listSearchByName, cmd = m.search.listSearchByName.Update(msg)
						return m, nil
					case 1:
						m.search.listSearchByOracle, cmd = m.search.listSearchByOracle.Update(msg)
						return m, nil
					}
				} else if m.search.focusViewport {
					m.search.rulingViewport, cmd = m.search.rulingViewport.Update(msg)
					return m, nil
				}

			case "up":
				if m.search.focusList {
					switch m.search.menuSearchIndex {
					case 0:
						m.search.listSearchByName, cmd = m.search.listSearchByName.Update(msg)
						return m, nil
					case 1:
						m.search.listSearchByOracle, cmd = m.search.listSearchByOracle.Update(msg)
						return m, nil
					}
				} else if m.search.focusViewport {
					m.search.rulingViewport, cmd = m.search.rulingViewport.Update(msg)
					return m, nil
				}

			case "ctrl+c":
				m.search.searchInput.Reset()
				m.search.searchQuery = ""
				m.search.listSearchByName.SetItems([]list.Item{})
				m.search.listSearchByOracle.SetItems([]list.Item{})
				m.search.menuSearchIndex = 0
				if m.search.focusList || m.search.focusViewport {
					m.search.focusList = false
					m.search.focusViewport = false
					m.search.focusInput = true
					m.search.searchInput.Focus()
					return m, nil
				}

			case "esc":
				if m.search.focusInput {
					m.search.focusInput = false
					m.search.searchInput.Blur()
				}
				m.focusTabs = true
				m.search.focusList = false
				m.search.focusViewport = false
				return m, nil

			case "ctrl+l":
				m.search.selectedLang++
				if m.search.selectedLang > len(lang)-1 {
					m.search.selectedLang = 0
				}
				switch m.search.menuSearchIndex {
				case 0:
					m.search.searchQuery = m.search.searchInput.Value()
					m.search.searchQueryLang = m.search.selectedLang
					m.search.searchID++
					return m, tea.Batch(cmd, m.debounceSearch(m.search.searchID, m.search.searchInput.Value(), m.search.selectedLang))
				case 1:
					return m, tea.Batch(cmd, m.fetchCardsByOracleID(m.search.oracleCardID.String()))
				}

			case "backspace":
				switch m.search.menuSearchIndex {
				case 0:
					if m.search.focusList {
						m.search.focusInput = true
						m.search.focusList = false
						m.search.searchInput.Focus()
						return m, nil
					}
				case 1:
					m.search.menuSearchIndex = 0
					m.search.selectedLang = m.search.searchQueryLang
					m.search.focusList = true
					m.search.focusViewport = false
					return m, nil
				}
			}

			if m.search.focusInput {
				m.search.searchInput, cmd = m.search.searchInput.Update(msg)
			}

			if m.search.menuSearchIndex == 0 {
				if m.search.searchInput.Value() != m.search.searchQuery || m.search.searchQueryLang != m.search.selectedLang {
					m.search.searchQuery = m.search.searchInput.Value()
					m.search.searchQueryLang = m.search.selectedLang
					m.search.searchID++
					return m, tea.Batch(cmd, m.debounceSearch(m.search.searchID, m.search.searchInput.Value(), m.search.selectedLang))
				}
			}

			return m, cmd
		}

	case debounceMsg:
		switch tuiTabs[m.activeTabIndex].title {
		case "Card Search":
			if msg.id == m.search.searchID && msg.query != "" {
				m.searching = true
				return m, m.fetchCardsByName(msg.query)
			} else if msg.query == "" {
				m.search.listSearchByName.SetItems([]list.Item{})
			}
			return m, nil
		}

	case searchResultMsg:
		switch tuiTabs[m.activeTabIndex].title {
		case "Card Search":
			m.searching = false
			m.search.listSearchByName.Select(0)
			m.search.listSearchByName.SetItems(msg)
			return m, nil
		}

	case oracleSearchResultMsg:
		switch tuiTabs[m.activeTabIndex].title {
		case "Card Search":
			m.searching = false
			m.search.listSearchByOracle.Select(0)
			m.search.listSearchByOracle.SetItems(msg.list)
			m.search.oracleCard = msg.oracleCard
			if len(m.search.oracleCard.Rulings) > 0 {
				m.search.rulingViewport.SetContent(renderRulings(m))
			} else {
				m.search.rulingViewport.SetContent("")
			}
			return m, nil
		}

	case error:
		m.err = msg
		return m, nil
	}

	return m, cmd
}
