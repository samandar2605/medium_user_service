package storage

import (
	"github.com/jmoiron/sqlx"

	"github.com/samandar2605/medium_user_service/storage/repo"
	"github.com/samandar2605/medium_user_service/storage/postgres"
)

type StorageI interface {
	User() repo.UserStorageI
}

type storagePg struct {
	userRepo repo.UserStorageI
}

func NewStoragePg(db *sqlx.DB) StorageI {
	return &storagePg{
		userRepo: postgres.NewUser(db),
	}
}

func (s *storagePg) User() repo.UserStorageI {
	return s.userRepo
}
