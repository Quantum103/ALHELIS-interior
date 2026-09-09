package main

import (
	"log"
	"net"
	"net/http"
	"os"

	pb "amelli/proto"

	"user-service/internal/database"
	grpcuser "user-service/internal/grpc"
	"user-service/internal/handlers"
	"user-service/internal/middleware"
	"user-service/internal/repo"
	"user-service/internal/service"

	"github.com/gorilla/mux"
	grpcgo "google.golang.org/grpc"
)

func main() {
	cfg := database.LoadConfigFromEnv()

	db, err := database.NewPostgresPool(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	userRepo := repo.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET не задан")
	}

	r := mux.NewRouter()

	r.Handle(
		"/api/user/profile",
		middleware.JWTAuthMiddleware([]byte(jwtSecret))(
			http.HandlerFunc(userHandler.ShowLK),
		),
	).Methods(http.MethodGet)

	r.Handle(
		"/api/user/profile",
		middleware.JWTAuthMiddleware([]byte(jwtSecret))(
			http.HandlerFunc(userHandler.CreateProfile),
		),
	).Methods(http.MethodPost)

	grpcHandler := grpcuser.NewServer(userService)

	lis, err := net.Listen("tcp", ":50053")
	if err != nil {
		log.Fatalf("Не удалось запустить gRPC listener: %v", err)
	}

	grpcServer := grpcgo.NewServer()

	pb.RegisterAlhelisServiceServer(grpcServer, grpcHandler)

	go func() {
		log.Println("User Service gRPC запущен на :50053")

		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Ошибка работы gRPC сервера: %v", err)
		}
	}()

	log.Println("User Service HTTP запущен на :8082")

	if err := http.ListenAndServe(":8082", r); err != nil {
		log.Fatal(err)
	}
}
