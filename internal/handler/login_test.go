package handler

import (
	"github.com/Holocron1/authservice/internal/handler/mocks"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

type LoginTest struct {
	caseName     string
	credMatch    bool
	outError     error
	outCode      int
	needBAuth    bool
	needGenToken bool
}

func TestLogin(t *testing.T) {
	tests := []LoginTest{
		{caseName: "Successful Login", credMatch: true, outCode: http.StatusOK, outError: nil, needBAuth: true, needGenToken: true},
		{caseName: "Login Failed", credMatch: false, outCode: http.StatusUnauthorized, outError: nil, needBAuth: false, needGenToken: false},
		{caseName: "Invalid Auth", credMatch: false, outCode: http.StatusUnauthorized, outError: nil, needBAuth: false, needGenToken: false},
	}

	for _, test := range tests {
		t.Run(test.caseName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockService := mocks.NewMockUserService(ctrl)

			loginHandler := &LoginHandler{userService: mockService}

			request := httptest.NewRequest(http.MethodPost, "/login", nil)

			if test.needBAuth {
				mockService.EXPECT().ValidateCredentials(gomock.Any(), "testguy1", "secret").Return(test.credMatch, test.outError)
				request.SetBasicAuth("testguy1", "secret")
			}

			if test.needGenToken {
				mockService.EXPECT().GenerateToken(gomock.Any(), "testguy1").Return("newToken", nil)
			}

			recorder := httptest.NewRecorder()
			loginHandler.Handle(recorder, request)

			if recorder.Code != test.outCode {
				t.Errorf("Test failed. Expected %d. Got %d", test.outCode, recorder.Code)
			}

			if test.outCode == http.StatusOK {
				expected := "Bearer newToken"
				if recorder.Header().Get("Authorization") != expected {
					t.Errorf("expected %q, got %q", expected, recorder.Header().Get("Authorization"))
				}
			}
		})
	}
}
