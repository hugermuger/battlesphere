package main

import (
	"net/http"

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
	if err != nil {
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
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't check hashed password"})
		return
	}
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Wrong password"})
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
	})
}
