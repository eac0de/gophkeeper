package storage

import (
	"context"

	"github.com/eac0de/gophkeeper/shared/pkg/psql"
)

type GophKeeperStorage struct {
	*psql.PSQLStorage
}

func NewGophKeeperStorage(
	ctx context.Context,
	host string,
	port string,
	username string,
	password string,
	dbName string,
) (*GophKeeperStorage, error) {
	storage, err := psql.New(ctx, host, port, username, password, dbName)
	if err != nil {
		return nil, err
	}
	return &GophKeeperStorage{storage}, nil
}
