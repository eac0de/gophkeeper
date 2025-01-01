package main

// A simple program that opens the alternate screen buffer then counts down
// from 5 and then exits.

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/eac0de/gophkeeper/client/internal/client"
	"github.com/eac0de/gophkeeper/client/internal/models"
)

func main() {
	var model tea.Model
	tokens := models.LoadTokens()
	apiClient := client.NewAPIClient(
		"localhost:8081",
		"localhost:8080",
	)
	if tokens.AccessToken != "" || tokens.RefreshToken != "" {
		apiClient.Tokens = tokens
		model = models.InitialInputModel(apiClient)

	}
	if model == nil {
		model = models.InitialLoginModel(apiClient)
	}
	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
