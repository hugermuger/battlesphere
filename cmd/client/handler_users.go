package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"charm.land/bubbles/v2/textinput"
)

type Config struct {
	LastUsername string `json:"last_username"`
	Password     string `json:"password"`
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

	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&m.user)
	if err != nil {
		m.login.err = fmt.Sprintf("Internal error: %v", err)
		return
	}

	err = saveUserConfig(m.user.UserName, input[2].Value())
	if err != nil {
		m.login.err = fmt.Sprintf("Internal error: %v", err)
		return
	}

	m.login.loggedIn = true
	m.login.err = ""
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
	if res.StatusCode == http.StatusUnauthorized {
		error := error{}
		decoder := json.NewDecoder(res.Body)
		_ = decoder.Decode(&error)
		m.login.err = error.Error
		return
	}
	if res.StatusCode != http.StatusOK {
		error := error{}
		decoder := json.NewDecoder(res.Body)
		_ = decoder.Decode(&error)
		m.login.err = error.Error
		return
	}

	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&m.user)
	if err != nil {
		m.login.err = fmt.Sprintf("Internal error: %v", err)
		return
	}

	err = saveUserConfig(username, password)
	if err != nil {
		m.login.err = fmt.Sprintf("Internal error: %v", err)
		return
	}

	m.login.loggedIn = true
	m.login.err = ""
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

func saveUserConfig(username, password string) error {
	path, err := getConfigPath()
	if err != nil {
		return err
	}

	config := Config{LastUsername: username,
		Password: password}
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
