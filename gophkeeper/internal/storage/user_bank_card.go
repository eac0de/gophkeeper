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
	query := `INSERT INTO user_bank_card (id, user_id, name, number, card_holder, expire_date, csc, metadata) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := s.Exec(ctx, query,
		userBankCardData.ID,
		userBankCardData.UserID,
		userBankCardData.Name,
		userBankCardData.Number,
		userBankCardData.CardHolder,
		userBankCardData.ExpireDate,
		userBankCardData.CSC,
		userBankCardData.Metadata,
	)
	return err
}

func (s *GophKeeperStorage) UpdateUserBankCard(ctx context.Context, userBankCardData *models.UserBankCard) error {
	query := `UPDATE user_bank_card SET name=$3, number=$4, card_holder=$5, expire_date=$6, csc=$7, metadata=$8 WHERE id=$1 AND user_id=$2`
	_, err := s.Exec(ctx, query,
		userBankCardData.ID,
		userBankCardData.UserID,
		userBankCardData.Name,
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
	query := `SELECT name, number, card_holder, expire_date, csc, metadata FROM user_bank_card WHERE id=$1 AND user_id=$2`
	row := s.QueryRow(ctx, query, dataID, userID)
	userBankCard := models.UserBankCard{BaseUserData: models.BaseUserData{ID: dataID, UserID: userID}}
	err := row.Scan(
		&userBankCard.Name,
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
