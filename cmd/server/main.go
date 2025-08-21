package main

import (
	"log"
	"net"
	"os"

	connection "users-grpc/internal/database"
	"users-grpc/internal/seed"
	"users-grpc/internal/service"
	"users-grpc/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	conn := connection.New()

	if os.Getenv("SEED") == "true" {
		seed.Run(conn)
	}

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	proto.RegisterUserServiceServer(grpcServer, service.NewUserServiceServer(conn))
	reflection.Register(grpcServer)

	log.Println("gRPC server running on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
