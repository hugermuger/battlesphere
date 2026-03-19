package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hugermuger/battlesphere/internal/auth"
	"github.com/hugermuger/battlesphere/internal/database"
)

type mockDB struct {
	database.Querier
	createUserFn        func(ctx context.Context, arg database.CreateUserParams) (database.User, error)
	getUserByUserNameFn func(ctx context.Context, userName string) (database.User, error)
}

func (m *mockDB) CreateUser(ctx context.Context, arg database.CreateUserParams) (database.User, error) {
	return m.createUserFn(ctx, arg)
}

func (m *mockDB) GetUserByUserName(ctx context.Context, userName string) (database.User, error) {
	return m.getUserByUserNameFn(ctx, userName)
}

func TestHandlerUsersCreate(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("successful user creation", func(t *testing.T) {
		mock := &mockDB{
			createUserFn: func(ctx context.Context, arg database.CreateUserParams) (database.User, error) {
				return database.User{
					ID:        uuid.New(),
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
					Email:     arg.Email,
					UserName:  arg.UserName,
				}, nil
			},
		}

		cfg := &apiConfig{db: mock}
		router := gin.New()
		router.POST("/users", cfg.handlerUsersCreate)

		body := map[string]string{
			"password":  "password123",
			"user_name": "testuser",
			"email":     "test@example.com",
		}
		jsonBody, _ := json.Marshal(body)

		req, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("expected status 201, got %d", w.Code)
		}

		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)

		if resp["user_name"] != "testuser" {
			t.Errorf("expected user_name testuser, got %v", resp["user_name"])
		}
	})
}

func TestHandlerLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	password := "password123"
	hashedPassword, _ := auth.HashPassword(password)

	t.Run("successful login", func(t *testing.T) {
		mock := &mockDB{
			getUserByUserNameFn: func(ctx context.Context, userName string) (database.User, error) {
				return database.User{
					ID:             uuid.New(),
					UserName:       userName,
					HashedPassword: hashedPassword,
				}, nil
			},
		}

		cfg := &apiConfig{db: mock}
		router := gin.New()
		router.POST("/login", cfg.handlerLogin)

		body := map[string]string{
			"password":  password,
			"user_name": "testuser",
		}
		jsonBody, _ := json.Marshal(body)

		req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}
	})

	t.Run("wrong password", func(t *testing.T) {
		mock := &mockDB{
			getUserByUserNameFn: func(ctx context.Context, userName string) (database.User, error) {
				return database.User{
					ID:             uuid.New(),
					UserName:       userName,
					HashedPassword: hashedPassword,
				}, nil
			},
		}

		cfg := &apiConfig{db: mock}
		router := gin.New()
		router.POST("/login", cfg.handlerLogin)

		body := map[string]string{
			"password":  "wrongpassword",
			"user_name": "testuser",
		}
		jsonBody, _ := json.Marshal(body)

		req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("expected status 401, got %d", w.Code)
		}
	})
}
