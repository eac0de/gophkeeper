package psql

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"auth/internal/models"
	"auth/pkg/httperror"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (storage *PSQLStorage) GetEmailCodeByID(ctx context.Context, emailCodeID uuid.UUID) (*models.EmailCode, error) {
	query := "SELECT id, email, code, expires_at, number_of_attempts FROM email_codes WHERE id=$1"
	row := storage.QueryRow(ctx, query, emailCodeID)
	emailCode := models.EmailCode{}
	err := row.Scan(
		&emailCode.ID,
		&emailCode.Email,
		&emailCode.Code,
		&emailCode.ExpiresAt,
		&emailCode.NumberOfAttempts,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, httperror.New(err, "EmailCode not found", http.StatusNotFound)
		}
		println(err.Error())
		return nil, err
	}
	return &emailCode, nil
}

func (storage *PSQLStorage) UpdateEmailCode(ctx context.Context, emailCode *models.EmailCode) error {
	if emailCode.ID == uuid.Nil {
		return fmt.Errorf("failed to update emailCode: ID cannot be empty")
	}
	query := "UPDATE email_codes SET email=$2, code=$3, expires_at=$4, number_of_attempts=$5 WHERE id=$1"
	_, err := storage.Exec(
		ctx,
		query,
		emailCode.ID,
		emailCode.Email,
		emailCode.Code,
		emailCode.ExpiresAt,
		emailCode.NumberOfAttempts,
	)
	if err != nil {
		return err
	}
	return nil
}

func (storage *PSQLStorage) DeleteEmailCode(ctx context.Context, emailCodeID uuid.UUID) error {
	query := "DELETE FROM email_codes WHERE id=$1"
	_, err := storage.Exec(ctx, query, emailCodeID)
	if err != nil {
		return err
	}
	return nil
}

func (storage *PSQLStorage) InsertEmailCode(ctx context.Context, emailCode *models.EmailCode) error {
	query := "INSERT INTO email_codes (id, email, code, expires_at, number_of_attempts) VALUES($1,$2,$3,$4,$5)"
	_, err := storage.Exec(
		ctx,
		query,
		emailCode.ID,
		emailCode.Email,
		emailCode.Code,
		emailCode.ExpiresAt,
		emailCode.NumberOfAttempts,
	)
	if err != nil {
		return err
	}
	return nil
}
