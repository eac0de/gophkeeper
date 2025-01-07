package main

// A simple program that opens the alternate screen buffer then counts down
// from 5 and then exits.

import (
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/eac0de/gophkeeper/client/internal/client"
	"github.com/eac0de/gophkeeper/client/internal/models/home"
	"github.com/eac0de/gophkeeper/client/internal/models/login"
	"github.com/eac0de/gophkeeper/client/internal/utils"
)

func main() {
	file, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("Не удалось открыть файл для логов: %v", err)
	}
	defer file.Close()

	// Установить вывод логов в файл
	log.SetOutput(file)

	var model tea.Model
	tokens := utils.LoadTokens()
	apiClient := client.NewAPIClient(
		"localhost:8081",
		"localhost:8080",
	)
	if tokens.AccessToken != "" || tokens.RefreshToken != "" {
		apiClient.Tokens = tokens
		model = home.New(apiClient)

	}
	if model == nil {
		model = login.New(apiClient, home.New(apiClient), "")
	}
	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
