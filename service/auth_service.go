package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/samandar2605/medium_user_service/config"
	pbn "github.com/samandar2605/medium_user_service/genproto/notification_service"
	pb "github.com/samandar2605/medium_user_service/genproto/user_service"
	grpcPkg "github.com/samandar2605/medium_user_service/pkg/grpc_client"
	"github.com/samandar2605/medium_user_service/pkg/utils"
	"github.com/samandar2605/medium_user_service/storage"
	"github.com/samandar2605/medium_user_service/storage/repo"
)

type AuthService struct {
	pb.UnimplementedAuthServiceServer
	storage    storage.StorageI
	inMemory   storage.InMemoryStorageI
	grpcClient grpcPkg.GrpcClientI
	cfg        *config.Config
}

func NewAuthService(strg storage.StorageI, inMemory storage.InMemoryStorageI, grpcConn grpcPkg.GrpcClientI, cfg config.Config) *AuthService {
	return &AuthService{
		storage:                        strg,
		inMemory:                       inMemory,
		grpcClient:                     grpcConn,
		UnimplementedAuthServiceServer: pb.UnimplementedAuthServiceServer{},
		cfg:                            &cfg,
	}
}

const (
	RegisterCodeKey   = "register_code_"
	ForgotPasswordKey = "forgot_password_code_"
)

func (s *AuthService) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.Empty, error) {
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Internal server error: %v", err)
	}

	user := repo.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Gender:    req.Gender,
		Type:      repo.UserTypeUser,
		Password:  hashedPassword,
	}

	userData, err := json.Marshal(user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Internal server error: %v", err)
	}

	err = s.inMemory.Set("user_"+user.Email, string(userData), 10*time.Minute)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Internal server error: %v", err)
	}

	go func() {
		err := s.sendVerificationCode(RegisterCodeKey, req.Email)
		if err != nil {
			fmt.Printf("failed to send verification code: %v", err)
		}
	}()

	return &pb.Empty{}, nil
}

func (s *AuthService) Verify(ctx context.Context, req *pb.VerifyRequest) (*pb.AuthResponse, error) {
	var user repo.User

	userData, err := s.inMemory.Get("user_" + req.Email)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Internal server error: %v", err)
	}

	err = json.Unmarshal([]byte(userData), &user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Internal server error: %v", err)
	}

	code, err := s.inMemory.Get(RegisterCodeKey + user.Email)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Internal server error: %v", err)
	}

	if req.Code != code {
		return nil, status.Errorf(codes.Internal, "Internal server error: %v", err)
	}

	result, err := s.storage.User().Create(&repo.User{
		ID:              user.ID,
		FirstName:       user.FirstName,
		LastName:        user.LastName,
		PhoneNumber:     user.PhoneNumber,
		Email:           user.Email,
		Gender:          user.Gender,
		Password:        user.Password,
		Username:        user.Username,
		ProfileImageUrl: user.ProfileImageUrl,
		Type:            user.Type,
		CreatedAt:       user.CreatedAt,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Internal server error: %v", err)
	}

	token, _, err := utils.CreateToken(s.cfg, &utils.TokenParams{
		UserID:   result.ID,
		Username: result.Username,
		Email:    result.Email,
		UserType: result.Type,
		Password: result.Password,
		Duration: time.Hour * 24,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Internal server error: %v", err)
	}

	return &pb.AuthResponse{
		Id:          result.ID,
		FirstName:   result.FirstName,
		LastName:    result.LastName,
		Username:    result.Username,
		Email:       result.Email,
		Type:        result.Type,
		Password:    result.Password,
		CreatedAt:   result.CreatedAt.Format(time.RFC822),
		AccessToken: token,
	}, nil
}

func (s *AuthService) Login(ctx context.Context, req *pb.VerifyRequest) (*pb.AuthResponse, error) {

	result, err := s.storage.User().GetByEmail(&req.Email)
	if err != nil {
		return nil, err
	}

	err = utils.CheckPassword(req.Code, result.Password)
	if err != nil {
		return nil, err
	}

	token, _, err := utils.CreateToken(s.cfg, &utils.TokenParams{
		UserID:   result.ID,
		Username: result.Username,
		Email:    result.Email,
		UserType: result.Type,
		Password: result.Password,
		Duration: time.Hour * 24,
	})
	if err != nil {
		return nil, err
	}

	return &pb.AuthResponse{
		Id:          result.ID,
		FirstName:   result.FirstName,
		LastName:    result.LastName,
		Password:    result.Password,
		Email:       result.Email,
		Username:    result.Username,
		Type:        result.Type,
		CreatedAt:   result.CreatedAt.Format(time.RFC822),
		AccessToken: token,
	}, nil
}

func (s *AuthService) ForgotPassword(ctx context.Context, req *pb.UserEmail) (*pb.Empty, error) {

	_, err := s.storage.User().GetByEmail(&req.Email)
	if err != nil {
		return nil, err
	}

	go func() {
		err := s.sendVerificationCode(ForgotPasswordKey, req.Email)
		if err != nil {
			fmt.Printf("failed to send verification code: %v", err)
		}
	}()

	return &pb.Empty{}, nil
}

func (s *AuthService) VerifyForgotPassword(ctx context.Context, req *pb.VerifyRequest) (*pb.AuthResponse, error) {
	code, err := s.inMemory.Get(ForgotPasswordKey + req.Email)
	if err != nil {
		return nil, errors.New("verification code has been expired")

	}

	if req.Code != code {
		return nil, errors.New("incorrect verification code")
	}

	result, err := s.storage.User().GetByEmail(&req.Email)
	if err != nil {
		return nil, err
	}

	token, _, err := utils.CreateToken(s.cfg, &utils.TokenParams{
		UserID:   result.ID,
		UserType: result.Type,
		Password: result.Password,
		Username: result.Username,
		Email:    result.Email,
		Duration: time.Hour * 24,
	})
	if err != nil {
		return nil, err
	}

	return &pb.AuthResponse{
		Id:          result.ID,
		FirstName:   result.FirstName,
		LastName:    result.LastName,
		Email:       result.Email,
		Username:    result.Username,
		Password:    result.Password,
		Type:        result.Type,
		CreatedAt:   result.CreatedAt.Format(time.RFC822),
		AccessToken: token,
	}, nil
}

func (s *AuthService) sendVerificationCode(key, email string) error {
	code, err := utils.GenerateRandomCode(6)
	if err != nil {
		return err
	}

	err = s.inMemory.Set(key+email, code, time.Minute*2)
	if err != nil {
		return err
	}

	_, err = s.grpcClient.NotificationService().SendEmail(context.Background(), &pbn.SendEmailRequest{
		To:      email,
		Subject: "Verification email",
		Body: map[string]string{
			"code": code,
		},
		Type: "verification_email",
	})

	if err != nil {
		return err
	}
	return nil
}

func (s *AuthService) VerifyToken(ctx context.Context, req *pb.VerifyTokenRequest) (*pb.AuthPayload, error) {
	accessToken := req.AccessToken

	payload, err := utils.VerifyToken(s.cfg, accessToken)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
	}
	hasPermission, err := s.storage.Permission().CheckPermission(payload.UserType, req.Resource, req.Action)
	if err != nil {
		return nil, err
	}

	return &pb.AuthPayload{
		Id:            payload.ID.String(),
		UserId:        payload.UserID,
		Email:         payload.Email,
		UserType:      payload.UserType,
		IssuedAt:      payload.IssuedAt.Format(time.RFC3339),
		ExpiredAt:     payload.ExpiredAt.Format(time.RFC3339),
		Password:      payload.Password,
		HasPermission: hasPermission,
	}, nil
}

func (s *AuthService) UpdatePassword(con context.Context, req *pb.NewPassword) (*pb.Empty, error) {
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}
	err = s.storage.User().UpdatePassword(&repo.UpdatePassword{
		UserID:   req.UserId,
		Password: hashedPassword,
	})
	if err != nil {
		return nil, err
	}

	return nil, nil
}
