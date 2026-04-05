package grpc

import (
	"context"
	"time"

	"github.com/alexe0110/chat-system/internal/service"
	"github.com/alexe0110/chat-system/pb"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ChatServiceServer struct {
	pb.UnimplementedChatServiceServer
	service *service.MessageService
}

func NewChatServiceServer(service *service.MessageService) *ChatServiceServer {
	return &ChatServiceServer{
		service: service,
	}
}

func (s *ChatServiceServer) SendMessage(ctx context.Context, req *pb.SendMessageRequest) (*pb.Message, error) {
	SenderID, err := uuid.Parse(req.SenderId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid sender id: %v", err)
	}

	ReceiverID, err := uuid.Parse(req.ReceiverId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid receiver id: %v", err)
	}

	result, err := s.service.SendMessage(ctx, SenderID, ReceiverID, req.MessageContent)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Problem with sending message: %v", err)
	}

	return &pb.Message{
		Id:             result.ID.String(),
		SenderId:       result.SenderID.String(),
		ReceiverId:     result.ReceiverID.String(),
		MessageContent: result.MessageContent,
		CreatedAt:      result.CreatedAt.Format(time.RFC3339),
		UpdatedAt:      result.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func (s *ChatServiceServer) GetMessageHistory(
	req *pb.MessageHistoryRequest,
	stream grpc.ServerStreamingServer[pb.Message],
) error {
	SenderID, err := uuid.Parse(req.SenderId)
	if err != nil {
		return status.Errorf(codes.InvalidArgument, "invalid sender id: %v", err)
	}

	ReceiverID, err := uuid.Parse(req.ReceiverId)
	if err != nil {
		return status.Errorf(codes.InvalidArgument, "invalid receiver id: %v", err)
	}

	messages, err := s.service.GetConversation(stream.Context(), SenderID, ReceiverID)
	if err != nil {
		return status.Errorf(codes.Internal, "problem with getting conversation: %v", err)
	}

	for _, msg := range messages {
		err := stream.Send(&pb.Message{
			Id:             msg.ID.String(),
			SenderId:       msg.SenderID.String(),
			ReceiverId:     msg.ReceiverID.String(),
			MessageContent: msg.MessageContent,
			CreatedAt:      msg.CreatedAt.Format(time.RFC3339),
			UpdatedAt:      msg.UpdatedAt.Format(time.RFC3339),
		})
		if err != nil {
			return err
		}
	}
	return nil

}
