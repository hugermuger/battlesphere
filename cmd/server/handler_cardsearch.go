package main

import (
	"database/sql"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hugermuger/battlesphere/internal/database"
	"github.com/hugermuger/battlesphere/internal/dbutils"
	"github.com/hugermuger/battlesphere/internal/scryfall"
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
			c.JSON(http.StatusBadRequest, gin.H{"error": "database error"})
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
			c.JSON(http.StatusBadRequest, gin.H{"error": "database error"})
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

	results := []scryfall.CardResponseList{}

	if lang == "en" {
		cardArgs := database.SearchCardsByNameListEngParams{
			Column1: dbutils.ToNullString(&name),
			Limit:   int32(limit),
			Offset:  int32(offset),
		}

		cards, err := cfg.db.SearchCardsByNameListEng(c.Request.Context(), cardArgs)
		if err != nil {
			c.Error(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "database error"})
			return
		}

		resultsBuffer := make([]scryfall.CardResponseList, len(cards))
		for i, card := range cards {
			resultsBuffer[i] = scryfall.CardResponseList{
				Name:    card.Name,
				SetName: card.SetName,
				SetCode: card.SetCode,
				Rarity:  card.Rarity,
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
			c.JSON(http.StatusBadRequest, gin.H{"error": "database error"})
			return
		}

		resultsBuffer := make([]scryfall.CardResponseList, len(cards))
		for i, card := range cards {
			resultsBuffer[i] = scryfall.CardResponseList{
				Name:    card.PrintedName.String,
				SetName: card.SetName,
				SetCode: card.SetCode,
				Rarity:  card.Rarity,
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
