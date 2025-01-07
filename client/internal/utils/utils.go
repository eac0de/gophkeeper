package utils

import (
	"encoding/json"
	"errors"
	"log"
	"os"

	"github.com/eac0de/gophkeeper/client/internal/schemes"
)

func SaveTokens(tokens schemes.Tokens) {
	file, err := os.OpenFile(os.TempDir()+"/gophkeeper_auth.json", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	err = json.NewEncoder(file).Encode(schemes.Tokens{AccessToken: tokens.AccessToken, RefreshToken: tokens.RefreshToken})
	if err != nil {
		log.Fatal(err)
	}
}

func LoadTokens() schemes.Tokens {
	file, err := os.Open(os.TempDir() + "/gophkeeper_auth.json")
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return schemes.Tokens{}
		}
		log.Fatal(err)
	}
	defer file.Close()
	var tokens schemes.Tokens
	err = json.NewDecoder(file).Decode(&tokens)
	if err != nil {
		log.Fatal(err)
	}
	return tokens
}
