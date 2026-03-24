package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"charm.land/bubbles/v2/textinput"
	"github.com/hugermuger/battlesphere/internal/types"
)

type Config struct {
	LastToken string `json:"last_token"`
}

func handlerCreateUser(input []textinput.Model, m *model) {
	if input[2].Value() != input[3].Value() {
		m.login.err = "Passwords do not match!"
		return
	}

	url := website + "/users"
	type parameters struct {
		Password string `json:"password"`
		UserName string `json:"user_name"`
		Email    string `json:"email"`
	}
	type error struct {
		Error string `json:"error"`
	}
	type response struct {
		types.User
	}

	params := parameters{
		Password: input[2].Value(),
		UserName: input[0].Value(),
		Email:    input[1].Value(),
	}

	jsonData, err := json.Marshal(params)
	if err != nil {
		m.login.err = fmt.Sprintf("Internal error: %v", err)
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		m.login.err = fmt.Sprintf("Internal error: %v", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		m.login.err = fmt.Sprintf("Internal error: %v", err)
		return
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		error := error{}
		decoder := json.NewDecoder(res.Body)
		_ = decoder.Decode(&error)
		m.login.err = error.Error
		return
	}

	resp := response{}

	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&resp)
	if err != nil {
		m.login.err = fmt.Sprintf("Internal error: %v", err)
		return
	}

	m.username = resp.User.UserName
	m.login.err = ""
	m.login.registerSucces = true
}

func handlerLogin(password, username string, m *model) {
	url := website + "/login"
	type parameters struct {
		Password string `json:"password"`
		UserName string `json:"user_name"`
	}
	type error struct {
		Error string `json:"error"`
	}
	type response struct {
		types.User
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	params := parameters{
		Password: password,
		UserName: username,
	}

	jsonData, err := json.Marshal(params)
	if err != nil {
		m.login.err = fmt.Sprintf("Internal error: %v", err)
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		m.login.err = fmt.Sprintf("Internal error: %v", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		m.login.err = fmt.Sprintf("Internal error: %v", err)
		return
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		error := error{}
		decoder := json.NewDecoder(res.Body)
		_ = decoder.Decode(&error)
		m.login.err = error.Error
		return
	}

	resp := response{}

	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&resp)
	if err != nil {
		m.login.err = fmt.Sprintf("Internal error: %v", err)
		return
	}

	m.username = resp.User.UserName

	err = saveUserConfig(username, resp.RefreshToken)
	if err != nil {
		m.login.err = fmt.Sprintf("Internal error: %v", err)
		return
	}

	m.jwtToken = resp.Token
	m.login.loggedIn = true
	m.login.err = ""
}

func handlerRefresh(token string, m *model) {
	url := website + "/refresh"

	type response struct {
		Token    string `json:"token"`
		UserName string `json:"user_name"`
	}
	type error struct {
		Error string `json:"error"`
	}

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		m.login.err = fmt.Sprintf("Internal error: %v", err)
		return
	}

	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		m.login.err = fmt.Sprintf("Internal error: %v", err)
		return
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		error := error{}
		decoder := json.NewDecoder(res.Body)
		_ = decoder.Decode(&error)
		m.login.err = error.Error
		return
	}

	resp := response{}

	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&resp)
	if err != nil {
		m.login.err = fmt.Sprintf("Internal error: %v", err)
		return
	}

	m.username = resp.UserName
	m.jwtToken = resp.Token
	m.login.loggedIn = true
}

func handlerLogout(token string, m *model) {
	url := website + "/revoke"

	type error struct {
		Error string `json:"error"`
	}

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		m.login.err = fmt.Sprintf("Internal error: %v", err)
		return
	}

	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		m.login.err = fmt.Sprintf("Internal error: %v", err)
		return
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusNoContent {
		error := error{}
		decoder := json.NewDecoder(res.Body)
		_ = decoder.Decode(&error)
		m.login.err = error.Error
	}

	m.username = ""
	m.login.loggedIn = false
	cleanLoginInput(m)
	_ = saveUserConfig("", "")
}

func cleanLoginInput(m *model) {
	for i, _ := range m.login.loginInput {
		m.login.loginInput[i].Reset()
	}
	for i, _ := range m.login.registerInput {
		m.login.registerInput[i].Reset()
	}
}

func getConfigPath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	appDir := filepath.Join(dir, "battlesphere")
	err = os.MkdirAll(appDir, 0755)

	return filepath.Join(appDir, "config.json"), err
}

func saveUserConfig(username, token string) error {
	path, err := getConfigPath()
	if err != nil {
		return err
	}

	config := Config{
		LastToken: token}
	data, _ := json.MarshalIndent(config, "", "  ")

	return os.WriteFile(path, data, 0644)
}

func loadUserConfig() (Config, error) {
	path, err := getConfigPath()
	if err != nil {
		return Config{}, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	var config Config
	json.Unmarshal(data, &config)
	return config, nil
}
