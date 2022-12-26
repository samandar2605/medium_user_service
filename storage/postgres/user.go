package postgres

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/samandar2605/medium_user_service/pkg/utils"
	"github.com/samandar2605/medium_user_service/storage/repo"
)

type userRepo struct {
	db *sqlx.DB
}

func NewUser(db *sqlx.DB) repo.UserStorageI {
	return &userRepo{
		db: db,
	}
}

func (ur *userRepo) Create(usr *repo.User) (*repo.User, error) {
	var user repo.User
	query := `
		INSERT INTO users(
			first_name,
			last_name,
			phone_number,
			email,
			gender,
			password,
			username,
			profile_image_url,
			type
		) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING 
			id,
			first_name,
			last_name,
			COALESCE(phone_number,'') as phone_number,
			email,
			gender,
			password,
			COALESCE(username,'') as username,
			COALESCE(profile_image_url,'') as profile_image_url,
			type,
			created_at
	`

	row := ur.db.QueryRow(
		query,
		usr.FirstName,
		usr.LastName,
		utils.NullString(usr.PhoneNumber),
		usr.Email,
		usr.Gender,
		usr.Password,
		utils.NullString(usr.Username),
		utils.NullString(usr.ProfileImageUrl),
		usr.Type,
	)

	err := row.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.PhoneNumber,
		&user.Email,
		&user.Gender,
		&user.Password,
		&user.Username,
		&user.ProfileImageUrl,
		&user.Type,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (ur *userRepo) Get(id int64) (*repo.User, error) {
	var result repo.User

	query := `
		SELECT
			id,
			first_name,
			last_name,
			COALESCE(phone_number,'') as phone_number,
			email,
			gender,
			password,
			COALESCE(username,'') as username,
			COALESCE(profile_image_url, '') as profile_image_url,
			type,
			created_at
		FROM users
		WHERE id=$1
	`

	row := ur.db.QueryRow(query, id)
	err := row.Scan(
		&result.ID,
		&result.FirstName,
		&result.LastName,
		&result.PhoneNumber,
		&result.Email,
		&result.Gender,
		&result.Password,
		&result.Username,
		&result.ProfileImageUrl,
		&result.Type,
		&result.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (ur *userRepo) GetAll(params *repo.GetAllUsersParams) (*repo.GetAllUsersResult, error) {
	result := repo.GetAllUsersResult{
		Users: make([]*repo.User, 0),
	}

	offset := (params.Page - 1) * params.Limit

	limit := fmt.Sprintf(" LIMIT %d OFFSET %d ", params.Limit, offset)

	filter := ""
	if params.Search != "" {
		str := "%" + params.Search + "%"
		filter += fmt.Sprintf(`
			WHERE first_name ILIKE '%s' OR last_name ILIKE '%s' OR email ILIKE '%s' 
				OR username ILIKE '%s' OR phone_number ILIKE '%s'`,
			str, str, str, str, str,
		)
	}

	query := `
		SELECT
			id,
			first_name,
			last_name,
			COALESCE(phone_number,'') as phone_number,
			email,
			gender,
			password,
			COALESCE(username,'') as username,
			COALESCE(profile_image_url,'') as profile_image_url,
			type,
			created_at
		FROM users
		` + filter + `
		ORDER BY created_at desc
		` + limit

	rows, err := ur.db.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var u repo.User

		err := rows.Scan(
			&u.ID,
			&u.FirstName,
			&u.LastName,
			&u.PhoneNumber,
			&u.Email,
			&u.Gender,
			&u.Password,
			&u.Username,
			&u.ProfileImageUrl,
			&u.Type,
			&u.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		result.Users = append(result.Users, &u)
	}

	queryCount := `SELECT count(1) FROM users ` + filter
	err = ur.db.QueryRow(queryCount).Scan(&result.Count)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (ur *userRepo) Update(usr *repo.UpdateUser) (*repo.User, error) {
	var user repo.User
	query := `
		update users set 
			first_name=$1,
			last_name=$2,
			phone_number=$3,
			gender=$4,
			username=$5,
			profile_image_url=$6
		where id=$7 
		returning
			id,
			first_name,
			last_name,
			COALESCE(phone_number,'') as phone_number,
			email,
			gender,
			password,
			COALESCE(username,'') as username,
			COALESCE(profile_image_url, '') as profile_image_url,
			type,
			created_at
	`

	row := ur.db.QueryRow(
		query,
		usr.FirstName,
		usr.LastName,
		utils.NullString(usr.PhoneNumber),
		utils.NullString(usr.Gender),
		usr.Username,
		utils.NullString(usr.ProfileImageUrl),
		usr.Id,
	)

	if err := row.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.PhoneNumber,
		&user.Email,
		&user.Gender,
		&user.Password,
		&user.Username,
		&user.ProfileImageUrl,
		&user.Type,
		&user.CreatedAt,
	); err != nil {
		return nil, err
	}

	return &user, nil
}

func (ur *userRepo) GetByEmail(email *string) (*repo.User, error) {
	var result repo.User

	query := `
		SELECT
			id,
			first_name,
			COALESCE(last_name,'') as last_name,
			COALESCE(phone_number,'') as phone_number,
			email,
			COALESCE(gender,'male') as gender,
			password,
			COALESCE(username,'') as username,
			COALESCE(profile_image_url,'') as profile_image_url,
			type,
			created_at
		FROM users
		WHERE email=$1
	`
	row := ur.db.QueryRow(query, email)
	err := row.Scan(
		&result.ID,
		&result.FirstName,
		&result.LastName,
		&result.Password,
		&result.Email,
		&result.Gender,
		&result.Password,
		&result.Username,
		&result.ProfileImageUrl,
		&result.Type,
		&result.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (ur *userRepo) UpdatePassword(req *repo.UpdatePassword) error {
	query := `UPDATE users SET password=$1 WHERE id=$2`

	_, err := ur.db.Exec(query, req.Password, req.UserID)
	if err != nil {
		return err
	}

	return nil
}

func (ur *userRepo) Delete(req *repo.DeleteUserRequest) error {
	_, err := ur.db.Exec("delete from users where id=$1", req.Id)
	if err != nil {
		return err
	}

	return nil
}
