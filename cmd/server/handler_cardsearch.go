package main

import (
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
				ManaCost: &card.ManaCost.String,
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
				ManaCost: &card.ManaCost.String,
				TypeLine: card.PrintedTypeLine.String,
			}
		}

		results = resultsBuffer
	}

	next_page := ""
	if numberPages > 1 && page < numberPages {
		next_page = getNextPageURL(name, limit, lang, page)
	}

	c.JSON(http.StatusOK, gin.H{
		"page":           page,
		"number_pages":   numberPages,
		"results":        results,
		"number_results": numberResults,
		"next_page":      next_page,
	})
}

func getNextPageURL(name string, limit int, lang string, page int) string {
	base, _ := url.Parse("/cards/search")

	params := url.Values{}
	params.Add("name", name)
	params.Add("limit", fmt.Sprintf("%d", limit))
	params.Add("lang", lang)
	params.Add("page", fmt.Sprintf("%d", page+1))

	base.RawQuery = params.Encode()

	return base.String()
}

func calculatePages(totalCards int64, pageSize int) int {
	if totalCards == 0 {
		return 0
	}

	return int(math.Ceil(float64(totalCards) / float64(pageSize)))
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
			Source:      &rule.Source.String,
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
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	multifaced := true

	multifaces, err := cfg.db.GetCardFaces(c.Request.Context(), oracleCard.ID)
	if err == sql.ErrNoRows {
		multifaced = false
	} else if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	oracleCardJSON := types.ResponseByOracleID{}

	if lang == "en" {
		if multifaced {
			responseFaces := make([]types.SingleCardFacesByOracleID, len(multifaces))

			for i, face := range multifaces {
				responseFaces[i] = types.SingleCardFacesByOracleID{
					Name:       face.Name,
					ManaCost:   face.ManaCost,
					TypeLine:   &face.TypeLine.String,
					OracleText: &face.OracleText.String,
					Power:      &face.Power.String,
					Toughness:  &face.Toughness.String,
					Loyalty:    &face.Loyalty.String,
					Defense:    &face.Defense.String,
				}
			}

			oracleCardJSON = types.ResponseByOracleID{
				Name:       oracleCard.Name,
				Layout:     oracleCard.Layout,
				ManaCost:   &oracleCard.ManaCost.String,
				TypeLine:   oracleCard.TypeLine,
				OracleText: &oracleCard.OracleText.String,
				Power:      &oracleCard.Power.String,
				Toughness:  &oracleCard.Toughness.String,
				Loyalty:    &oracleCard.Loyalty.String,
				Defense:    &oracleCard.Defense.String,
				CardFaces:  &responseFaces,
				Multifaced: true,
			}

		} else {
			oracleCardJSON = types.ResponseByOracleID{
				Name:       oracleCard.Name,
				Layout:     oracleCard.Layout,
				ManaCost:   &oracleCard.ManaCost.String,
				TypeLine:   oracleCard.TypeLine,
				OracleText: &oracleCard.OracleText.String,
				Power:      &oracleCard.Power.String,
				Toughness:  &oracleCard.Toughness.String,
				Loyalty:    &oracleCard.Loyalty.String,
				Defense:    &oracleCard.Defense.String,
				Multifaced: false,
			}
		}
	} else {
		if multifaced {
			responseFaces := make([]types.SingleCardFacesByOracleID, len(multifaces))

			for i, face := range multifaces {
				responseFaces[i] = types.SingleCardFacesByOracleID{
					Name:       face.PrintedName.String,
					ManaCost:   face.ManaCost,
					TypeLine:   &face.PrintedTypeLine.String,
					OracleText: &face.PrintedText.String,
					Power:      &face.Power.String,
					Toughness:  &face.Toughness.String,
					Loyalty:    &face.Loyalty.String,
					Defense:    &face.Defense.String,
				}
			}

			oracleCardJSON = types.ResponseByOracleID{
				Name:       oracleCard.PrintedName.String,
				Layout:     oracleCard.Layout,
				ManaCost:   &oracleCard.ManaCost.String,
				TypeLine:   oracleCard.PrintedTypeLine.String,
				OracleText: &oracleCard.PrintedText.String,
				Power:      &oracleCard.Power.String,
				Toughness:  &oracleCard.Toughness.String,
				Loyalty:    &oracleCard.Loyalty.String,
				Defense:    &oracleCard.Defense.String,
				CardFaces:  &responseFaces,
				Multifaced: true,
			}

		} else {
			oracleCardJSON = types.ResponseByOracleID{
				Name:       oracleCard.PrintedName.String,
				Layout:     oracleCard.Layout,
				ManaCost:   &oracleCard.ManaCost.String,
				TypeLine:   oracleCard.PrintedTypeLine.String,
				OracleText: &oracleCard.PrintedText.String,
				Power:      &oracleCard.Power.String,
				Toughness:  &oracleCard.Toughness.String,
				Loyalty:    &oracleCard.Loyalty.String,
				Defense:    &oracleCard.Defense.String,
				Multifaced: false,
			}
		}
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
		if lang == "en" {
			results[i] = types.CardResponseSearchByOracleID{
				ID:              card.ID,
				Name:            card.Name,
				FlavorName:      &card.FlavorName.String,
				ReleasedAt:      card.ReleaseDate,
				Set:             card.SetCode,
				SetName:         card.SetName,
				CollectorNumber: card.CollectorNumber,
			}
		} else {
			results[i] = types.CardResponseSearchByOracleID{
				ID:              card.ID,
				Name:            card.PrintedName.String,
				FlavorName:      &card.FlavorName.String,
				ReleasedAt:      card.ReleaseDate,
				Set:             card.SetCode,
				SetName:         card.SetName,
				CollectorNumber: card.CollectorNumber,
			}
		}
	}

	legalities, err := cfg.db.GetCardLegalties(c.Request.Context(), oracleCard.ID)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	legalitiesJSON := types.Legalities{
		Standard:        &legalities.Standard.String,
		Future:          &legalities.Future.String,
		Historic:        &legalities.Historic.String,
		Timeless:        &legalities.Timeless.String,
		Gladiator:       &legalities.Gladiator.String,
		Pioneer:         &legalities.Pioneer.String,
		Modern:          &legalities.Modern.String,
		Legacy:          &legalities.Legacy.String,
		Pauper:          &legalities.Pauper.String,
		Vintage:         &legalities.Vintage.String,
		Penny:           &legalities.Penny.String,
		Commander:       &legalities.Commander.String,
		Oathbreaker:     &legalities.Oathbreaker.String,
		Standardbrawl:   &legalities.Standardbrawl.String,
		Brawl:           &legalities.Brawl.String,
		Alchemy:         &legalities.Alchemy.String,
		Paupercommander: &legalities.Paupercommander.String,
		Duel:            &legalities.Duel.String,
		Oldschool:       &legalities.Oldschool.String,
		Premodern:       &legalities.Premodern.String,
		Predh:           &legalities.Predh.String,
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
			Source:      &rule.Source.String,
			PublishedAt: rule.PublishedAt,
			Comment:     rule.Comment,
		}
	}

	path := fmt.Sprintf("cards/oracle/%v", idStr)
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
		"rulings":        rulingsJSON,
		"legalities":     legalitiesJSON,
		"results":        results,
	})
}
