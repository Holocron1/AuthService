package grpc

import (
	"context"
	"errors"
	"testing"

	"github.com/Holocron1/authservice/internal/grpc/mocks"
	"github.com/Holocron1/authservice/pkg/authpb"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestLogin(t *testing.T) {
	tests := []struct {
		name        string
		valid       bool
		validateErr error
		callGen     bool
		token       string
		tokenErr    error
		wantCode    codes.Code
	}{
		{name: "success", valid: true, callGen: true, token: "jwt", wantCode: codes.OK},
		{name: "wrong password", valid: false, wantCode: codes.Unauthenticated},
		{name: "validate fails", validateErr: errors.New("db down"), wantCode: codes.Internal},
		{name: "generate fails", valid: true, callGen: true, tokenErr: errors.New("sign"), wantCode: codes.Internal},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			svc := mocks.NewMockuserService(ctrl)
			ctx := context.Background()

			svc.EXPECT().ValidateCredentials(ctx, "user", "pass").Return(tt.valid, tt.validateErr)
			if tt.callGen {
				svc.EXPECT().GenerateToken(ctx, "user").Return(tt.token, tt.tokenErr)
			}

			resp, err := NewServer(svc).Login(ctx, &authpb.LoginRequest{Username: "user", Password: "pass"})

			if got := status.Code(err); got != tt.wantCode {
				t.Fatalf("code: got %v, want %v", got, tt.wantCode)
			}
			if tt.wantCode == codes.OK && resp.GetJwtToken() != tt.token {
				t.Errorf("token: got %q, want %q", resp.GetJwtToken(), tt.token)
			}
		})
	}
}

func TestVerify(t *testing.T) {
	tests := []struct {
		name        string
		token       string
		callRefresh bool
		newToken    string
		refreshErr  error
		wantCode    codes.Code
	}{
		{name: "success", token: "old", callRefresh: true, newToken: "new", wantCode: codes.OK},
		{name: "refresh fails", token: "bad", callRefresh: true, refreshErr: errors.New("expired"), wantCode: codes.Unauthenticated},
		{name: "empty token", token: "", wantCode: codes.Unauthenticated},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			svc := mocks.NewMockuserService(ctrl)
			ctx := context.Background()

			if tt.callRefresh {
				svc.EXPECT().RefreshToken(ctx, tt.token).Return(tt.newToken, tt.refreshErr)
			}

			resp, err := NewServer(svc).Verify(ctx, &authpb.VerifyRequest{JwtToken: tt.token})

			if got := status.Code(err); got != tt.wantCode {
				t.Fatalf("code: got %v, want %v", got, tt.wantCode)
			}
			if tt.wantCode == codes.OK && resp.GetJwtToken() != tt.newToken {
				t.Errorf("token: got %q, want %q", resp.GetJwtToken(), tt.newToken)
			}
		})
	}
}
