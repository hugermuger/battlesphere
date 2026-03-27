package main

import (
	"encoding/csv"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hugermuger/battlesphere/internal/auth"
	"github.com/hugermuger/battlesphere/internal/database"
)

func (cfg *apiConfig) handlerImportCollection(c *gin.Context) {
	type response struct {
		Missing [][]string `json:"missing"`
	}

	formatName := c.Query("format")

	token, err := auth.GetBearerToken(c.Request.Header)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Couldn't find JWT"})
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Couldn't validate JWT"})
		return
	}

	if formatName != "dragon_shield" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Fromat not supported"})
		return
	}

	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File upload failed"})
		return
	}
	defer file.Close()

	tx, err := cfg.dbConn.Begin()
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't connect to DB"})
		return
	}
	defer tx.Rollback()

	reader := csv.NewReader(file)
	if _, err := reader.Read(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Couldn't read csv"})
		return
	}
	if _, err := reader.Read(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Couldn't read csv"})
		return
	}

	unknown := [][]string{}
	qtx := cfg.db.WithTx(tx)

	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Couldn't read csv"})
			return
		}

		params := database.AddCardToCollectionsParams{}

		if formatName == "dragon_shield" {
			params, err = mapDragonShield(line, cfg, c.Request.Context())
			if err != nil {
				unknown = append(unknown, line)
			}
		}

		params.UserID = userID

		err = qtx.AddCardToCollections(c.Request.Context(), params)
		if err != nil {
			c.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't add card to DB"})
			return
		}
	}

	err = tx.Commit()
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't add collection to DB"})
		return
	}

	if len(unknown) == 0 {
		c.Status(http.StatusCreated)
		return
	} else {
		c.JSON(http.StatusCreated, response{
			Missing: unknown,
		})
	}
}
