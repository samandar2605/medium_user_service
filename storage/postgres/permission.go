package postgres

import (
	"database/sql"
	"errors"

	_ "github.com/lib/pq"

	"github.com/jmoiron/sqlx"
	"github.com/samandar2605/medium_user_service/storage/repo"
)

type PermissionRepo struct {
	db *sqlx.DB
}

func NewPermission(db *sqlx.DB) repo.PermissionStorageI {
	return &PermissionRepo{
		db: db,
	}
}

func (ur *PermissionRepo)CheckPermission(userType,resourse,action string)(bool,error){
	query:=`
	select id from permissions
	where user_type=$1 and resource=$2 and action=$3
	`
	var id int64
	err:=ur.db.QueryRow(query,userType,resourse,action).Scan(&id)
	if err!=nil{
		if errors.Is(err,sql.ErrNoRows){
			return false,nil
		}
		return false,err
	}
	return true,nil 
}