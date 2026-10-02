package handler

import (
	"errors"
	"github.com/Holocron1/authservice/internal/handler/mocks"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

type VerifyTest struct {
	caseName string
	outToken string
	outError error
	outCode  int
}

func TestVerify(t *testing.T) {
	tests := []VerifyTest{
		{caseName: "Valid Token", outToken: "newToken", outError: nil, outCode: http.StatusOK},
		{caseName: "Expired Token", outToken: "", outError: errors.New("expired"), outCode: http.StatusUnauthorized},
		{caseName: "Fake Token", outToken: "", outError: errors.New("invalid token"), outCode: http.StatusUnauthorized},
	}

	for _, test := range tests {
		t.Run(test.caseName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockService := mocks.NewMockUserService(ctrl)
			mockService.EXPECT().RefreshToken(gomock.Any(), "someToken").Return(test.outToken, test.outError)

			verifyHandler := NewVerifyHandler(mockService)
			request := httptest.NewRequest(http.MethodPost, "/verify", nil)
			request.Header.Set("Authorization", "Bearer someToken")
			recorder := httptest.NewRecorder()
			verifyHandler.Verify(recorder, request)

			if recorder.Code != test.outCode {
				t.Errorf("got %d, want %d", recorder.Code, test.outCode)
			}

			if test.outCode == http.StatusOK {
				expected := "Bearer " + test.outToken
				if recorder.Header().Get("Authorization") != expected {
					t.Errorf("expected %q, got %q", expected, recorder.Header().Get("Authorization"))
				}
			}
		})
	}
}
