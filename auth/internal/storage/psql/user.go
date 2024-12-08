package psql

import (
	"context"
	"errors"
	"net/http"

	"auth/internal/models"
	"auth/pkg/httperror"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (storage *PSQLStorage) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	query := "SELECT id, email, created_at, is_super FROM users WHERE email=$1"
	row := storage.QueryRow(ctx, query, email)
	user := models.User{}
	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.CreatedAt,
		&user.IsSuper,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, httperror.New(err, "User not found", http.StatusNotFound)
		}
		return nil, err
	}
	return &user, nil
}

func (storage *PSQLStorage) GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	query := "SELECT id, email, created_at, is_super FROM users WHERE id=$1"
	row := storage.QueryRow(ctx, query, userID)
	user := models.User{}
	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.CreatedAt,
		&user.IsSuper,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, httperror.New(err, "User not found", http.StatusNotFound)
		}
		return nil, err
	}
	return &user, nil
}

func (storage *PSQLStorage) InsertUser(ctx context.Context, user *models.User) error {
	query := "INSERT INTO users (id, email, created_at, is_super) VALUES($1,$2,$3,$4)"
	_, err := storage.Exec(
		ctx,
		query,
		user.ID,
		user.Email,
		user.CreatedAt,
		user.IsSuper,
	)
	if err != nil {
		return err
	}
	return nil
}

func (storage *PSQLStorage) UpdateUserEmail(ctx context.Context, userID uuid.UUID, email string) error {
	query := "UPDATE users SET email=$2 WHERE id=$1"
	_, err := storage.Exec(
		ctx,
		query,
		userID,
		email,
	)
	if err != nil {
		return err
	}
	return nil
}
