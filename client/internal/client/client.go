package client

import (
	"errors"
	"fmt"
	"strings"

	"github.com/eac0de/gophkeeper/client/internal/schemes"
	"github.com/go-resty/resty/v2"
)

type APIClient struct {
	gophKeeperServerAddress string
	authServerAddress       string
	client                  *resty.Client
}

func NewAPIClient(gophkeeperServerAddress string, authServerAddress string) *APIClient {
	return &APIClient{
		gophKeeperServerAddress: gophkeeperServerAddress,
		authServerAddress:       authServerAddress,
		client:                  resty.New(),
	}
}

func (c *APIClient) GenerateEmailCode(email string) (string, error) {
	var result struct {
		EmailCodeId string `json:"email_code_id"`
	}
	resp, err := c.client.R().
		SetResult(&result).
		SetBody(map[string]string{"email": email}).
		Post(fmt.Sprintf("http://%s/api/auth/code/generate/", c.authServerAddress))
	if err != nil {
		return "", err
	}
	if resp.StatusCode() != 201 {
		return "", errors.New("Response error: " + resp.String() + " Status: " + resp.Status())
	}
	return result.EmailCodeId, nil
}

func (c *APIClient) VerifyEmailCode(emailCodeId string, code int) (int, *schemes.Tokens, error) {
	requestBody := struct {
		EmailCodeId string `json:"email_code_id"`
		Code        int    `json:"code"`
	}{
		EmailCodeId: emailCodeId,
		Code:        int(code),
	}
	var responseBody struct {
		AccessToken string `json:"access_token"`
	}
	var tokens schemes.Tokens
	client := resty.New()
	resp, err := client.R().
		SetResult(&responseBody).
		SetBody(requestBody).
		Post("http://localhost:8080/api/auth/code/verify/")
	if err != nil {
		return 0, nil, err
	}
	tokens.AccessToken = responseBody.AccessToken
	cookies := resp.Header()["Set-Cookie"]
	for _, cookie := range cookies {
		subString1 := strings.Split(cookie, "; ")
		subString2 := strings.Split(subString1[0], "=")
		if len(subString2) != 2 {
			continue
		}

		if subString2[0] == "atlas_rt" {
			tokens.RefreshToken = subString2[1]
			break
		}
	}
	if tokens.RefreshToken == "" || tokens.AccessToken == "" {
		return resp.StatusCode(), nil, errors.New("an error occurred while receiving tokens")
	}
	return resp.StatusCode(), &tokens, nil
}
