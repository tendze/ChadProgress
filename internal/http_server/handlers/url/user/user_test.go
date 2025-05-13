package userhandler

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"ChadProgress/internal/models"
	service "ChadProgress/internal/services"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCreateTrainer(t *testing.T) {
	tests := []struct {
		name         string
		email        string
		requestBody  string
		mockError    error
		expectedCode int
		expectedResp string
	}{
		{
			name:         "Success",
			email:        "trainer@example.com",
			requestBody:  `{"qualification":"Certified","experience":"5 years","achievement":"Champion"}`,
			mockError:    nil,
			expectedCode: http.StatusOK,
			expectedResp: `{"status":"OK"}`,
		},
		{
			name:         "Invalid role (client exists)",
			email:        "client@example.com",
			requestBody:  `{"qualification":"Certified","experience":"5 years","achievement":"Champion"}`,
			mockError:    service.ErrInvalidRoleRequest,
			expectedCode: http.StatusBadRequest,
			expectedResp: `"cannot create trainer profile while being client"`,
		},
		{
			name:         "Validation error",
			email:        "user@example.com",
			requestBody:  `{"experience":"5 years"}`,
			mockError:    nil,
			expectedCode: http.StatusBadRequest,
			expectedResp: `{"status":"Error","error":"field Qualification is a required field, field Achievement is a required field"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockService := NewMockUserService(ctrl)
			logger := slog.New(slog.NewTextHandler(io.Discard, nil))

			ctx := context.WithValue(context.Background(), models.ContextUserKey, tt.email)
			req, _ := http.NewRequestWithContext(ctx, "POST", "/trainer", strings.NewReader(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")

			if tt.mockError != nil || tt.expectedCode == http.StatusOK {
				mockService.EXPECT().
					CreateTrainer(tt.email, gomock.Any(), gomock.Any(), gomock.Any()).
					Return(tt.mockError)
			}

			handler := NewUserHandler(logger, mockService)
			rr := httptest.NewRecorder()
			handler.CreateTrainer(rr, req)

			assert.Equal(t, tt.expectedCode, rr.Code)
			if tt.expectedResp != "" {
				assert.Contains(t, rr.Body.String(), tt.expectedResp)
			}
		})
	}
}
