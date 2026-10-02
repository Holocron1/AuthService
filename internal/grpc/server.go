package grpc

import (
	"context"
	"github.com/Holocron1/authservice/pkg/authpb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type userService interface {
	ValidateCredentials(ctx context.Context, username, password string) (bool, error)
	GenerateToken(ctx context.Context, username string) (string, error)
	RefreshToken(ctx context.Context, token string) (string, error)
}

type Server struct {
	authpb.UnimplementedAuthServiceServer
	userService userService
}

func NewServer(userService userService) *Server {
	return &Server{userService: userService}
}

func (s *Server) Login(ctx context.Context, req *authpb.LoginRequest) (*authpb.LoginResponse, error) {
	ok, err := s.userService.ValidateCredentials(ctx, req.GetUsername(), req.GetPassword())
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to validate credentials")
	}
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "Invalid username or password")
	}

	jwtToken, err := s.userService.GenerateToken(ctx, req.GetUsername())
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to generate token")
	}
	return &authpb.LoginResponse{JwtToken: jwtToken}, nil
}

func (s *Server) Verify(ctx context.Context, req *authpb.VerifyRequest) (*authpb.VerifyResponse, error) {
	jwtToken := req.GetJwtToken()
	if jwtToken == "" {
		return nil, status.Error(codes.Unauthenticated, "Invalid token")
	}

	newToken, err := s.userService.RefreshToken(ctx, jwtToken)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "Failed to refresh token")
	}

	return &authpb.VerifyResponse{JwtToken: newToken}, nil
}
