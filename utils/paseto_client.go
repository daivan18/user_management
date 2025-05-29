package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

type TokenRequest struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

func GeneratePasetoToken(authURL string, id int, username string) (string, error) {
	reqBody, _ := json.Marshal(TokenRequest{
		ID:       id,
		Username: username,
	})

	resp, err := http.Post(authURL, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result map[string]string
	json.NewDecoder(resp.Body).Decode(&result)

	if token, ok := result["token"]; ok {
		return token, nil
	}
	return "", errors.New("no token in response")
}
