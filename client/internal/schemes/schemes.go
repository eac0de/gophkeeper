package schemes

import "github.com/google/uuid"

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type DataRow struct {
	ID    uuid.UUID `json:"id"`
	Title string    `json:"title"`
	Type  string    `json:"type"`
}
