package main

import (
	"context"
	"strconv"
	"strings"
	"time"

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

func mapDragonShield(line []string, cfg *apiConfig, c context.Context) (database.AddCardToCollectionsParams, error) {
	scryfallParams := database.GetOneUniqueCardParams{
		Name:            line[3],
		SetCode:         strings.ToLower(line[4]),
		CollectorNumber: line[6],
		Lang:            langMap[line[9]],
	}

	scryfallID, err := cfg.db.GetOneUniqueCard(c, scryfallParams)
	if err != nil {
		return database.AddCardToCollectionsParams{}, err
	}

	parsedDate, _ := time.Parse("2006-01-02", line[11])

	condition := ""

	switch line[7] {
	case "Normal":
		condition = "nonfoil"
	case "Foil":
		condition = "foil"
	default:
		condition = "etched"
	}

	params := database.AddCardToCollectionsParams{
		ScryfallID:    scryfallID,
		PurchaseDate:  parsedDate,
		PurchasePrice: toFloat(line[10]),
		Finish:        line[8],
		Condition:     condition,
	}

	return params, nil
}

func toInt(s string) int {
	v, _ := strconv.Atoi(s)
	return v
}

func toFloat(s string) float64 {
	v, _ := strconv.ParseFloat(s, 64)
	return v
}
