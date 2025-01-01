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

func (s *GophKeeperStorage) InsertUserBankCard(ctx context.Context, userBankCardData *models.UserBankCard) error {
	query := `INSERT INTO user_bank_card (id, user_id, name, created_at, updated_at, number, card_holder, expire_date, csc, metadata) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	_, err := s.Exec(ctx, query,
		userBankCardData.ID,
		userBankCardData.UserID,
		userBankCardData.Name,
		userBankCardData.CreatedAt,
		userBankCardData.UpdatedAt,
		userBankCardData.Number,
		userBankCardData.CardHolder,
		userBankCardData.ExpireDate,
		userBankCardData.CSC,
		userBankCardData.Metadata,
	)
	return err
}

func (s *GophKeeperStorage) UpdateUserBankCard(ctx context.Context, userBankCardData *models.UserBankCard) error {
	query := `UPDATE user_bank_card SET name=$3, updated_at=$4, number=$5, card_holder=$6, expire_date=$7, csc=$8, metadata=$8 WHERE id=$1 AND user_id=$2`
	_, err := s.Exec(ctx, query,
		userBankCardData.ID,
		userBankCardData.UserID,
		userBankCardData.Name,
		userBankCardData.UpdatedAt,
		userBankCardData.Number,
		userBankCardData.CardHolder,
		userBankCardData.ExpireDate,
		userBankCardData.CSC,
		userBankCardData.Metadata,
	)
	return err
}

func (s *GophKeeperStorage) DeleteUserBankCard(ctx context.Context, dataID uuid.UUID, userID uuid.UUID) error {
	query := `DELETE FROM user_bank_card WHERE id=$1 AND user_id=$2`
	_, err := s.Exec(ctx, query, dataID, userID)
	return err
}

func (s *GophKeeperStorage) GetUserBankCard(ctx context.Context, dataID uuid.UUID, userID uuid.UUID) (*models.UserBankCard, error) {
	query := `SELECT name, created_at, updated_at, number, card_holder, expire_date, csc, metadata FROM user_bank_card WHERE id=$1 AND user_id=$2`
	row := s.QueryRow(ctx, query, dataID, userID)
	userBankCard := models.UserBankCard{BaseUserData: models.BaseUserData{ID: dataID, UserID: userID}}
	err := row.Scan(
		&userBankCard.Name,
		&userBankCard.CreatedAt,
		&userBankCard.UpdatedAt,
		&userBankCard.Number,
		&userBankCard.CardHolder,
		&userBankCard.ExpireDate,
		&userBankCard.CSC,
		&userBankCard.Metadata,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, httperror.New(err, "UserBankCard not found", http.StatusNotFound)
		}
		return nil, err
	}
	return &userBankCard, nil
}

func (s *GophKeeperStorage) GetUserBankCardList(ctx context.Context, userID uuid.UUID, offset int) ([]models.UserBankCard, error) {
	query := `SELECT id, name, created_at, updated_at, number, card_holder, expire_date, csc, metadata FROM user_bank_card WHERE user_id=$1 LIMIT 20 OFFSET $2`

	rows, err := s.Query(ctx, query, userID, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userBankCardList []models.UserBankCard
	for rows.Next() {
		var userBankCard models.UserBankCard
		err := rows.Scan(
			&userBankCard.ID,
			&userBankCard.Name,
			&userBankCard.CreatedAt,
			&userBankCard.UpdatedAt,
			&userBankCard.Number,
			&userBankCard.CardHolder,
			&userBankCard.ExpireDate,
			&userBankCard.CSC,
			&userBankCard.Metadata,
		)
		if err != nil {
			return nil, err
		}
		userBankCard.UserID = userID
		userBankCardList = append(userBankCardList, userBankCard)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return userBankCardList, nil
}
