package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/IainMosima/gomart/domains/auth/schema"
	"github.com/google/uuid"
)

// HTTPClient interface for dependency injection in tests
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// TestableValidateAccessToken is a version that accepts an HTTP client for testing
func (c *CognitoService) TestableValidateAccessToken(ctx context.Context, accessToken string, client HTTPClient) (*schema.UserInfoResponse, error) {
	userInfoURL := fmt.Sprintf("https://%s/oauth2/userInfo", c.config.CognitoDomain)

	req, err := http.NewRequestWithContext(ctx, "GET", userInfoURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call userInfo endpoint: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("userInfo request failed with status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var userInfoClaims struct {
		Sub           string `json:"sub"`
		Username      string `json:"username"`
		Email         string `json:"email"`
		EmailVerified string `json:"email_verified"`
		PhoneNumber   string `json:"phone_number"`
	}

	if err := json.Unmarshal(body, &userInfoClaims); err != nil {
		return nil, fmt.Errorf("failed to parse userInfo response: %w", err)
	}
	userID, _ := uuid.Parse(userInfoClaims.Sub)

	userInfo := &schema.UserInfoResponse{
		UserID:        userID,
		UserName:      userInfoClaims.Username,
		Email:         userInfoClaims.Email,
		EmailVerified: userInfoClaims.EmailVerified == "true",
		PhoneNumber:   userInfoClaims.PhoneNumber,
	}

	return userInfo, nil
}
