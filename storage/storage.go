package storage

import (
	"github.com/jmoiron/sqlx"

	"github.com/samandar2605/medium_user_service/storage/postgres"
	"github.com/samandar2605/medium_user_service/storage/repo"
)

type StorageI interface {
	User() repo.UserStorageI
	Permission() repo.PermissionStorageI
}

type storagePg struct {
	userRepo       repo.UserStorageI
	permissionRepo repo.PermissionStorageI
}

func NewStoragePg(db *sqlx.DB) StorageI {
	return &storagePg{
		userRepo:       postgres.NewUser(db),
		permissionRepo: postgres.NewPermission(db),
	}
}

func (s *storagePg) User() repo.UserStorageI {
	return s.userRepo
}

func (s *storagePg) Permission() repo.PermissionStorageI {
	return s.permissionRepo
}
