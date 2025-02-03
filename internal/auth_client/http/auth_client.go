package auth_client

import (
	"ChadProgress/internal/auth_client"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"
)

type AuthServiceClient struct {
	baseUrl    string
	log        *slog.Logger
	httpClient *http.Client
}

func NewAuthClient(baseUrl string, log *slog.Logger, timeOut time.Duration) *AuthServiceClient {
	return &AuthServiceClient{
		baseUrl:    baseUrl,
		log:        log,
		httpClient: &http.Client{Timeout: timeOut},
	}
}

func (c *AuthServiceClient) RegisterUser(ctx context.Context, authReq auth_client.UserAuthRequestInterface) (*auth_client.UserRegistrationResponse, error) {
	const op = "auth_client.http.auth_client.RegisterUser"
	log := c.log.With(
		slog.String("op", op),
	)
	regRequest := auth_client.UserRegistrationRequest{
		Login:    authReq.GetLogin(),
		Password: authReq.GetPassword(),
	}

	jsonPayload, err := json.Marshal(regRequest)
	if err != nil {
		log.Error(
			"error occurred: " + err.Error(),
		)
		return nil, err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.baseUrl+"/register",
		bytes.NewBuffer(jsonPayload),
	)
	if err != nil {
		log.Error("error occurred: " + err.Error())
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.Error("error occurred: " + err.Error())
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer resp.Body.Close()

	contentLen := resp.Header.Get("Content-Length")
	if resp.StatusCode != http.StatusOK && contentLen == "0" {
		log.Error("auth client response status code: " + fmt.Sprintf("%v", resp.StatusCode))
		body, _ := io.ReadAll(resp.Body)
		return nil, errors.New("failed to register user: " + string(body))
	}

	var regResp auth_client.UserRegistrationResponse
	err = json.NewDecoder(resp.Body).Decode(&regResp)
	if err != nil {
		log.Error("error occurred: " + err.Error())
		return nil, errors.New("failed to parse response from auth service")
	}
	return &regResp, nil
}
