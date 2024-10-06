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
		Email:    "11alic22e@example.com",
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
	commentClient := NewCommentsClient(conn)
	commentReply, err := commentClient.CreateComment(ctx, &CreateCommentRequest{
		PostID: postReply.Id,
		Text:   "Hello from client",
	})
	if err != nil {
		log.Fatalf("CreateComment failed: %v", err)
	}
	fmt.Println("CreateComment response:", commentReply)
	commentReply, err = commentClient.ReadComment(ctx, &ReadCommentRequest{
		Id: commentReply.Id,
	})
	if err != nil {

		log.Fatalf("GetComment failed: %v", err)
	}
	fmt.Println("GetComment response:", commentReply)
	commentReply, err = commentClient.UpdateComment(ctx, &UpdateCommentRequest{
		Id:   commentReply.Id,
		Text: "Hello from client updated",
	})
	if err != nil {
		log.Fatalf("UpdateComment failed: %v", err)
	}
	fmt.Println("UpdateComment response:", commentReply)
	commentDelete, err := commentClient.DeleteComment(ctx, &DeleteCommentRequest{
		CommentId: commentReply.Id,
	})
	if err != nil {
		log.Fatalf("DeleteComment failed: %v", err)
	}
	fmt.Println("DeleteComment response:", commentDelete)
	postDelete, err := postClient.DeletePost(ctx, &DeletePostRequest{
		PostId: postReply.Id,
	})
	if err != nil {
		log.Fatalf("DeletePost failed: %v", err)
	}
	fmt.Println("DeletePost response", postDelete)

}
