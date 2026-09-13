package handler

import (
	"AuthService/internal/service/mocks"
	"errors"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestValidToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockService := mocks.NewMockUserService(ctrl)
	mockService.EXPECT().RefreshToken("someToken").Return("newToken", nil)

	verifyHandler := &VerifyHandler{UserService: mockService}

	request := httptest.NewRequest(http.MethodPost, "/verify", nil)
	request.Header.Set("Authorization", "Bearer someToken")

	recorder := httptest.NewRecorder()
	verifyHandler.Verify(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Errorf("got %d, want %d", recorder.Code, http.StatusOK)
	}

	expected := "Bearer newToken"
	if recorder.Header().Get("Authorization") != expected {
		t.Errorf("expected %q, got %q", expected, recorder.Header().Get("Authorization"))
	}
}

func TestExpiredToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockService := mocks.NewMockUserService(ctrl)
	mockService.EXPECT().RefreshToken("someToken").Return("", errors.New("expired"))

	verifyHandler := &VerifyHandler{UserService: mockService}

	request := httptest.NewRequest(http.MethodPost, "/verify", nil)
	request.Header.Set("Authorization", "Bearer someToken")
	recorder := httptest.NewRecorder()

	verifyHandler.Verify(recorder, request)

	if recorder.Code != http.StatusUnauthorized {
		t.Errorf("got %d, want %d", recorder.Code, http.StatusUnauthorized)
	}
}

func TestFakeToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockService := mocks.NewMockUserService(ctrl)
	mockService.EXPECT().RefreshToken("someToken").Return("", errors.New("invalid token"))

	verifyHandler := &VerifyHandler{UserService: mockService}

	request := httptest.NewRequest(http.MethodPost, "/verify", nil)
	request.Header.Set("Authorization", "Bearer someToken")
	recorder := httptest.NewRecorder()

	verifyHandler.Verify(recorder, request)

	if recorder.Code != http.StatusUnauthorized {
		t.Errorf("got %d, want %d", recorder.Code, http.StatusUnauthorized)
	}
}
