package main

import (
	"auth-service/internal/config"
	"auth-service/internal/database"
	"auth-service/internal/handlers"
	"auth-service/internal/middleware"
	repository "auth-service/internal/repo"
	"auth-service/internal/service"

	"log"
	"net"
	"net/http"

	pb "amelli/proto"

	"google.golang.org/grpc"
)

func main() {

	dbConfig := database.LoadConfigFromEnv()

	db, err := database.NewPostgresPool(dbConfig)
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}
	defer db.Close()

	cfg := config.Load()
	userRepo := repository.NewUserRepository(db)
	authSvc := service.NewAuthService(userRepo, string(cfg.JWTSecret))
	srv := handlers.NewServer(authSvc)
	jwtSecret := []byte("super-secret-key-change-me")
	authMiddleware := middleware.JWTAuthMiddleware(jwtSecret)

	mux := http.NewServeMux()
	mux.HandleFunc("/auth", srv.HandleAuthPage)
	mux.HandleFunc("/api/auth/login", srv.HandleLogin)
	mux.HandleFunc("/api/auth/register", srv.HandleRegister)

	mux.Handle("/api/auth/me", authMiddleware(http.HandlerFunc(srv.HandleGetMe)))
	go func() {
		if err := http.ListenAndServe(":8082", mux); err != nil {
			log.Fatalf("Ошибка HTTP сервера: %v", err)
		}
	}()

	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("Не удалось запустить gRPC сервер: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterAlhelisServiceServer(grpcServer, srv)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Ошибка работы gRPC сервера: %v", err)
	}
}
