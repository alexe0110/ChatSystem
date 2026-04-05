package grpc

import (
	"context"
	"time"

	"github.com/alexe0110/chat-system/internal/service"
	"github.com/alexe0110/chat-system/pb"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserServiceServer struct {
	pb.UnimplementedUserServiceServer
	service *service.UserService
}

func NewUserServiceServer(service *service.UserService) *UserServiceServer {
	return &UserServiceServer{
		service: service,
	}
}

func (s *UserServiceServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	IDStr := req.Id
	ID, err := uuid.Parse(IDStr)

	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user id: %v", err)
	}

	result, err := s.service.GetByID(ctx, ID)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Problem with getting user: %v", err)
	}

	return &pb.GetUserResponse{
		Id:        result.ID.String(),
		Login:     result.Login,
		Name:      result.Name,
		CreatedAt: result.CreatedAt.Format(time.RFC3339),
		UpdatedAt: result.UpdatedAt.Format(time.RFC3339),
	}, nil

}
