package main

import (
	"auth-service/internal/database"
	"auth-service/internal/handlers"

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

	srv := &handlers.Server{DB: db}

	mux := http.NewServeMux()

	mux.HandleFunc("/auth", srv.HandleAuthPage)
	mux.HandleFunc("/api/auth/login", srv.HandleLogin)       // Исправлено
	mux.HandleFunc("/api/auth/register", srv.HandleRegister) // Исправлено

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
