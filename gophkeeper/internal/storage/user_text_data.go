package storage

import (
	"context"
	"errors"
	"net/http"

	"github.com/eac0de/gophkeeper/internal/models"
	"github.com/eac0de/gophkeeper/shared/pkg/httperror"
	"github.com/google/uuid"
	"github.com/jackc/pgx"
)

func (s *GophKeeperStorage) InsertUserTextData(ctx context.Context, userTextData *models.UserTextData) error {
	query := `INSERT INTO user_text_data (id, user_id, name, data, metadata) VALUES ($1, $2, $3, $4, $5)`
	_, err := s.Exec(ctx, query, userTextData.ID, userTextData.UserID, userTextData.Name, userTextData.Data, userTextData.Metadata)
	return err
}

func (s *GophKeeperStorage) UpdateUserTextData(ctx context.Context, userTextData *models.UserTextData) error {
	query := `UPDATE user_text_data SET name=$3, data=$4, metadata=$5 WHERE id=$1 AND user_id=$2`
	_, err := s.Exec(ctx, query, userTextData.ID, userTextData.UserID, userTextData.Name, userTextData.Data, userTextData.Metadata)
	return err
}

func (s *GophKeeperStorage) GetUserTextData(ctx context.Context, dataID uuid.UUID, userID uuid.UUID) (*models.UserTextData, error) {
	query := `SELECT name, data, metadata FROM user_text_data WHERE id=$1 AND user_id=$2`
	row := s.QueryRow(ctx, query, dataID, userID)
	userTextData := models.UserTextData{BaseUserData: models.BaseUserData{ID: dataID, UserID: userID}}
	err := row.Scan(
		&userTextData.Name,
		&userTextData.Data,
		&userTextData.Metadata,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, httperror.New(err, "UserTextData not found", http.StatusNotFound)
		}
		return nil, err
	}
	return &userTextData, nil
}

func (s *GophKeeperStorage) DeleteUserTextData(ctx context.Context, dataID uuid.UUID, userID uuid.UUID) error {
	query := `DELETE FROM user_text_data WHERE id=$1 AND user_id=$2`
	_, err := s.Exec(ctx, query, dataID, userID)
	return err
}
