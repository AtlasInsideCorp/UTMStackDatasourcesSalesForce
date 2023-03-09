package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/c3s4rfred/sforceds/configs"
	"net/http"
	"net/url"
	"strings"
)

type LoginResponse struct {
	AccessToken string `json:"access_token"`
	InstanceURL string `json:"instance_url"`
	ID          string `json:"id"`
	TokenType   string `json:"token_type"`
	IssuedAt    string `json:"issued_at"`
	Signature   string `json:"signature"`
}

func Login() (*LoginResponse, error) {
	body := url.Values{}
	body.Set("grant_type", configs.GrantType)
	body.Set("client_id", configs.ClientId)
	body.Set("client_secret", configs.ClientSecret)
	body.Set("username", configs.Username)
	body.Set("password", configs.Password+configs.SecurityToken)

	ctx, cancelFn := context.WithTimeout(context.Background(), configs.OAuthDialTimeout)
	defer cancelFn()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, configs.OAuthService+configs.LoginEndpoint, strings.NewReader(body.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	httpResp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non-200 status code returned on OAuth authentication call: %v", httpResp.StatusCode)
	}

	var loginResponse LoginResponse
	err = json.NewDecoder(httpResp.Body).Decode(&loginResponse)
	if err != nil {
		return nil, err
	}

	return &loginResponse, nil
}
