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

func (s *GophKeeperStorage) InsertUserAuthInfo(ctx context.Context, userAuthInfo *models.UserAuthInfo) error {
	query := `INSERT INTO user_auth_info (id, user_id, name, login, password, metadata) VALUES ($1, $2, $3, $4)`
	_, err := s.Exec(ctx, query, userAuthInfo.ID, userAuthInfo.UserID, userAuthInfo.Name, userAuthInfo.Login, userAuthInfo.Password, userAuthInfo.Metadata)
	return err
}

func (s *GophKeeperStorage) UpdateUserAuthInfo(ctx context.Context, userAuthInfo *models.UserAuthInfo) error {
	query := `UPDATE user_auth_info SET name=$3, login=$4, password=$5, metadata=$6 WHERE id=$1 AND user_id=$2`
	_, err := s.Exec(ctx, query, userAuthInfo.ID, userAuthInfo.UserID, userAuthInfo.Name, userAuthInfo.Login, userAuthInfo.Password, userAuthInfo.Metadata)
	return err
}

func (s *GophKeeperStorage) GetUserAuthInfo(ctx context.Context, userID uuid.UUID, dataID uuid.UUID) (*models.UserAuthInfo, error) {
	query := `SELECT id, user_id, name, login, password, metadata FROM user_auth_info WHERE id=$1 AND user_id=$2`
	userAuthInfo := models.UserAuthInfo{BaseUserData: models.BaseUserData{ID: dataID, UserID: userID}}
	row := s.QueryRow(ctx, query, dataID, userID)
	err := row.Scan(&userAuthInfo.ID, &userAuthInfo.UserID, &userAuthInfo.Name, &userAuthInfo.Login, &userAuthInfo.Password, &userAuthInfo.Metadata)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, httperror.New(err, "UserAuthInfo not found", http.StatusNotFound)
		}
		return nil, err
	}
	return &userAuthInfo, nil
}

func (s *GophKeeperStorage) DeleteUserAuthInfo(ctx context.Context, dataID uuid.UUID, userID uuid.UUID) error {
	query := `DELETE FROM user_auth_info WHERE id=$1 AND user_id=$1`
	_, err := s.Exec(ctx, query, dataID, userID)
	return err
}
