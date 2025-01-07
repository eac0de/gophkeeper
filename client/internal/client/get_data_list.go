package client

import (
	"fmt"
	"net/http"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/eac0de/gophkeeper/client/internal/schemes"
	"github.com/go-resty/resty/v2"
)

func (c *APIClient) GetUserTextDataList(offset string) tea.Msg {
	var resp *resty.Response
	var statusCode int
	var err error
	var token string
	var responseBody []schemes.UserTextData

	for i := 0; i < 2; i++ {
		token, err = c.getAccessToken()
		if err != nil {
			return schemes.GetUserTextDataListMsg{
				Err:        err,
				StatusCode: http.StatusUnauthorized,
			}
		}

		req := c.client.R().SetAuthToken(token).SetResult(&responseBody)
		if offset != "" {
			req.SetQueryParam("offset", offset)
		}
		resp, err = req.Get(fmt.Sprintf("http://%s/api/gophkeeper/text_data/", c.gophKeeperServerAddress))
		if err != nil {
			return schemes.GetUserTextDataListMsg{
				Err: err,
			}
		}
		statusCode = resp.StatusCode()
		if statusCode == http.StatusUnauthorized {
			c.Tokens.AccessToken = ""
			continue
		}
		if statusCode == http.StatusOK {
			return schemes.GetUserTextDataListMsg{
				Err:        nil,
				StatusCode: statusCode,
				List:       responseBody,
			}
		}
		break
	}
	return schemes.GetUserTextDataListMsg{
		Err:        getResponseError(resp),
		StatusCode: statusCode,
	}
}

func (c *APIClient) GetUserFileDataList(offset string) tea.Msg {
	var resp *resty.Response
	var statusCode int
	var err error
	var token string
	var responseBody []schemes.UserFileData

	for i := 0; i < 2; i++ {
		token, err = c.getAccessToken()
		if err != nil {
			return schemes.GetUserFileDataListMsg{
				Err:        err,
				StatusCode: http.StatusUnauthorized,
			}
		}
		req := c.client.R().SetAuthToken(token).SetResult(&responseBody)
		if offset != "" {
			req.SetQueryParam("offset", offset)
		}
		resp, err = req.Get(fmt.Sprintf("http://%s/api/gophkeeper/file_data/", c.gophKeeperServerAddress))
		if err != nil {
			return schemes.GetUserFileDataListMsg{
				Err: err,
			}
		}
		statusCode = resp.StatusCode()
		if statusCode == http.StatusUnauthorized {
			c.Tokens.AccessToken = ""
			continue
		}
		if statusCode == http.StatusOK {
			return schemes.GetUserFileDataListMsg{
				Err:        nil,
				StatusCode: statusCode,
				List:       responseBody,
			}
		}
		break
	}
	return schemes.GetUserFileDataListMsg{
		Err:        getResponseError(resp),
		StatusCode: statusCode,
	}
}

func (c *APIClient) GetUserAuthInfoList(offset string) tea.Msg {
	var resp *resty.Response
	var statusCode int
	var err error
	var token string
	var responseBody []schemes.UserAuthInfo

	for i := 0; i < 2; i++ {
		token, err = c.getAccessToken()
		if err != nil {
			return schemes.GetUserAuthInfoListMsg{
				Err:        err,
				StatusCode: http.StatusUnauthorized,
			}
		}
		req := c.client.R().SetAuthToken(token).SetResult(&responseBody)
		if offset != "" {
			req.SetQueryParam("offset", offset)
		}
		resp, err = req.Get(fmt.Sprintf("http://%s/api/gophkeeper/auth_info/", c.gophKeeperServerAddress))
		if err != nil {
			return schemes.GetUserAuthInfoListMsg{
				Err: err,
			}
		}
		statusCode = resp.StatusCode()
		if statusCode == http.StatusUnauthorized {
			c.Tokens.AccessToken = ""
			continue
		}
		if statusCode == http.StatusOK {
			return schemes.GetUserAuthInfoListMsg{
				Err:        err,
				List:       responseBody,
				StatusCode: statusCode,
			}

		}
		break
	}
	return schemes.GetUserAuthInfoListMsg{
		Err:        getResponseError(resp),
		StatusCode: statusCode,
	}
}

func (c *APIClient) GetUserBankCardList(offset string) tea.Msg {
	var resp *resty.Response
	var statusCode int
	var err error
	var token string
	var responseBody []schemes.UserBankCard

	for i := 0; i < 2; i++ {
		token, err = c.getAccessToken()
		if err != nil {
			return schemes.GetUserBankCardListMsg{
				Err:        err,
				StatusCode: http.StatusUnauthorized,
			}
		}
		req := c.client.R().SetAuthToken(token).SetResult(&responseBody)
		if offset != "" {
			req.SetQueryParam("offset", offset)
		}
		resp, err = req.Get(fmt.Sprintf("http://%s/api/gophkeeper/bank_cards/", c.gophKeeperServerAddress))
		if err != nil {
			return schemes.GetUserBankCardListMsg{
				Err: err,
			}
		}
		statusCode = resp.StatusCode()
		if statusCode == http.StatusUnauthorized {
			c.Tokens.AccessToken = ""
			continue
		}
		if statusCode == http.StatusOK {
			return schemes.GetUserBankCardListMsg{
				Err:        nil,
				StatusCode: statusCode,
				List:       responseBody,
			}
		}
		break
	}
	return schemes.GetUserBankCardListMsg{
		Err:        getResponseError(resp),
		StatusCode: statusCode,
	}

}
