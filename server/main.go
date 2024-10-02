package main

import (
	pb "authenticator"
	"context"
	"db"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var jwtSecret = []byte("your-256-bit-secret")

type Claims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func generateJWT(email string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	// Tworzymy token z podpisem HMAC
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func verifyJWT(tokenStr string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid token: %v", err)
	}
	return claims, nil
}

type Server struct {
	pb.UnimplementedAuthServer
	prismaClient *db.PrismaClient
}

func CurrentUser(ctx context.Context) (string, error) {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		return "", fmt.Errorf("missing metadata")
	}
	current_user := md["current_user"]
	return current_user[0], nil
}

func (s *Server) SampleProtected(ctx context.Context, in *pb.ProtectedRequest) (*pb.ProtectedReply, error) {
	currentUser, err := CurrentUser(ctx)
	if err != nil {
		return nil, err
	}
	return &pb.ProtectedReply{
		Result: in.Text + " " + currentUser,
	}, nil
}

func (s *Server) Login(ctx context.Context, in *pb.LoginRequest) (*pb.LoginReply, error) {
	log.Println("Login attempt for email:", in.Email)

	user, err := s.prismaClient.User.FindUnique(
		db.User.Email.Equals(in.Email),
	).Exec(ctx)

	if err != nil {
		log.Printf("User not found: %v", err)
		return nil, fmt.Errorf("incorrect email or password")
	}

	if user.Password != in.Password {
		log.Println("Invalid password")
		return nil, fmt.Errorf("incorrect email or password")
	}

	token, err := generateJWT(in.Email)
	fmt.Println(token)
	if err != nil {
		log.Printf("Error generating token: %v", err)
		return nil, fmt.Errorf("could not generate token: %v", err)
	}

	log.Printf("Generated token: %s", token)

	// Zwracanie tokena w odpowiedzi
	return &pb.LoginReply{
		Token: token,
	}, nil
}

func (s *Server) Register(ctx context.Context, in *pb.RegisterRequest) (*pb.RegisterReply, error) {
	obj, err := s.prismaClient.User.CreateOne(
		db.User.Name.Set(in.Name),
		db.User.Password.Set(in.Password),
		db.User.Email.Set(in.Email),
	).Exec(ctx)

	if err != nil {
		log.Printf("failed to create user: %v", err)
		return nil, fmt.Errorf("failed to register user")
	}

	return &pb.RegisterReply{
		Reply: fmt.Sprintf("Congratulations, User id: %s got created!", obj.ID),
	}, nil
}

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

	claims, err := verifyJWT(token[0])
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
	pb.RegisterAuthServer(grpcServer, &Server{
		prismaClient: client,
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
