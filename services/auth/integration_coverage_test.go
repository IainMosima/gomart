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

// CustomRoundTripper allows us to intercept HTTP requests in the original ValidateAccessToken method
type CustomRoundTripper struct {
	RoundTripFunc func(*http.Request) (*http.Response, error)
}

func (c *CustomRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return c.RoundTripFunc(req)
}

func TestCognitoService_ValidateAccessToken_WithCustomTransport(t *testing.T) {
	// This test uses a custom HTTP transport to intercept the real ValidateAccessToken method
	cfg := &configs.Config{
		CognitoDomain: "test-domain.auth.us-east-1.amazoncognito.com",
	}

	service := &CognitoService{
		config: cfg,
	}

	tests := []struct {
		name           string
		setupTransport func() *CustomRoundTripper
		expectError    bool
		errorContains  string
	}{
		{
			name: "successful response",
			setupTransport: func() *CustomRoundTripper {
				return &CustomRoundTripper{
					RoundTripFunc: func(req *http.Request) (*http.Response, error) {
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
			setupTransport: func() *CustomRoundTripper {
				return &CustomRoundTripper{
					RoundTripFunc: func(req *http.Request) (*http.Response, error) {
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
			setupTransport: func() *CustomRoundTripper {
				return &CustomRoundTripper{
					RoundTripFunc: func(req *http.Request) (*http.Response, error) {
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
		{
			name: "body read error",
			setupTransport: func() *CustomRoundTripper {
				return &CustomRoundTripper{
					RoundTripFunc: func(req *http.Request) (*http.Response, error) {
						// Create a response with a body that will fail to read
						return &http.Response{
							StatusCode: 200,
							Body:       &errorReadCloser{},
							Header:     make(http.Header),
						}, nil
					},
				}
			},
			expectError:   true,
			errorContains: "failed to read response body",
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

// errorReadCloser implements io.ReadCloser but always returns an error on Read
type errorReadCloser struct{}

func (e *errorReadCloser) Read(p []byte) (int, error) {
	return 0, io.ErrUnexpectedEOF
}

func (e *errorReadCloser) Close() error {
	return nil
}

func TestCognitoService_ValidateAccessToken_EmptyResponseBody(t *testing.T) {
	cfg := &configs.Config{
		CognitoDomain: "test-domain.auth.us-east-1.amazoncognito.com",
	}

	service := &CognitoService{
		config: cfg,
	}

	// Test with empty response body
	originalTransport := http.DefaultTransport
	http.DefaultTransport = &CustomRoundTripper{
		RoundTripFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader("")),
				Header:     make(http.Header),
			}, nil
		},
	}
	defer func() {
		http.DefaultTransport = originalTransport
	}()

	userInfo, err := service.ValidateAccessToken(context.Background(), "test-token")

	assert.Error(t, err)
	assert.Nil(t, userInfo)
	assert.Contains(t, err.Error(), "failed to parse userInfo response")
}
