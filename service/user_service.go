package service

import (
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/samandar2605/medium_user_service/genproto/user_service"
	"github.com/samandar2605/medium_user_service/pkg/utils"
	"github.com/samandar2605/medium_user_service/storage"
	"github.com/samandar2605/medium_user_service/storage/repo"
)

func parseUserModel(user *repo.User) *pb.User {
	return &pb.User{
		Id:              user.ID,
		FirstName:       user.FirstName,
		LastName:        user.LastName,
		PhoneNumber:     user.PhoneNumber,
		Email:           user.Email,
		Gender:          user.Gender,
		Password:        user.Password,
		Username:        user.Username,
		ProfileImageUrl: user.ProfileImageUrl,
		Type:            user.Type,
		CreatedAt:       user.CreatedAt.Format(time.RFC3339),
	}
}

type UserService struct {
	pb.UnimplementedUserServiceServer
	storage  storage.StorageI
	inMemory storage.InMemoryStorageI
}

func NewUserService(strg storage.StorageI, inMemory storage.InMemoryStorageI) *UserService {
	return &UserService{
		storage:                        strg,
		inMemory:                       inMemory,
		UnimplementedUserServiceServer: pb.UnimplementedUserServiceServer{},
	}
}

func (s *UserService) Create(ctx context.Context, req *pb.User) (*pb.User, error) {
	user, err := s.storage.User().Create(&repo.User{
		FirstName:       req.FirstName,
		LastName:        req.LastName,
		PhoneNumber:     req.PhoneNumber,
		Email:           req.Email,
		Gender:          req.Gender,
		Password:        req.Password,
		Username:        req.Username,
		ProfileImageUrl: req.ProfileImageUrl,
		Type:            req.Type,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Internal server error: %v", err)
	}

	return parseUserModel(user), nil
}
func (s *UserService) Get(ctx context.Context, req *pb.IdRequest) (*pb.User, error) {
	resp, err := s.storage.User().Get(req.Id)
	if err != nil {
		return nil, err
	}
	return parseUserModel(resp), nil
}

func (s *UserService) GetAll(ctx context.Context, req *pb.GetAllUsersRequest) (*pb.GetAllUsersResponse, error) {
	result, err := s.storage.User().GetAll(&repo.GetAllUsersParams{
		Page:   req.Page,
		Limit:  req.Limit,
		Search: req.Search,
	})
	if err != nil {
		return nil, err
	}
	response := pb.GetAllUsersResponse{}
	response.Count = result.Count
	for _, i := range result.Users {
		response.Users = append(response.Users, parseUserModel(i))
	}
	return &response, nil
}

func (s *UserService) Update(ctx context.Context, req *pb.User) (*pb.User, error) {
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}
	user, err := s.storage.User().Update(&repo.User{
		ID:              req.Id,
		FirstName:       req.FirstName,
		LastName:        req.LastName,
		PhoneNumber:     req.PhoneNumber,
		Email:           req.Email,
		Username:        req.Username,
		Password:        hashedPassword,
		ProfileImageUrl: req.ProfileImageUrl,
		Type:            req.Type,
	})
	if err != nil {
		return nil, err
	}

	return parseUserModel(user), nil
}

func (s *UserService) Delete(ctx context.Context, req *pb.IdRequest) (*pb.Empty, error) {
	err := s.storage.User().Delete(int(req.Id))
	if err != nil {
		return nil, err
	}
	return nil, nil
}
