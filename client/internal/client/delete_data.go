package client

import (
	"fmt"
	"net/http"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/eac0de/gophkeeper/client/internal/schemes"
	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
)

func (c *APIClient) DeleteData(itemID uuid.UUID, itemIndex, activeTab int) tea.Msg {
	req := c.client.R()

	var url string
	switch activeTab {
	case 0:
		url = fmt.Sprintf("http://%s/api/gophkeeper/text_data/%s/", c.gophKeeperServerAddress, itemID.String())
	case 1:
		url = fmt.Sprintf("http://%s/api/gophkeeper/file_data/%s/", c.gophKeeperServerAddress, itemID.String())
	case 2:
		url = fmt.Sprintf("http://%s/api/gophkeeper/bank_cards/%s/", c.gophKeeperServerAddress, itemID.String())
	case 3:
		url = fmt.Sprintf("http://%s/api/gophkeeper/auth_info/%s/", c.gophKeeperServerAddress, itemID.String())
	}

	var resp *resty.Response
	var statusCode int
	var err error
	var token string

	for i := 0; i < 2; i++ {
		token, err = c.getAccessToken()
		if err != nil {
			return schemes.DeleteListItemMsg{
				Err:        err,
				StatusCode: http.StatusUnauthorized,
			}
		}
		resp, err = req.SetAuthToken(token).Delete(url)
		if err != nil {
			return schemes.DeleteListItemMsg{Err: err}
		}
		statusCode = resp.StatusCode()
		if statusCode == http.StatusUnauthorized {
			c.Tokens.AccessToken = ""
			continue
		}
		if statusCode == http.StatusNoContent {
			return schemes.DeleteListItemMsg{
				ItemIndex:  itemIndex,
				ActiveTab:  activeTab,
				StatusCode: statusCode,
			}
		}
		break
	}
	return schemes.DeleteListItemMsg{Err: getResponseError(resp), StatusCode: statusCode}
}
