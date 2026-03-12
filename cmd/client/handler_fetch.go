package main

import (
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"github.com/hugermuger/battlesphere/internal/types"
)

func (m model) debounceSearch(id int, query string, langID int) tea.Cmd {
	return func() tea.Msg {
		time.Sleep(300 * time.Millisecond)
		return debounceMsg{id: id, query: query, langID: langID}
	}
}

func (m model) fetchCardsByName(query string) tea.Cmd {
	return func() tea.Msg {
		endpoint := "http://localhost:8080/cards/search?name=" + url.QueryEscape(query) + "&limit=600" + "&lang=" + lang[m.selectedLang]
		resp, err := http.Get(endpoint)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		var data struct {
			Results []types.CardResponseSearchByName `json:"results"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return err
		}

		var items []list.Item
		for _, v := range data.Results {
			items = append(items, v)
		}

		return searchResultMsg(items)
	}
}

func (m model) fetchCardsByOracleID(oracleID string) tea.Cmd {
	return func() tea.Msg {
		endpoint := "http://localhost:8080/cards/oracle/" + oracleID + "?limit=600" + "&lang=" + lang[m.selectedLang]
		resp, err := http.Get(endpoint)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		var data struct {
			Results    []types.CardResponseSearchByOracleID `json:"results"`
			OracleCard types.ResponseByOracleID             `json:"oracle_card"`
			Legalities types.Legalities                     `json:"legalities"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return err
		}

		var items []list.Item
		for _, v := range data.Results {
			items = append(items, v)
		}
		return oracleSearchResultMsg{list: items, oracleCard: data.OracleCard, legalities: data.Legalities}
	}
}

func (m model) fetchCardByID(cardID string) tea.Cmd {
	return func() tea.Msg {
		endpoint := "http://localhost:8080/cards/" + cardID
		resp, err := http.Get(endpoint)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		var data struct {
			Card    types.ResponseCard    `json:"card"`
			Rulings types.ResponseRulings `json:"rulings"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return err
		}

		return cardSearchResultMsg{card: data.Card, rulings: data.Rulings}
	}
}
