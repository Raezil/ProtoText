package main

import (
	pb "backend"
	"context"
	"db"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func authUnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if info.FullMethod == "/authenticator.Auth/Login" || info.FullMethod == "/authenticator.Auth/Register" {
		return handler(ctx, req)
	}
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("missing metadata")
	}
	token := md["authorization"]
	if len(token) == 0 {
		return nil, fmt.Errorf("missing token")
	}

	claims, err := pb.VerifyJWT(token[0])
	ctx = metadata.AppendToOutgoingContext(ctx, "current_user", claims.Email)
	if err != nil {
		return nil, fmt.Errorf("unauthorized: %v", err)
	}
	return handler(ctx, req)
}

func main() {
	client := db.NewClient()
	if err := client.Prisma.Connect(); err != nil {
		panic(err)
	}

	defer func() {
		if err := client.Prisma.Disconnect(); err != nil {
			panic(err)
		}
	}()
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(authUnaryInterceptor),
	)
	pb.RegisterAuthServer(grpcServer, &pb.Server{
		PrismaClient: client,
	})
	pb.RegisterPostsServer(grpcServer, &pb.PostServer{
		PrismaClient: client,
	})

	// Rejestracja serwisu Auth na serwerze

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
