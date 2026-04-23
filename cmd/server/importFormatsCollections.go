package main

import (
	"context"
	"database/sql"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hugermuger/battlesphere/internal/database"
)

var langMap = map[string]string{
	"English":             "en",
	"Spanish":             "es",
	"French":              "fr",
	"German":              "de",
	"Italian":             "it",
	"Portuguese":          "pt",
	"Japanese":            "ja",
	"Korean":              "ko",
	"Russian":             "ru",
	"Simplified Chinese":  "zhs",
	"Traditional Chinese": "zht",
	"Hebrew":              "he",
	"Latin":               "la",
	"Ancient Greek":       "grc",
	"Arabic":              "ar",
	"Sanskrit":            "sa",
	"Phyrexian":           "ph",
	"Quenya":              "qya",
}

func mapDragonShield(line []string, cfg *apiConfig, c context.Context) (database.AddCardToCollectionsParams, int, string, error) {
	quantity := toInt(line[1])
	folderName := line[0]

	var scryfallID uuid.UUID

	line[3] = strings.ReplaceAll(line[3], " Token", "")

	setCode := strings.ToLower(line[4])
	if setCode == "fwb" {
		setCode = "3ed"
	}
	if setCode == "vthb" {
		setCode = "thb"
	}
	if setCode == "pw09" {
		setCode = "dci"
	}
	if setCode == "gk1_dimir" {
		setCode = "gk1"
	}
	if setCode == "mb1" || setCode == "fmb1" {
		setCode = "plst"
		scryfallParams := database.GetOneUniqueCardTheListParams{
			Column1: sql.NullString{
				String: line[3],
				Valid:  true,
			},
			SetCode: setCode,
			Lang:    langMap[line[9]],
		}

		ID, err := cfg.db.GetOneUniqueCardTheList(c, scryfallParams)
		if err != nil {
			return database.AddCardToCollectionsParams{}, 0, "", err
		}

		scryfallID = ID
	} else {
		scryfallParams := database.GetOneUniqueCardParams{
			Column1: sql.NullString{
				String: line[3],
				Valid:  true,
			},
			SetCode:         setCode,
			CollectorNumber: line[6],
			Lang:            langMap[line[9]],
		}

		ID, err := cfg.db.GetOneUniqueCard(c, scryfallParams)
		if err != nil {
			return database.AddCardToCollectionsParams{}, 0, "", err
		}

		scryfallID = ID
	}

	parsedDate, _ := time.Parse("2006-01-02", line[11])

	finish := ""

	switch line[8] {
	case "Normal":
		finish = "nonfoil"
	case "Foil":
		finish = "foil"
	default:
		finish = "etched"
	}

	params := database.AddCardToCollectionsParams{
		ScryfallID:    scryfallID,
		PurchaseDate:  parsedDate,
		PurchasePrice: toFloat(line[10]),
		Finish:        finish,
		Condition:     line[7],
	}

	return params, quantity, folderName, nil
}

func toInt(s string) int {
	v, _ := strconv.Atoi(s)
	return v
}

func toFloat(s string) float64 {
	v, _ := strconv.ParseFloat(s, 64)
	return v
}
