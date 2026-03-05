package scryfall

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/hugermuger/battlesphere/internal/database"
	"github.com/hugermuger/battlesphere/internal/dbutils"
	"github.com/hugermuger/battlesphere/internal/types"
)

func SingleCardImport(ctx context.Context, qtx *database.Queries, card types.CardJSON) error {
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
		Image:           toNullImage(card.ImageUris, "normal"),
		ImagePng:        toNullImage(card.ImageUris, "png"),
		ImageLarge:      toNullImage(card.ImageUris, "large"),
		ImageSmall:      toNullImage(card.ImageUris, "small"),
		ImageCrop:       toNullImage(card.ImageUris, "crop"),
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
		CardID:          card.ID,
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

	err = qtx.InsertCard(ctx, cardParams)
	if err != nil {
		return err
	}

	err = qtx.InsertLegality(ctx, legalitiesParam)
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
				Image:           toNullImage(c.ImageUris, "normal"),
				ImagePng:        toNullImage(c.ImageUris, "png"),
				ImageLarge:      toNullImage(c.ImageUris, "large"),
				ImageSmall:      toNullImage(c.ImageUris, "small"),
				ImageCrop:       toNullImage(c.ImageUris, "crop"),
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

			err := qtx.InsertCardFace(ctx, faceParams)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func SingleRuleImport(ctx context.Context, qtx *database.Queries, rule types.Rulings) error {
	layout := "2006-01-02"
	publishedAt, err := time.Parse(layout, rule.PublishedAt)
	if err != nil {
		return err
	}

	ruleParams := database.InsertRulingsParams{
		OracleID:    rule.OracleID,
		Source:      dbutils.ToNullString(rule.Source),
		PublishedAt: publishedAt,
		Comment:     rule.Comment,
	}

	err = qtx.InsertRulings(ctx, ruleParams)
	if err != nil {
		return err
	}

	return nil
}

func GetBulkURL(url, bulktype string) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", "BattlesphereApp/1.0")
	req.Header.Set("Accept", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("scryfall error: %s", res.Status)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	var urls types.URLS
	err = json.Unmarshal(body, &urls)
	if err != nil {
		return "", err
	}

	for _, u := range urls.Data {
		if u.Type == bulktype {
			return u.DownloadURI, nil
		}
	}

	return "", fmt.Errorf("Bulk Type not supoorted")
}

func toNullImage(images *types.ImageUris, size string) sql.NullString {
	if images == nil {
		return sql.NullString{Valid: false}
	}
	switch size {
	case "normal":
		return dbutils.ToNullString(images.Normal)
	case "png":
		return dbutils.ToNullString(images.Png)
	case "large":
		return dbutils.ToNullString(images.Large)
	case "crop":
		return dbutils.ToNullString(images.ArtCrop)
	case "small":
		return dbutils.ToNullString(images.Small)
	default:
		return sql.NullString{Valid: false}
	}
}
