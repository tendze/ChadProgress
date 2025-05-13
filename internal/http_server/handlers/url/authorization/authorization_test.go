package authorization

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	tests := []struct {
		name         string
		requestBody  string
		mockReturn   string
		mockError    error
		expectedCode int
		expectedResp string
	}{
		{
			name:         "Success",
			requestBody:  `{"email": "test@example.com", "password": "password123", "name": "Test User", "role": "client"}`,
			mockReturn:   "fake-jwt-token",
			mockError:    nil,
			expectedCode: http.StatusOK,
			expectedResp: `{"status":"OK","token":"fake-jwt-token"}`,
		},
		{
			name:         "Invalid role",
			requestBody:  `{"email": "test@example.com", "password": "password123", "name": "Test", "role": "invalid"}`,
			mockReturn:   "",
			mockError:    nil,
			expectedCode: http.StatusBadRequest,
			expectedResp: `{"status":"Error","error":"field Role is not valid"}`,
		},
		{
			name:         "Empty request",
			requestBody:  "",
			mockReturn:   "",
			mockError:    nil,
			expectedCode: http.StatusBadRequest,
			expectedResp: `"empty request"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAuthService := NewMockUserAuthService(ctrl)
			logger := slog.New(slog.NewTextHandler(io.Discard, nil))

			if tt.mockError != nil || tt.mockReturn != "" {
				mockAuthService.EXPECT().
					RegisterUser(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(tt.mockReturn, tt.mockError)
			}

			handler := NewUserAuthHandler(mockAuthService, logger)
			req, _ := http.NewRequest("POST", "/register", strings.NewReader(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			handler.Register(rr, req)

			assert.Equal(t, tt.expectedCode, rr.Code)
			assert.Contains(t, rr.Body.String(), tt.expectedResp)
		})
	}
}

func TestLogin(t *testing.T) {
    tests := []struct {
        name         string
        requestBody  string
        mockReturn   string
        mockError    error
        expectedCode int
        expectedResp string
    }{
        {
            name:         "Success",
            requestBody:  `{"email": "test@example.com", "password": "valid"}`,
            mockReturn:   "fake-jwt-token",
            mockError:    nil,
            expectedCode: http.StatusOK,
            expectedResp: `{"status":"OK","token":"fake-jwt-token"}`,
        },
        {
            name:         "Empty email",
            requestBody:  `{"password": "password123"}`,
            mockReturn:   "",
            mockError:    nil,
            expectedCode: http.StatusBadRequest,
            expectedResp: `{"status":"Error","error":"field Email is a required field"}`,
        },
		{
			name:         "Empty request",
			requestBody:  "",
			mockReturn:   "",
			mockError:    nil,
			expectedCode: http.StatusBadRequest,
			expectedResp: `"empty request"`,
		},
		{
			name:         "Invalid json",
			requestBody:  `{invalid}`,
			mockReturn:   "",
			mockError:    nil,
			expectedCode: http.StatusBadRequest,
			expectedResp: "failed to decode request bod",
		},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            ctrl := gomock.NewController(t)
            defer ctrl.Finish()

            mockAuthService := NewMockUserAuthService(ctrl)
            logger := slog.New(slog.NewTextHandler(io.Discard, nil))

            if tt.mockError != nil || tt.mockReturn != "" {
                mockAuthService.EXPECT().
                    Login(gomock.Any(), gomock.Any()).
                    Return(tt.mockReturn, tt.mockError)
            }

            handler := NewUserAuthHandler(mockAuthService, logger)
            req, _ := http.NewRequest("POST", "/login", strings.NewReader(tt.requestBody))
            req.Header.Set("Content-Type", "application/json")

            rr := httptest.NewRecorder()
            handler.Login(rr, req)

            assert.Equal(t, tt.expectedCode, rr.Code)
            assert.Contains(t, rr.Body.String(), tt.expectedResp)
        })
    }
}
