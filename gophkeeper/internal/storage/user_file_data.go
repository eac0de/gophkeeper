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

func (s *GophKeeperStorage) InsertUserFileData(ctx context.Context, userFileData *models.UserFileData) error {
	query := `INSERT INTO user_file_data (id, user_id, name, created_at, updated_at, path_to_file, metadata) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := s.Exec(
		ctx,
		query,
		userFileData.ID,
		userFileData.UserID,
		userFileData.Name,
		userFileData.CreatedAt,
		userFileData.UpdatedAt,
		userFileData.PathToFile,
		userFileData.Metadata,
	)
	return err
}

func (s *GophKeeperStorage) UpdateUserFileData(ctx context.Context, userFileData *models.UserFileData) error {
	query := `UPDATE user_file_data SET name=$3, updated_at=$4, path_to_file=$5, metadata=$6 WHERE id=$1 AND user_id=$2`
	_, err := s.Exec(
		ctx,
		query, userFileData.ID,
		userFileData.UserID,
		userFileData.Name,
		userFileData.UpdatedAt,
		userFileData.PathToFile,
		userFileData.Metadata,
	)
	return err
}

func (s *GophKeeperStorage) GetUserFileData(ctx context.Context, dataID uuid.UUID, userID uuid.UUID) (*models.UserFileData, error) {
	query := `SELECT name, created_at, updated_at, path_to_file, metadata FROM user_file_data WHERE id=$1 AND user_id=$2`
	row := s.QueryRow(ctx, query, dataID, userID)
	userFileData := models.UserFileData{BaseUserData: models.BaseUserData{ID: dataID, UserID: userID}}
	err := row.Scan(
		&userFileData.Name,
		&userFileData.CreatedAt,
		&userFileData.UpdatedAt,
		&userFileData.PathToFile,
		&userFileData.Metadata,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, httperror.New(err, "UserFileData not found", http.StatusNotFound)
		}
		return nil, err
	}
	return &userFileData, nil
}

func (s *GophKeeperStorage) DeleteUserFileData(ctx context.Context, dataID uuid.UUID, userID uuid.UUID) error {
	query := `DELETE FROM user_file_data WHERE id=$1 AND user_id=$2`
	_, err := s.Exec(ctx, query, dataID, userID)
	return err
}

func (s *GophKeeperStorage) GetUserFileDataList(ctx context.Context, userID uuid.UUID, offset int) ([]models.UserFileData, error) {
	query := `SELECT id, name, created_at, updated_at, path_to_file, metadata FROM user_file_data WHERE user_id=$1 LIMIT 20 OFFSET $2`

	rows, err := s.Query(ctx, query, userID, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userFileDataList []models.UserFileData
	for rows.Next() {
		var userFileData models.UserFileData
		err := rows.Scan(
			&userFileData.ID,
			&userFileData.Name,
			&userFileData.CreatedAt,
			&userFileData.UpdatedAt,
			&userFileData.PathToFile,
			&userFileData.Metadata,
		)
		if err != nil {
			return nil, err
		}
		userFileData.UserID = userID
		userFileDataList = append(userFileDataList, userFileData)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return userFileDataList, nil
}
