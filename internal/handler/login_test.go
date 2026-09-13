package handler

import (
	"AuthService/internal/service/mocks"
	gomock "go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLoginSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockService := mocks.NewMockUserService(ctrl)
	mockService.EXPECT().ValidateCredentials("testguy1", "secret").Return(true, nil)
	mockService.EXPECT().GenerateToken("testguy1").Return("ourFakeToken", nil)
	loginHandler := &LoginHandler{UserService: mockService}

	request := httptest.NewRequest(http.MethodPost, "/login", nil)
	request.SetBasicAuth("testguy1", "secret")

	recorder := httptest.NewRecorder()
	loginHandler.Handle(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Errorf("Test failed. Expected %d. Got %d", http.StatusOK, recorder.Code)
	}

	expected := "Bearer ourFakeToken"
	if recorder.Header().Get("Authorization") != expected {
		t.Errorf("expected %q, got %q", expected, recorder.Header().Get("Authorization"))
	}

}

func TestLoginFail(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockService := mocks.NewMockUserService(ctrl)
	mockService.EXPECT().ValidateCredentials("testguy1", "secret").Return(false, nil)
	loginHandler := &LoginHandler{UserService: mockService}

	request := httptest.NewRequest(http.MethodPost, "/login", nil)
	request.SetBasicAuth("testguy1", "secret")

	recorder := httptest.NewRecorder()
	loginHandler.Handle(recorder, request)

	if recorder.Code != http.StatusUnauthorized {
		t.Errorf("Test failed. Expected %d. Got %d", http.StatusUnauthorized, recorder.Code)
	}
}

func TestLoginInvalidAuth(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockService := mocks.NewMockUserService(ctrl)
	loginHandler := &LoginHandler{UserService: mockService}

	request := httptest.NewRequest(http.MethodPost, "/login", nil)
	recorder := httptest.NewRecorder()
	loginHandler.Handle(recorder, request)

	if recorder.Code != http.StatusUnauthorized {
		t.Errorf("Test failed. Expected %d. Got %d", http.StatusUnauthorized, recorder.Code)
	}
}
