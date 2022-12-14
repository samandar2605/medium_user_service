package main

import (
	"fmt"
	"log"
	"net"

	"github.com/go-redis/redis/v9"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	pb "github.com/samandar2605/medium_user_service/genproto/user_service"
	grpcPkg "github.com/samandar2605/medium_user_service/pkg/grpc_client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/samandar2605/medium_user_service/config"
	"github.com/samandar2605/medium_user_service/service"
	"github.com/samandar2605/medium_user_service/storage"
)

func main() {
	cfg := config.Load(".")

	psqlUrl := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.Database,
	)

	psqlConn, err := sqlx.Connect("postgres", psqlUrl)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: cfg.Redis.Addr,
	})
	strg := storage.NewStoragePg(psqlConn)
	inMemory := storage.NewInMemoryStorage(rdb)

	grpcConn, err := grpcPkg.New(cfg)
	if err != nil {
		log.Fatalf("failed to get grpc connections: %v", err)
	}

	userService := service.NewUserService(strg, inMemory)
	authService := service.NewAuthService(strg, inMemory, grpcConn, cfg)
	lis, err := net.Listen("tcp", cfg.GrpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	reflection.Register(s)

	pb.RegisterUserServiceServer(s, userService)
	pb.RegisterAuthServiceServer(s, authService)

	log.Println("Grpc server started in port ", cfg.GrpcPort)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Error while listening: %v", err)
	}

}
