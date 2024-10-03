package backend

import (
	"context"
	"db"
	"fmt"
	"log"
)

type PostServer struct {
	UnimplementedPostsServer
	PrismaClient *db.PrismaClient
}

func (s *PostServer) CreatePost(ctx context.Context, in *CreatePostRequest) (*PostReply, error) {
	currentUser, err := CurrentUser(ctx)
	if err != nil {
		return nil, err
	}
	obj, err := s.PrismaClient.Post.CreateOne(
		db.Post.Content.Set(in.Text),
		db.Post.Title.Set(in.Title),

		db.Post.Author.Link(
			db.User.Email.Equals(currentUser),
		),
	).Exec(ctx)

	if err != nil {
		log.Printf("failed to create post: %v", err)
		return nil, fmt.Errorf("failed to create post")
	}

	return &PostReply{
		Id:      obj.ID,
		Title:   obj.Title,
		Content: obj.Content,
		Author:  currentUser,
	}, nil
}

func (s *PostServer) UpdatePost(ctx context.Context, in *UpdatePostRequest) (*PostReply, error) {
	obj, err := s.PrismaClient.Post.FindUnique(
		db.Post.ID.Equals(in.Id),
	).Update(
		db.Post.Content.Set(in.Text),
		db.Post.Title.Set(in.Title),
	).Exec(ctx)

	if err != nil {
		log.Printf("failed to update post: %v", err)
		return nil, fmt.Errorf("failed to update post")
	}

	return &PostReply{
		Id:      obj.ID,
		Title:   obj.Title,
		Content: obj.Content,
	}, nil
}

func (s *PostServer) ReadPost(ctx context.Context, in *ReadPostRequest) (*PostReply, error) {
	obj, err := s.PrismaClient.Post.FindUnique(
		db.Post.ID.Equals(in.Id),
	).Exec(ctx)

	if err != nil {
		log.Printf("failed to read post: %v", err)
		return nil, fmt.Errorf("failed to read post")
	}

	return &PostReply{
		Id:      obj.ID,
		Title:   obj.Title,
		Content: obj.Content,
	}, nil
}
