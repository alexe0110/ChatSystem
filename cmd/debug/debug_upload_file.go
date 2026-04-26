package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/alexe0110/chat-system/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func main() {
	fileName := "kek.png"
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := pb.NewChatServiceClient(conn)

	ctx := metadata.AppendToOutgoingContext(context.Background(), "authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.")

	// Открываем стрим
	stream, err := client.UploadFile(ctx)
	if err != nil {
		log.Fatal(err)
	}

	data, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}

	// Отправляем чанками по 1024 байт
	chunkSize := 1024
	for i := 0; i < len(data); i += chunkSize {
		end := i + chunkSize
		if end > len(data) {
			end = len(data)
		}

		err := stream.Send(&pb.FileChunk{
			FileName: fileName,
			Data:     data[i:end],
		})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Sent chunk %d-%d\n", i, end)
	}

	resp, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Uploaded! ID: %s, URL: %s, Size: %d\n", resp.FileId, resp.Url, resp.Size)
}
