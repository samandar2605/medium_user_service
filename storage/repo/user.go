package repo

import "time"

const (
	UserTypeSuperadmin = "superadmin"
	UserTypeUser       = "user"
)

type User struct {
	ID              int64
	FirstName       string
	LastName        string
	PhoneNumber     string
	Email           string
	Gender          string
	Password        string
	Username        string
	ProfileImageUrl string
	Type            string
	CreatedAt       time.Time
}

type CreateUser struct {
	ID              int64
	FirstName       string
	LastName        string
	PhoneNumber     string
	Email           string
	Gender          string
	Password        string
	Username        string
	ProfileImageUrl string
	Type            string
	CreatedAt       time.Time
}

type UpdateUser struct {
	Id              int64
	FirstName       string
	LastName        string
	PhoneNumber     string
	Gender          string
	Username        string
	ProfileImageUrl string
}

type GetAllUsersParams struct {
	Limit  int32
	Page   int32
	Search string
}

type GetAllUsersResult struct {
	Users []*User
	Count int32
}

type UpdatePassword struct {
	UserID   int64
	Password string
}
type DeleteUserRequest struct {
	Id int64
}

type UserStorageI interface {
	Create(user *User) (*User, error)
	Get(id int64) (*User, error)
	GetAll(params *GetAllUsersParams) (*GetAllUsersResult, error)
	Delete(*DeleteUserRequest) error
	Update(usr *UpdateUser) (*User, error)
	UpdatePassword(req *UpdatePassword) error
	GetByEmail(email *string) (*User, error)
}
