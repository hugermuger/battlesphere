package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hugermuger/battlesphere/internal/types"
)

func handlerCreateUser(email, userName, password string) (types.User, error) {
	url := website + "/users"
	type parameters struct {
		Password string `json:"password"`
		UserName string `json:"user_name"`
		Email    string `json:"email"`
	}
	params := parameters{
		Password: password,
		UserName: userName,
		Email:    email,
	}

	jsonData, err := json.Marshal(params)
	if err != nil {
		return types.User{}, fmt.Errorf("Internal error: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return types.User{}, fmt.Errorf("Internal error: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return types.User{}, fmt.Errorf("Internal error: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		return types.User{}, fmt.Errorf("expected status 201, got %d. Error: %s", res.StatusCode, res.Header.Get("error"))
	}

	user := types.User{}
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&user)
	if err != nil {
		return types.User{}, fmt.Errorf("Internal error: %v", err)
	}

	return user, nil
}

func handlerLogin(userName, password string) (types.User, error) {
	url := website + "/login"
	type parameters struct {
		Password string `json:"password"`
		UserName string `json:"user_name"`
	}
	params := parameters{
		Password: password,
		UserName: userName,
	}

	jsonData, err := json.Marshal(params)
	if err != nil {
		return types.User{}, fmt.Errorf("Internal error: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return types.User{}, fmt.Errorf("Internal error: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return types.User{}, fmt.Errorf("Internal error: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode == http.StatusUnauthorized {
		return types.User{}, fmt.Errorf("%v", res.Header.Get("error"))
	}
	if res.StatusCode != http.StatusOK {
		return types.User{}, fmt.Errorf("expected status 200, got %d. Error: %s", res.StatusCode, res.Header.Get("error"))
	}

	user := types.User{}
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&user)
	if err != nil {
		return types.User{}, fmt.Errorf("Internal error: %v", err)
	}

	return user, nil
}
