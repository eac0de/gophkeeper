package client

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/eac0de/gophkeeper/client/internal/schemes"
	"github.com/go-resty/resty/v2"
	"github.com/golang-jwt/jwt/v4"
)

type APIClient struct {
	gophKeeperServerAddress string
	authServerAddress       string
	Tokens                  schemes.Tokens
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
		return "", getResponseError(resp)
	}
	return result.EmailCodeId, nil
}

func (c *APIClient) VerifyEmailCode(emailCodeId string, code int) (int, error) {
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
	resp, err := c.client.R().
		SetResult(&responseBody).
		SetBody(requestBody).
		Post(fmt.Sprintf("http://%s/api/auth/code/verify/", c.authServerAddress))
	if err != nil {
		return 0, err
	}
	statusCode := resp.StatusCode()
	c.Tokens.AccessToken = responseBody.AccessToken
	cookies := resp.Header()["Set-Cookie"]
	for _, cookie := range cookies {
		subString1 := strings.Split(cookie, "; ")
		subString2 := strings.Split(subString1[0], "=")
		if len(subString2) != 2 {
			continue
		}

		if subString2[0] == "atlas_rt" {
			c.Tokens.RefreshToken = subString2[1]
			break
		}
	}
	if c.Tokens.RefreshToken == "" || c.Tokens.AccessToken == "" {
		return statusCode, errors.New("an error occurred while receiving tokens")
	}
	return statusCode, nil
}

func (c *APIClient) getAccessToken() (string, error) {
	if c.Tokens.AccessToken != "" {
		token, _, err := new(jwt.Parser).ParseUnverified(c.Tokens.AccessToken, jwt.MapClaims{})
		if err != nil {
			return "", err
		}
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			// Получаем exp (если существует)
			if exp, ok := claims["exp"].(float64); ok {
				if time.Now().Unix() < int64(exp) {
					return c.Tokens.AccessToken, nil
				} else {
					err = c.refreshTokens()
					if err != nil {
						return "", err
					}
					return c.Tokens.AccessToken, nil
				}
			} else {
				return "", errors.New("exp not found in token")
			}
		} else {
			return "", errors.New("access token isn't valid")
		}
	}
	err := c.refreshTokens()
	if err != nil {
		return "", err
	}
	return c.Tokens.AccessToken, nil
}

func (c *APIClient) refreshTokens() error {
	if c.Tokens.RefreshToken == "" {
		return errors.New("refresh token not found")
	}
	var responseBody struct {
		AccessToken string `json:"access_token"`
	}
	resp, err := c.client.R().
		SetResult(&responseBody).
		SetCookie(&http.Cookie{Name: "atlas_rt", Value: c.Tokens.RefreshToken}).
		Post(fmt.Sprintf("http://%s/api/auth/token/", c.authServerAddress))
	if err != nil {
		return err
	}
	c.Tokens.AccessToken = responseBody.AccessToken
	cookies := resp.Header()["Set-Cookie"]
	for _, cookie := range cookies {
		subString1 := strings.Split(cookie, "; ")
		subString2 := strings.Split(subString1[0], "=")
		if len(subString2) != 2 {
			continue
		}

		if subString2[0] == "atlas_rt" {
			c.Tokens.RefreshToken = subString2[1]
			break
		}
	}
	if c.Tokens.RefreshToken == "" || c.Tokens.AccessToken == "" {
		return errors.New("an error occurred while receiving tokens")
	}
	return nil
}

func getResponseError(resp *resty.Response) error {
	if resp == nil {
		return errors.New("response error: nil")
	}
	return errors.New("response error: " + resp.String() + " status: " + resp.Status())
}
