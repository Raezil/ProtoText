package main

import (
	. "backend"
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := NewAuthClient(conn)
	_, err = client.Register(context.Background(), &RegisterRequest{
		Email:    "alic22242e@example.com",
		Password: "password",
	})
	if err != nil {
		log.Fatalf("Register should have failed")
	}
	loginReply, err := client.Login(context.Background(), &LoginRequest{
		Email:    "alic22e@example.com",
		Password: "password",
	})
	if err != nil {
		log.Fatalf("Login failed: %v", err)
	}

	token := loginReply.Token
	fmt.Println("Received JWT token:", token)
	md := metadata.Pairs("authorization", token)
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	protectedReply, err := client.SampleProtected(ctx, &ProtectedRequest{
		Text: "Hello from client",
	})
	if err != nil {
		log.Fatalf("SampleProtected failed: %v", err)
	}
	fmt.Println("SampleProtected response:", protectedReply.Result)
	postClient := NewPostsClient(conn)
	postReply, err := postClient.CreatePost(ctx, &CreatePostRequest{
		Title: "Hello from client",
		Text:  "Hello from client",
	})
	if err != nil {
		log.Fatalf("CreatePost failed: %v", err)
	}
	fmt.Println("CreatePost response:", postReply)
	postReply, err = postClient.ReadPost(ctx, &ReadPostRequest{
		Id: postReply.Id,
	})
	if err != nil {
		log.Fatalf("GetPost failed: %v", err)
	}
	fmt.Println("GetPost response:", postReply)

	postReply, err = postClient.UpdatePost(ctx, &UpdatePostRequest{
		Id:    postReply.Id,
		Title: "Hello from client updated",
		Text:  "Hello from client updated",
	})
	if err != nil {
		log.Fatalf("UpdatePost failed: %v", err)
	}
	fmt.Println("UpdatePost response:", postReply)

}
