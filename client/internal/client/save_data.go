package client

import (
	"fmt"
	"net/http"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/eac0de/gophkeeper/client/internal/schemes"
	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
)

func (c *APIClient) SaveUserTextData(id uuid.UUID, name, textData string, metadata map[string]interface{}) tea.Msg {
	var resp *resty.Response
	var statusCode int
	var err error
	var url, method, token string
	if id != uuid.Nil {
		url = fmt.Sprintf("http://%s/api/gophkeeper/text_data/%s/", c.gophKeeperServerAddress, id.String())
		method = resty.MethodPut
	} else {
		url = fmt.Sprintf("http://%s/api/gophkeeper/text_data/", c.gophKeeperServerAddress)
		method = resty.MethodPost
	}
	requestBody := struct {
		Name     string                 `json:"name"`
		TextData string                 `json:"text_data"`
		Metadata map[string]interface{} `json:"metadata"`
	}{
		Name:     name,
		TextData: textData,
		Metadata: metadata,
	}
	var responseBody schemes.UserTextData
	for i := 0; i < 2; i++ {
		token, err = c.getAccessToken()
		if err != nil {
			return schemes.SaveTextDataMsg{
				Err:        err,
				StatusCode: http.StatusUnauthorized,
			}
		}
		resp, err = c.client.R().SetAuthToken(token).SetBody(requestBody).SetResult(&responseBody).Execute(method, url)
		if err != nil {
			return schemes.SaveTextDataMsg{
				Err: err,
			}
		}
		statusCode = resp.StatusCode()
		if statusCode == http.StatusUnauthorized {
			c.Tokens.AccessToken = ""
			continue
		}
		if statusCode == http.StatusOK || statusCode == http.StatusCreated {
			return schemes.SaveTextDataMsg{
				StatusCode: statusCode,
				TextData:   responseBody,
			}
		}
		break
	}
	return schemes.SaveTextDataMsg{
		Err:        getResponseError(resp),
		StatusCode: statusCode,
	}

}

func (c *APIClient) SaveUserFileData(id uuid.UUID, name, filePath string, metadata map[string]interface{}) tea.Msg {
	var resp *resty.Response
	var statusCode int
	var err error
	var url, method, token string
	req := c.client.R()
	if id != uuid.Nil {
		requestBody := struct {
			Name     string                 `json:"name"`
			Metadata map[string]interface{} `json:"metadata"`
		}{
			Name:     name,
			Metadata: metadata,
		}
		url = fmt.Sprintf("http://%s/api/gophkeeper/file_data/%s/", c.gophKeeperServerAddress, id.String())
		method = resty.MethodPut
		req.SetBody(requestBody)
	} else {
		url = fmt.Sprintf("http://%s/api/gophkeeper/file_data/", c.gophKeeperServerAddress)
		method = resty.MethodPost
		file, err := os.OpenFile(filePath, os.O_RDONLY, 0644)
		if err != nil {
			return schemes.SaveFileDataMsg{
				Err: err,
			}
		}
		req.SetFileReader("file", file.Name(), file)
	}
	var responseBody schemes.UserFileData
	req.SetResult(&responseBody)
	for i := 0; i < 2; i++ {
		token, err = c.getAccessToken()
		if err != nil {
			return schemes.SaveFileDataMsg{
				Err:        err,
				StatusCode: http.StatusUnauthorized,
			}
		}
		resp, err = req.SetAuthToken(token).Execute(method, url)
		if err != nil {
			return schemes.SaveFileDataMsg{
				Err: err,
			}
		}
		statusCode = resp.StatusCode()
		if statusCode == http.StatusUnauthorized {
			c.Tokens.AccessToken = ""
			continue
		}
		if statusCode == http.StatusOK || statusCode == http.StatusCreated {
			return schemes.SaveFileDataMsg{
				StatusCode: statusCode,
				FileData:   responseBody,
			}
		}
		break
	}
	return schemes.SaveFileDataMsg{
		Err:        getResponseError(resp),
		StatusCode: statusCode,
	}
}
func (c *APIClient) SaveUserBankCard(id uuid.UUID, name, number, cardHolder, expireDate, csc string, metadata map[string]interface{}) tea.Msg {
	var resp *resty.Response
	var statusCode int
	var err error
	var url, method, token string
	if id != uuid.Nil {
		url = fmt.Sprintf("http://%s/api/gophkeeper/bank_cards/%s/", c.gophKeeperServerAddress, id.String())
		method = resty.MethodPut
	} else {
		url = fmt.Sprintf("http://%s/api/gophkeeper/bank_cards/", c.gophKeeperServerAddress)
		method = resty.MethodPost
	}
	requestBody := struct {
		Name       string                 `json:"name"`
		Number     string                 `json:"number"`
		CardHolder string                 `json:"card_holder"`
		ExpireDate string                 `json:"expire_date"`
		CSC        string                 `json:"csc"`
		Metadata   map[string]interface{} `json:"metadata"`
	}{
		Name:       name,
		Number:     number,
		CardHolder: cardHolder,
		ExpireDate: expireDate,
		CSC:        csc,
		Metadata:   metadata,
	}
	var responseBody schemes.UserBankCard
	for i := 0; i < 2; i++ {
		token, err = c.getAccessToken()
		if err != nil {
			return schemes.SaveBankCardMsg{
				Err:        err,
				StatusCode: http.StatusUnauthorized,
			}
		}
		resp, err = c.client.R().SetAuthToken(token).SetBody(requestBody).SetResult(&responseBody).Execute(method, url)
		if err != nil {
			return schemes.SaveBankCardMsg{
				Err: err,
			}
		}
		statusCode = resp.StatusCode()
		if statusCode == http.StatusUnauthorized {
			c.Tokens.AccessToken = ""
			continue
		}
		if statusCode == http.StatusOK || statusCode == http.StatusCreated {
			return schemes.SaveBankCardMsg{
				StatusCode: statusCode,
				BankCard:   responseBody,
			}
		}
		break
	}
	return schemes.SaveBankCardMsg{
		Err:        getResponseError(resp),
		StatusCode: statusCode,
	}
}

func (c *APIClient) SaveUserAuthInfo(id uuid.UUID, name, login, password string, metadata map[string]interface{}) tea.Msg {
	var resp *resty.Response
	var statusCode int
	var err error
	var url, method, token string
	if id != uuid.Nil {
		url = fmt.Sprintf("http://%s/api/gophkeeper/auth_info/%s/", c.gophKeeperServerAddress, id.String())
		method = resty.MethodPut
	} else {
		url = fmt.Sprintf("http://%s/api/gophkeeper/auth_info/", c.gophKeeperServerAddress)
		method = resty.MethodPost
	}

	requestBody := struct {
		Name     string                 `json:"name"`
		Login    string                 `json:"login"`
		Password string                 `json:"password"`
		Metadata map[string]interface{} `json:"metadata"`
	}{
		Name:     name,
		Login:    login,
		Password: password,
		Metadata: metadata,
	}
	var responseBody schemes.UserAuthInfo
	for i := 0; i < 2; i++ {
		token, err = c.getAccessToken()
		if err != nil {
			return schemes.SaveTextDataMsg{
				Err:        err,
				StatusCode: http.StatusUnauthorized,
			}
		}
		resp, err = c.client.R().SetAuthToken(token).SetBody(requestBody).SetResult(&responseBody).Execute(method, url)
		if err != nil {
			return schemes.SaveAuthInfoMsg{
				Err: err,
			}
		}
		statusCode = resp.StatusCode()
		if statusCode == http.StatusUnauthorized {
			c.Tokens.AccessToken = ""
			continue
		}
		if statusCode == http.StatusOK || statusCode == http.StatusCreated {
			return schemes.SaveAuthInfoMsg{
				StatusCode: statusCode,
				AuthInfo:   responseBody,
			}
		}
		break
	}
	return schemes.SaveAuthInfoMsg{
		Err:        getResponseError(resp),
		StatusCode: statusCode,
	}
}
