package main

import (
	"database/sql"
	"net/http"
	"net/mail"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hugermuger/battlesphere/internal/auth"
	"github.com/hugermuger/battlesphere/internal/database"
	"github.com/hugermuger/battlesphere/internal/types"
)

func (cfg *apiConfig) handlerUsersCreate(c *gin.Context) {
	type parameters struct {
		Password string `json:"password"`
		UserName string `json:"user_name"`
		Email    string `json:"email"`
	}

	type response struct {
		types.User
	}

	params := parameters{}

	if err := c.BindJSON(&params); err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't decode parameters"})
		return
	}

	if params.Email == "" || params.Password == "" || params.UserName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Empty Parameter"})
		return
	}

	if len(params.Password) < 8 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password must be 8 characters long"})
		return
	}

	_, err := mail.ParseAddress(params.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Mail in wrong format"})
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't hash password"})
		return
	}

	user, err := cfg.db.CreateUser(c.Request.Context(), database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hashedPassword,
		UserName:       params.UserName,
	})
	if err == sql.ErrNoRows {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username or Mail does already exist"})
		return
	} else if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't create user"})
		return
	}

	c.JSON(http.StatusCreated, response{
		User: types.User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
			UserName:  user.UserName,
		},
	})
}

func (cfg *apiConfig) handlerLogin(c *gin.Context) {
	type parameters struct {
		Password string `json:"password"`
		UserName string `json:"user_name"`
	}
	type response struct {
		types.User
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	params := parameters{}

	if err := c.BindJSON(&params); err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't decode parameters"})
		return
	}

	user, err := cfg.db.GetUserByUserName(c.Request.Context(), params.UserName)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unknown Username"})
		return
	}

	exists, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil || !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Wrong password or username"})
		return
	}

	accessToken, err := auth.MakeJWT(
		user.ID,
		cfg.jwtSecret,
		time.Hour,
	)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't create access JWT"})
		return
	}

	refreshToken := auth.MakeRefreshToken()

	_, err = cfg.db.CreateRefreshToken(c.Request.Context(), database.CreateRefreshTokenParams{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: time.Now().UTC().Add(time.Hour * 24 * 60),
	})
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't save refresh token"})
		return
	}

	c.JSON(http.StatusOK, response{
		User: types.User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
			UserName:  user.UserName,
		},
		Token:        accessToken,
		RefreshToken: refreshToken,
	})
}

func (cfg *apiConfig) handlerRefresh(c *gin.Context) {
	type response struct {
		Token    string `json:"token"`
		UserName string `json:"user_name"`
	}

	refreshToken, err := auth.GetBearerToken(c.Request.Header)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Couldn't find token"})
		return
	}

	user, err := cfg.db.GetUserFromRefreshToken(c.Request.Context(), refreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Couldn't get user for refresh token"})
		return
	}

	accessToken, err := auth.MakeJWT(
		user.ID,
		cfg.jwtSecret,
		time.Hour,
	)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Couldn't validate token"})
		return
	}

	c.JSON(http.StatusOK, response{
		Token:    accessToken,
		UserName: user.UserName,
	})
}

func (cfg *apiConfig) handlerRevoke(c *gin.Context) {
	refreshToken, err := auth.GetBearerToken(c.Request.Header)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Couldn't find token"})
		return
	}

	_, err = cfg.db.RevokeRefreshToken(c.Request.Context(), refreshToken)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't revoke session"})
		return
	}

	c.Status(http.StatusNoContent)
}
