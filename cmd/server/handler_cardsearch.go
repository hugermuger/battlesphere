package main

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hugermuger/battlesphere/internal/database"
	"github.com/hugermuger/battlesphere/internal/dbutils"
	"github.com/hugermuger/battlesphere/internal/types"
)

func (cfg *apiConfig) handlerSearchCards(c *gin.Context) {
	const defaultLimit = "20"

	name := c.Query("name")
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", defaultLimit)
	lang := c.DefaultQuery("lang", "en")

	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name parameter is required"})
		return
	}

	_, err := cfg.db.DoesLangExist(c.Request.Context(), lang)
	if err == sql.ErrNoRows {
		lang = "en"
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "page parameter is in wrong format"})
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "limit parameter is in wrong format"})
		return
	}

	var numberResults int64

	if lang == "en" {
		numberResults, err = cfg.db.CountCardsByNameListEng(c.Request.Context(), dbutils.ToNullString(&name))
		if err != nil {
			c.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
			return
		}
	} else {
		paramNumber := database.CountCardsByNameListParams{
			Column1: dbutils.ToNullString(&name),
			Lang:    lang,
		}

		numberResults, err = cfg.db.CountCardsByNameList(c.Request.Context(), paramNumber)
		if err != nil {
			c.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
			return
		}
	}

	numberPages := calculatePages(numberResults, limit)

	if page > numberPages {
		if numberPages == 0 {
			page = 1
		} else {
			page = numberPages
		}
	}

	offset := (page - 1) * limit

	results := []types.CardResponseSearchByName{}

	if lang == "en" {
		cardArgs := database.SearchCardsByNameListEngParams{
			Column1: dbutils.ToNullString(&name),
			Limit:   int32(limit),
			Offset:  int32(offset),
		}

		cards, err := cfg.db.SearchCardsByNameListEng(c.Request.Context(), cardArgs)
		if err != nil {
			c.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
			return
		}

		resultsBuffer := make([]types.CardResponseSearchByName, len(cards))
		for i, card := range cards {
			resultsBuffer[i] = types.CardResponseSearchByName{
				OracleID: dbutils.ToUUIDPtr(card.OracleID),
				Name:     card.Name,
				Layout:   card.Layout,
				ManaCost: dbutils.FromNullString(card.ManaCost),
				TypeLine: card.TypeLine,
			}
		}

		results = resultsBuffer
	} else {
		cardArgs := database.SearchCardsByNameListParams{
			Column1: dbutils.ToNullString(&name),
			Lang:    lang,
			Limit:   int32(limit),
		}

		cards, err := cfg.db.SearchCardsByNameList(c.Request.Context(), cardArgs)
		if err != nil {
			c.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
			return
		}

		resultsBuffer := make([]types.CardResponseSearchByName, len(cards))
		for i, card := range cards {
			resultsBuffer[i] = types.CardResponseSearchByName{
				OracleID: dbutils.ToUUIDPtr(card.OracleID),
				Name:     card.PrintedName.String,
				Layout:   card.Layout,
				ManaCost: dbutils.FromNullString(card.ManaCost),
				TypeLine: card.PrintedTypeLine.String,
			}
		}

		results = resultsBuffer
	}

	next_page := ""
	if numberPages > 1 && page < numberPages {
		base, _ := url.Parse("/cards/search")

		params := url.Values{}
		params.Add("name", name)
		params.Add("limit", fmt.Sprintf("%d", limit))
		params.Add("lang", lang)
		params.Add("page", fmt.Sprintf("%d", page+1))

		base.RawQuery = params.Encode()
		next_page = base.String()
	}

	c.JSON(http.StatusOK, gin.H{
		"page":           page,
		"number_pages":   numberPages,
		"results":        results,
		"number_results": numberResults,
		"next_page":      next_page,
	})
}

func (cfg *apiConfig) handlerRulings(c *gin.Context) {
	idStr := c.Param("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "wrong UUID format"})
		return
	}

	rules, err := cfg.db.GetOracleRulings(c.Request.Context(), id)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}
	rulingsJSON := make([]types.ResponseRulings, len(rules))

	for i, rule := range rules {
		rulingsJSON[i] = types.ResponseRulings{
			OracleID:    rule.OracleID,
			Source:      dbutils.FromNullString(rule.Source),
			PublishedAt: rule.PublishedAt,
			Comment:     rule.Comment,
		}
	}

	c.JSON(http.StatusOK, rulingsJSON)
}

func (cfg *apiConfig) handlerCardsByOracleID(c *gin.Context) {
	const defaultLimit = "20"

	idStr := c.Param("id")

	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", defaultLimit)
	lang := c.DefaultQuery("lang", "en")

	id, err := uuid.Parse(idStr)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "wrong UUID format"})
		return
	}

	_, err = cfg.db.DoesLangExist(c.Request.Context(), lang)
	if err == sql.ErrNoRows {
		lang = "en"
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "page parameter is in wrong format"})
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "limit parameter is in wrong format"})
		return
	}

	numberParams := database.CountCardsByOracleIDListParams{
		OracleID: uuid.NullUUID{UUID: id, Valid: true},
		Lang:     lang,
	}

	numberResults, err := cfg.db.CountCardsByOracleIDList(c.Request.Context(), numberParams)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	numberPages := calculatePages(numberResults, limit)

	if page > numberPages {
		if numberPages == 0 {
			page = 1
		} else {
			page = numberPages
		}
	}

	offset := (page - 1) * limit

	oracleParams := database.SearchCardByOracleIDParams{
		OracleID: uuid.NullUUID{UUID: id, Valid: true},
		Lang:     lang,
	}

	oracleCard, err := cfg.db.SearchCardByOracleID(c.Request.Context(), oracleParams)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found OracleID"})
		return
	} else if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	oracleCardJSON := types.ResponseByOracleID{
		Name:          oracleCard.Name,
		Layout:        oracleCard.Layout,
		Cmc:           oracleCard.Cmc,
		Colors:        &oracleCard.Colors,
		ColorIdentity: oracleCard.ColorIdentity,
		ManaCost:      dbutils.FromNullString(oracleCard.ManaCost),
		TypeLine:      oracleCard.TypeLine,
		OracleText:    dbutils.FromNullString(oracleCard.OracleText),
		Power:         dbutils.FromNullString(oracleCard.Power),
		Toughness:     dbutils.FromNullString(oracleCard.Toughness),
		Loyalty:       dbutils.FromNullString(oracleCard.Loyalty),
		Defense:       dbutils.FromNullString(oracleCard.Defense),
		Multifaced:    oracleCard.Multifaced,
		GameChanger:   dbutils.FromNullBool(oracleCard.GameChanger),
		EdhrecRank:    dbutils.FromNullInt32(oracleCard.EdhrecRank),
	}

	multifaces := []database.CardFace{}

	if oracleCard.Multifaced {
		multifaces, err = cfg.db.GetCardFaces(c.Request.Context(), oracleCard.ID)
		if err != nil {
			c.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
			return
		}

		responseFaces := make([]types.CardFacesByOracleID, len(multifaces))

		for i, face := range multifaces {
			responseFaces[i] = types.CardFacesByOracleID{
				Name:       face.Name,
				ManaCost:   face.ManaCost,
				Cmc:        dbutils.FromNullFloat64(face.Cmc),
				Colors:     &face.Colors,
				TypeLine:   dbutils.FromNullString(face.TypeLine),
				OracleText: dbutils.FromNullString(face.OracleText),
				Power:      dbutils.FromNullString(face.Power),
				Toughness:  dbutils.FromNullString(face.Toughness),
				Loyalty:    dbutils.FromNullString(face.Loyalty),
				Defense:    dbutils.FromNullString(face.Defense),
			}

			if lang != "en" {
				responseFaces[i].Name = face.PrintedName.String
				responseFaces[i].TypeLine = dbutils.FromNullString(face.PrintedTypeLine)
				responseFaces[i].PrintedText = dbutils.FromNullString(face.PrintedText)
			}
		}
		oracleCardJSON.CardFaces = &responseFaces
	}

	if lang != "en" {
		oracleCardJSON.Name = oracleCard.PrintedName.String
		oracleCardJSON.TypeLine = oracleCard.PrintedTypeLine.String
		oracleCardJSON.PrintedText = dbutils.FromNullString(oracleCard.PrintedText)
	}

	cardParams := database.SearchCardsByOracleIDListParams{
		OracleID: uuid.NullUUID{UUID: id, Valid: true},
		Lang:     lang,
		Limit:    int32(limit),
		Offset:   int32(offset),
	}

	cards, err := cfg.db.SearchCardsByOracleIDList(c.Request.Context(), cardParams)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	results := make([]types.CardResponseSearchByOracleID, len(cards))

	for i, card := range cards {
		results[i] = types.CardResponseSearchByOracleID{
			ID:              card.ID,
			Name:            card.Name,
			FlavorName:      dbutils.FromNullString(card.FlavorName),
			ReleasedAt:      card.ReleaseDate,
			Set:             card.SetCode,
			SetName:         card.SetName,
			CollectorNumber: card.CollectorNumber,
		}
		if lang != "en" {
			results[i].Name = card.PrintedName.String
		}
	}

	legalitiesJSON, err := cfg.getLegalities(oracleCard.ID, c.Request.Context())
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	rulings, err := cfg.db.GetOracleRulings(c.Request.Context(), id)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	rulingsJSON := make([]types.ResponseRulings, len(rulings))

	for i, rule := range rulings {
		rulingsJSON[i] = types.ResponseRulings{
			OracleID:    rule.OracleID,
			Source:      dbutils.FromNullString(rule.Source),
			PublishedAt: rule.PublishedAt,
			Comment:     rule.Comment,
		}
	}

	oracleCardJSON.Rulings = rulingsJSON
	oracleCardJSON.Legalities = legalitiesJSON

	path := fmt.Sprintf("/cards/oracle/%v", idStr)
	base, _ := url.Parse(path)

	params := url.Values{}
	params.Add("limit", fmt.Sprintf("%d", limit))
	params.Add("lang", lang)
	params.Add("page", fmt.Sprintf("%d", page+1))

	base.RawQuery = params.Encode()

	c.JSON(http.StatusOK, gin.H{
		"page":           page,
		"number_pages":   numberPages,
		"number_results": numberResults,
		"next_page":      base.String(),
		"oracle_card":    oracleCardJSON,
		"results":        results,
	})
}

func (cfg *apiConfig) handlerCardByID(c *gin.Context) {
	idStr := c.Param("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "wrong UUID format"})
		return
	}

	card, err := cfg.db.GetCardByID(c.Request.Context(), id)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found CardID"})
		return
	} else if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	legalitiesJSON, err := cfg.getLegalities(card.ID, c.Request.Context())
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	cardJSON := types.ResponseCard{
		ID:         card.ID,
		OracleID:   dbutils.ToUUIDPtr(card.OracleID),
		Name:       card.Name,
		FlavorName: dbutils.FromNullString(card.FlavorName),
		Lang:       card.Lang,
		ReleasedAt: card.ReleaseDate,
		Layout:     card.Layout,
		ImageUris: types.ResponseImages{
			ImageNormal: dbutils.FromNullString(card.Image),
			ImagePNG:    dbutils.FromNullString(card.ImagePng),
			ImageLarge:  dbutils.FromNullString(card.ImageLarge),
			ImageSmall:  dbutils.FromNullString(card.ImageSmall),
			ImageCrop:   dbutils.FromNullString(card.ImageCrop),
		},
		ManaCost:        dbutils.FromNullString(card.ManaCost),
		Cmc:             card.Cmc,
		TypeLine:        card.TypeLine,
		OracleText:      dbutils.FromNullString(card.OracleText),
		Power:           dbutils.FromNullString(card.Power),
		Toughness:       dbutils.FromNullString(card.Toughness),
		Loyalty:         dbutils.FromNullString(card.Loyalty),
		Colors:          &card.Colors,
		ColorIdentity:   card.ColorIdentity,
		Defense:         dbutils.FromNullString(card.Defense),
		Keywords:        card.Keywords,
		FlavorText:      dbutils.FromNullString(card.FlavorText),
		Legalities:      legalitiesJSON,
		GameChanger:     dbutils.FromNullBool(card.GameChanger),
		Finishes:        card.Finishes,
		Set:             card.SetCode,
		SetName:         card.SetName,
		CollectorNumber: card.CollectorNumber,
		Rarity:          card.Rarity,
		Artist:          dbutils.FromNullString(card.Artist),
		EdhrecRank:      dbutils.FromNullInt32(card.EdhrecRank),
		Multifaced:      card.Multifaced,
	}

	multifaces := []database.CardFace{}

	if card.Multifaced {
		multifaces, err = cfg.db.GetCardFaces(c.Request.Context(), card.ID)
		if err != nil {
			c.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
			return
		}

		responseFaces := make([]types.ResponseCardFaces, len(multifaces))

		for i, face := range multifaces {
			responseFaces[i] = types.ResponseCardFaces{
				Name:       face.Name,
				ManaCost:   face.ManaCost,
				Cmc:        dbutils.FromNullFloat64(face.Cmc),
				Colors:     &face.Colors,
				TypeLine:   dbutils.FromNullString(face.TypeLine),
				OracleText: dbutils.FromNullString(face.OracleText),
				Power:      dbutils.FromNullString(face.Power),
				Toughness:  dbutils.FromNullString(face.Toughness),
				Loyalty:    dbutils.FromNullString(face.Loyalty),
				Defense:    dbutils.FromNullString(face.Defense),
				FlavorText: dbutils.FromNullString(face.FlavorText),
				Artist:     dbutils.FromNullString(face.Artist),
				Layout:     dbutils.FromNullString(face.Layout),
				ImageUris: types.ResponseImages{
					ImageNormal: dbutils.FromNullString(card.Image),
					ImagePNG:    dbutils.FromNullString(card.ImagePng),
					ImageLarge:  dbutils.FromNullString(card.ImageLarge),
					ImageSmall:  dbutils.FromNullString(card.ImageSmall),
					ImageCrop:   dbutils.FromNullString(card.ImageCrop),
				},
			}

			if card.Lang != "en" {
				responseFaces[i].Name = face.PrintedName.String
				responseFaces[i].TypeLine = dbutils.FromNullString(face.PrintedTypeLine)
				responseFaces[i].PrintedText = dbutils.FromNullString(face.PrintedText)
			}
		}

		cardJSON.CardFaces = &responseFaces
	}

	if card.Lang != "en" {
		cardJSON.Name = card.PrintedName.String
		cardJSON.TypeLine = card.PrintedTypeLine.String
		cardJSON.PrintedText = dbutils.FromNullString(card.PrintedText)
	}

	rulings, err := cfg.db.GetOracleRulings(c.Request.Context(), id)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	rulingsJSON := make([]types.ResponseRulings, len(rulings))

	for i, rule := range rulings {
		rulingsJSON[i] = types.ResponseRulings{
			OracleID:    rule.OracleID,
			Source:      dbutils.FromNullString(rule.Source),
			PublishedAt: rule.PublishedAt,
			Comment:     rule.Comment,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"card":    cardJSON,
		"rulings": rulingsJSON,
	})
}

func calculatePages(totalCards int64, pageSize int) int {
	if totalCards == 0 {
		return 0
	}

	return int(math.Ceil(float64(totalCards) / float64(pageSize)))
}

func (cfg *apiConfig) getLegalities(oracleID uuid.UUID, ctx context.Context) (types.Legalities, error) {
	legalities, err := cfg.db.GetCardLegalties(ctx, oracleID)
	if err == sql.ErrNoRows {
		return types.Legalities{}, nil
	} else if err != nil {
		return types.Legalities{}, err
	}

	return types.Legalities{
		Standard:        dbutils.FromNullString(legalities.Standard),
		Future:          dbutils.FromNullString(legalities.Future),
		Historic:        dbutils.FromNullString(legalities.Historic),
		Timeless:        dbutils.FromNullString(legalities.Timeless),
		Gladiator:       dbutils.FromNullString(legalities.Gladiator),
		Pioneer:         dbutils.FromNullString(legalities.Pioneer),
		Modern:          dbutils.FromNullString(legalities.Modern),
		Legacy:          dbutils.FromNullString(legalities.Legacy),
		Pauper:          dbutils.FromNullString(legalities.Pauper),
		Vintage:         dbutils.FromNullString(legalities.Vintage),
		Penny:           dbutils.FromNullString(legalities.Penny),
		Commander:       dbutils.FromNullString(legalities.Commander),
		Oathbreaker:     dbutils.FromNullString(legalities.Oathbreaker),
		Standardbrawl:   dbutils.FromNullString(legalities.Standardbrawl),
		Brawl:           dbutils.FromNullString(legalities.Brawl),
		Alchemy:         dbutils.FromNullString(legalities.Alchemy),
		Paupercommander: dbutils.FromNullString(legalities.Paupercommander),
		Duel:            dbutils.FromNullString(legalities.Duel),
		Oldschool:       dbutils.FromNullString(legalities.Oldschool),
		Premodern:       dbutils.FromNullString(legalities.Premodern),
		Predh:           dbutils.FromNullString(legalities.Predh),
	}, nil
}
