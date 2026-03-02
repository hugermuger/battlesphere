package main

import (
	"context"
	"time"

	"github.com/hugermuger/battlesphere/internal/database"
	"github.com/hugermuger/battlesphere/internal/dbutils"
	"github.com/hugermuger/battlesphere/internal/scryfall"
)

func (cfg *apiConfig) handler_importSingleCardToDB(card scryfall.CardJSON) error {
	tx, err := cfg.dbConn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	qtx := cfg.db.WithTx(tx)

	layout := "2006-01-02"
	relaseDate, err := time.Parse(layout, card.ReleasedAt)
	if err != nil {
		return err
	}

	priceUsd, err := dbutils.StringToNullFloat64(card.Prices.Usd)
	if err != nil {
		return err
	}

	priceEur, err := dbutils.StringToNullFloat64(card.Prices.Eur)
	if err != nil {
		return err
	}

	priceUsdFoil, err := dbutils.StringToNullFloat64(card.Prices.UsdFoil)
	if err != nil {
		return err
	}

	priceEurFoil, err := dbutils.StringToNullFloat64(card.Prices.EurFoil)
	if err != nil {
		return err
	}

	priceUsdEtched, err := dbutils.StringToNullFloat64(card.Prices.UsdEtched)
	if err != nil {
		return err
	}

	cardParams := database.InsertCardParams{
		ID:              card.ID,
		ArenaID:         dbutils.ToNullInt64(card.ArenaID),
		MtgoID:          dbutils.ToNullInt64(card.MtgoID),
		CardmarketID:    dbutils.ToNullInt64(card.CardmarketID),
		OracleID:        dbutils.ToNullUUID(card.OracleID),
		ReleaseDate:     relaseDate,
		Lang:            card.Lang,
		Layout:          card.Layout,
		RulingUri:       card.RulingsURI,
		EdhrecRank:      dbutils.ToNullInt32(card.EdhrecRank),
		GameChanger:     dbutils.ToNullBool(card.GameChanger),
		Multifaced:      card.CardFaces != nil,
		Cmc:             card.Cmc,
		ColorIdentity:   card.ColorIdentity,
		Colors:          dbutils.SafeSlice(card.Colors),
		Defense:         dbutils.ToNullString(card.Defense),
		Keywords:        card.Keywords,
		Loyalty:         dbutils.ToNullString(card.Loyalty),
		ManaCost:        dbutils.ToNullString(card.ManaCost),
		Name:            card.Name,
		OracleText:      dbutils.ToNullString(card.OracleText),
		Power:           dbutils.ToNullString(card.Power),
		Toughness:       dbutils.ToNullString(card.Toughness),
		TypeLine:        card.TypeLine,
		Artist:          dbutils.ToNullString(card.Artist),
		CollectorNumber: card.CollectorNumber,
		Finishes:        card.Finishes,
		FlavorName:      dbutils.ToNullString(card.FlavorName),
		FlavorText:      dbutils.ToNullString(card.FlavorText),
		Games:           card.Games,
		Image:           dbutils.ToNullImage(card.ImageUris, "normal"),
		ImagePng:        dbutils.ToNullImage(card.ImageUris, "png"),
		ImageLarge:      dbutils.ToNullImage(card.ImageUris, "large"),
		ImageSmall:      dbutils.ToNullImage(card.ImageUris, "small"),
		ImageCrop:       dbutils.ToNullImage(card.ImageUris, "crop"),
		PriceUsd:        priceUsd,
		PriceEur:        priceEur,
		PriceFoilUsd:    priceUsdFoil,
		PriceFoilEur:    priceEurFoil,
		PriceEtchedUsd:  priceUsdEtched,
		PrintedName:     dbutils.ToNullString(card.PrintedName),
		PrintedText:     dbutils.ToNullString(card.PrintedText),
		PrintedTypeLine: dbutils.ToNullString(card.PrintedTypeLine),
		Rarity:          card.Rarity,
		SetName:         card.SetName,
		SetCode:         card.Set,
	}

	legalitiesParam := database.InsertLegalityParams{
		ID:              card.ID,
		Standard:        dbutils.ToNullString(card.Legalities.Standard),
		Pauper:          dbutils.ToNullString(card.Legalities.Pauper),
		Vintage:         dbutils.ToNullString(card.Legalities.Vintage),
		Pioneer:         dbutils.ToNullString(card.Legalities.Pioneer),
		Modern:          dbutils.ToNullString(card.Legalities.Modern),
		Legacy:          dbutils.ToNullString(card.Legalities.Legacy),
		Commander:       dbutils.ToNullString(card.Legalities.Commander),
		Future:          dbutils.ToNullString(card.Legalities.Future),
		Historic:        dbutils.ToNullString(card.Legalities.Historic),
		Timeless:        dbutils.ToNullString(card.Legalities.Timeless),
		Gladiator:       dbutils.ToNullString(card.Legalities.Gladiator),
		Penny:           dbutils.ToNullString(card.Legalities.Penny),
		Oathbreaker:     dbutils.ToNullString(card.Legalities.Oathbreaker),
		Standardbrawl:   dbutils.ToNullString(card.Legalities.Standardbrawl),
		Brawl:           dbutils.ToNullString(card.Legalities.Brawl),
		Alchemy:         dbutils.ToNullString(card.Legalities.Alchemy),
		Paupercommander: dbutils.ToNullString(card.Legalities.Paupercommander),
		Duel:            dbutils.ToNullString(card.Legalities.Duel),
		Oldschool:       dbutils.ToNullString(card.Legalities.Oldschool),
		Premodern:       dbutils.ToNullString(card.Legalities.Premodern),
		Predh:           dbutils.ToNullString(card.Legalities.Predh),
	}

	err = qtx.InsertCard(context.Background(), cardParams)
	if err != nil {
		return err
	}

	err = qtx.InsertLegality(context.Background(), legalitiesParam)
	if err != nil {
		return err
	}

	if card.CardFaces != nil {
		for i, c := range *card.CardFaces {

			faceParams := database.InsertCardFaceParams{
				CardID:          card.ID,
				FaceIndex:       int32(i),
				Artist:          dbutils.ToNullString(c.Artist),
				Cmc:             dbutils.ToNullFloat64(c.Cmc),
				Colors:          dbutils.SafeSlice(c.Colors),
				Defense:         dbutils.ToNullString(c.Defense),
				FlavorText:      dbutils.ToNullString(c.FlavorText),
				Image:           dbutils.ToNullImage(c.ImageUris, "normal"),
				ImagePng:        dbutils.ToNullImage(c.ImageUris, "png"),
				ImageLarge:      dbutils.ToNullImage(c.ImageUris, "large"),
				ImageSmall:      dbutils.ToNullImage(c.ImageUris, "small"),
				ImageCrop:       dbutils.ToNullImage(c.ImageUris, "crop"),
				Layout:          dbutils.ToNullString(c.Layout),
				Loyalty:         dbutils.ToNullString(c.Loyalty),
				ManaCost:        c.ManaCost,
				Name:            c.Name,
				OracleText:      dbutils.ToNullString(c.OracleText),
				Power:           dbutils.ToNullString(c.Power),
				PrintedName:     dbutils.ToNullString(c.PrintedName),
				PrintedText:     dbutils.ToNullString(c.PrintedText),
				PrintedTypeLine: dbutils.ToNullString(c.PrintedTypeLine),
				Toughness:       dbutils.ToNullString(c.Toughness),
				TypeLine:        dbutils.ToNullString(c.TypeLine),
			}

			err := qtx.InsertCardFace(context.Background(), faceParams)
			if err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}
