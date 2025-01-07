package client

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/eac0de/gophkeeper/client/internal/schemes"
	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
)

func (c *APIClient) DownloadFile(itemID uuid.UUID) tea.Msg {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return schemes.DownloadFileMsg{
			Err: err,
		}
	}
	url := fmt.Sprintf("http://%s/api/gophkeeper/file_data/%s/download/", c.gophKeeperServerAddress, itemID.String())

	var resp *resty.Response
	var statusCode int
	var token string

	for i := 0; i < 2; i++ {
		token, err = c.getAccessToken()
		if err != nil {
			return schemes.DownloadFileMsg{
				Err:        err,
				StatusCode: http.StatusUnauthorized,
			}
		}
		resp, err = c.client.R().
			SetDoNotParseResponse(true). // Не обрабатывать тело ответа (оставляем как поток)
			SetAuthToken(token).
			Get(url)
		if err != nil {
			return schemes.DownloadFileMsg{
				Err: err,
			}
		}
		defer resp.RawBody().Close() // Закрываем поток, чтобы избежать утечек

		statusCode = resp.StatusCode()
		if statusCode == http.StatusUnauthorized {
			c.Tokens.AccessToken = ""
			continue
		}
		if statusCode != http.StatusOK {
			return schemes.DownloadFileMsg{
				Err:        getResponseError(resp),
				StatusCode: statusCode,
			}
		}
		contentDisposition := resp.Header().Get("Content-Disposition")
		fileName := extractFileName(contentDisposition)
		if fileName == "" {
			fileName = "downloaded_file"
		}
		dir := fmt.Sprintf("%s/GophkeeperFiles", homeDir)
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return schemes.DownloadFileMsg{
				Err: err,
			}
		}

		outFile, err := os.Create(dir + "/" + fileName)
		if err != nil {
			return schemes.DownloadFileMsg{
				Err: err,
			}
		}
		defer outFile.Close()
		_, err = outFile.ReadFrom(resp.RawBody())
		if err != nil {
			return schemes.DownloadFileMsg{
				Err: err,
			}
		}
		return schemes.DownloadFileMsg{
			StatusCode: statusCode,
		}
	}
	return schemes.DownloadFileMsg{
		Err:        getResponseError(resp),
		StatusCode: statusCode,
	}
}

// Функция для извлечения имени файла из Content-Disposition
func extractFileName(contentDisposition string) string {
	if contentDisposition == "" {
		return ""
	}

	// Ищем "filename=" в заголовке
	parts := strings.Split(contentDisposition, ";")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if strings.HasPrefix(part, "filename=") {
			// Убираем кавычки вокруг имени файла, если они есть
			return strings.Trim(part[len("filename="):], `"`)
		}
	}

	return ""
}
