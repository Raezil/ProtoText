package backend

import (
	"context"
	"db"
	"fmt"
	"log"
)

type CommentServer struct {
	UnimplementedCommentsServer
	PrismaClient *db.PrismaClient
}

func (s *CommentServer) CreateComment(ctx context.Context, in *CreateCommentRequest) (*CommentReply, error) {
	currentUser, err := CurrentUser(ctx)
	if err != nil {
		return nil, err
	}
	obj, err := s.PrismaClient.Comment.CreateOne(
		db.Comment.Content.Set(in.Text),
		db.Comment.Post.Link(
			db.Post.ID.Equals(in.PostID),
		),
		db.Comment.User.Link(
			db.User.Email.Equals(currentUser),
		),
	).Exec(ctx)

	if err != nil {
		log.Printf("failed to create comment: %v", err)
		return nil, fmt.Errorf("failed to create comment")
	}

	return &CommentReply{
		Id:      obj.ID,
		Content: obj.Content,
	}, nil
}

func (s *CommentServer) UpdateComment(ctx context.Context, in *UpdateCommentRequest) (*CommentReply, error) {
	obj, err := s.PrismaClient.Comment.FindUnique(
		db.Comment.ID.Equals(in.Id),
	).Update(
		db.Comment.Content.Set(in.Text),
	).Exec(ctx)

	if err != nil {
		log.Printf("failed to update comment: %v", err)
		return nil, fmt.Errorf("failed to update comment")
	}

	return &CommentReply{
		Id:      obj.ID,
		Content: obj.Content,
	}, nil
}

func (s *CommentServer) ReadComment(ctx context.Context, in *ReadCommentRequest) (*CommentReply, error) {
	obj, err := s.PrismaClient.Comment.FindUnique(
		db.Comment.ID.Equals(in.Id),
	).Exec(ctx)

	if err != nil {
		log.Printf("failed to read comment: %v", err)
		return nil, fmt.Errorf("failed to read comment")
	}

	return &CommentReply{
		Id:      obj.ID,
		Content: obj.Content,
	}, nil
}
