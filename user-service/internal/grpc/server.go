package grpcserver

import (
	"context"
	"log"
	"strconv"

	pb "amelli/proto"

	"user-service/internal/service"
)

type Server struct {
	pb.UnimplementedAlhelisServiceServer
	userService *service.UserService
}

func NewServer(userService *service.UserService) *Server {
	return &Server{
		userService: userService,
	}
}

func (s *Server) CreateProfile(
	ctx context.Context,
	req *pb.CreateProfileRequest,
) (*pb.CreateProfileResponse, error) {

	log.Printf("gRPC CreateProfile получил запрос: userID=%s", req.GetUserId())

	userID, err := strconv.ParseInt(req.GetUserId(), 10, 64)
	if err != nil {
		log.Printf("gRPC CreateProfile ошибка ParseInt: %v", err)
		return nil, err
	}

	if err := s.userService.CreateProfile(ctx, userID); err != nil {
		log.Printf("gRPC CreateProfile ошибка: %v", err)
		return nil, err
	}

	log.Printf("gRPC CreateProfile успешно: userID=%d", userID)

	return &pb.CreateProfileResponse{
		Message: "Профиль создан",
	}, nil
}
