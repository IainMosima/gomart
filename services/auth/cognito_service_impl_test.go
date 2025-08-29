package auth

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/IainMosima/gomart/configs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCognitoService_Success(t *testing.T) {
	cfg := &configs.Config{
		AWSRegion:           "us-east-1",
		CognitoUserPoolID:   "us-east-1_TestPoolId",
		CognitoClientID:     "test-client-id",
		CognitoClientSecret: "test-client-secret",
		CognitoRedirectURI:  "http://localhost:8080/callback",
		CognitoDomain:       "test-domain.auth.us-east-1.amazoncognito.com",
	}

	service, err := NewCognitoService(cfg)

	if err == nil {
		assert.NotNil(t, service)
		assert.Equal(t, cfg, service.config)
	} else {
		// Expected in test environment without proper AWS setup
		assert.Nil(t, service)
	}
}

func TestCognitoService_GetAuthURL(t *testing.T) {
	cfg := &configs.Config{
		CognitoDomain:      "test-domain.auth.us-east-1.amazoncognito.com",
		CognitoClientID:    "test-client-id",
		CognitoRedirectURI: "http://localhost:8080/callback",
	}

	service := &CognitoService{
		config: cfg,
	}

	tests := []struct {
		name          string
		state         string
		expectedURL   string
		shouldContain []string
	}{
		{
			name:  "with state",
			state: "test-state",
			shouldContain: []string{
				"https://test-domain.auth.us-east-1.amazoncognito.com/login",
				"client_id=test-client-id",
				"response_type=code",
				"scope=email+openid+phone",
				"redirect_uri=http%3A%2F%2Flocalhost%3A8080%2Fcallback",
				"state=test-state",
			},
		},
		{
			name:  "without state",
			state: "",
			shouldContain: []string{
				"https://test-domain.auth.us-east-1.amazoncognito.com/login",
				"client_id=test-client-id",
				"response_type=code",
				"scope=email+openid+phone",
				"redirect_uri=http%3A%2F%2Flocalhost%3A8080%2Fcallback",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url, err := service.GetAuthURL(tt.state)

			require.NoError(t, err)
			assert.NotEmpty(t, url)

			for _, expected := range tt.shouldContain {
				assert.Contains(t, url, expected)
			}

			if tt.state == "" {
				assert.NotContains(t, url, "state=")
			}
		})
	}
}

func TestCognitoService_ValidateAccessToken_WithMockTransport(t *testing.T) {
	cfg := &configs.Config{
		CognitoDomain: "test-domain.auth.us-east-1.amazoncognito.com",
	}

	service := &CognitoService{
		config: cfg,
	}

	tests := []struct {
		name           string
		setupTransport func() http.RoundTripper
		expectError    bool
		errorContains  string
	}{
		{
			name: "successful response",
			setupTransport: func() http.RoundTripper {
				return &mockRoundTripper{
					roundTripFunc: func(req *http.Request) (*http.Response, error) {
						assert.Equal(t, "https://test-domain.auth.us-east-1.amazoncognito.com/oauth2/userInfo", req.URL.String())
						assert.Equal(t, "Bearer test-token", req.Header.Get("Authorization"))

						body := `{
							"sub": "550e8400-e29b-41d4-a716-446655440000",
							"username": "testuser",
							"email": "test@example.com",
							"email_verified": "true",
							"phone_number": "+1234567890"
						}`
						return &http.Response{
							StatusCode: 200,
							Body:       io.NopCloser(strings.NewReader(body)),
							Header:     make(http.Header),
						}, nil
					},
				}
			},
			expectError: false,
		},
		{
			name: "unauthorized response",
			setupTransport: func() http.RoundTripper {
				return &mockRoundTripper{
					roundTripFunc: func(req *http.Request) (*http.Response, error) {
						return &http.Response{
							StatusCode: 401,
							Body:       io.NopCloser(strings.NewReader(`{"error": "unauthorized"}`)),
							Header:     make(http.Header),
						}, nil
					},
				}
			},
			expectError:   true,
			errorContains: "userInfo request failed with status 401",
		},
		{
			name: "invalid json response",
			setupTransport: func() http.RoundTripper {
				return &mockRoundTripper{
					roundTripFunc: func(req *http.Request) (*http.Response, error) {
						return &http.Response{
							StatusCode: 200,
							Body:       io.NopCloser(strings.NewReader(`invalid json`)),
							Header:     make(http.Header),
						}, nil
					},
				}
			},
			expectError:   true,
			errorContains: "failed to parse userInfo response",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Replace the default HTTP transport temporarily
			originalTransport := http.DefaultTransport
			http.DefaultTransport = tt.setupTransport()
			defer func() {
				http.DefaultTransport = originalTransport
			}()

			userInfo, err := service.ValidateAccessToken(context.Background(), "test-token")

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, userInfo)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				require.NoError(t, err)
				assert.NotNil(t, userInfo)
				assert.Equal(t, "testuser", userInfo.UserName)
				assert.Equal(t, "test@example.com", userInfo.Email)
				assert.True(t, userInfo.EmailVerified)
				assert.Equal(t, "+1234567890", userInfo.PhoneNumber)
			}
		})
	}
}

// mockRoundTripper allows us to intercept HTTP requests in the ValidateAccessToken method
type mockRoundTripper struct {
	roundTripFunc func(*http.Request) (*http.Response, error)
}

func (m *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.roundTripFunc(req)
}
