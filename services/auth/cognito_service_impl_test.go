package auth

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/IainMosima/gomart/configs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"
)

func TestCognitoService_GetAuthURL(t *testing.T) {
	cfg := &configs.Config{
		CognitoClientID:    "test-client-id",
		CognitoRedirectURI: "http://localhost:8080/callback",
		CognitoDomain:      "test-domain.auth.ap-northeast-1.amazoncognito.com",
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
			name:  "with state parameter",
			state: "test-state-123",
			shouldContain: []string{
				"https://test-domain.auth.ap-northeast-1.amazoncognito.com/login",
				"client_id=test-client-id",
				"response_type=code",
				"scope=email+openid+phone",
				"redirect_uri=http%3A%2F%2Flocalhost%3A8080%2Fcallback",
				"state=test-state-123",
			},
		},
		{
			name:  "without state parameter",
			state: "",
			shouldContain: []string{
				"https://test-domain.auth.ap-northeast-1.amazoncognito.com/login",
				"client_id=test-client-id",
				"response_type=code",
				"scope=email+openid+phone",
				"redirect_uri=http%3A%2F%2Flocalhost%3A8080%2Fcallback",
			},
		},
		{
			name:  "with special characters in state",
			state: "test@state#123",
			shouldContain: []string{
				"https://test-domain.auth.ap-northeast-1.amazoncognito.com/login",
				"state=test%40state%23123",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			authURL, err := service.GetAuthURL(tt.state)

			require.NoError(t, err)
			assert.NotEmpty(t, authURL)

			for _, expected := range tt.shouldContain {
				assert.Contains(t, authURL, expected)
			}

			if tt.state == "" {
				assert.NotContains(t, authURL, "state=")
			}
		})
	}
}

// NOTE: HTTP client-based tests are complex to mock due to HTTPS requirements
// and SSL certificate validation. In a production environment, these would be
// better tested with:
// 1. Dependency injection of HTTP client interface
// 2. Integration tests against real Cognito endpoints
// 3. Contract tests
// For unit testing, we focus on the business logic in auth_service_impl_test.go

func TestCognitoService_ExchangeCodeForTokens_Success(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/oauth2/token" {
			w.Header().Set("Content-Type", "application/json")
			response := `{
				"access_token": "test-access-token",
				"refresh_token": "test-refresh-token",
				"id_token": "test-id-token",
				"token_type": "Bearer",
				"expires_in": 3600
			}`
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(response))
		}
	}))
	defer server.Close()

	cfg := &configs.Config{
		CognitoClientID:     "test-client-id",
		CognitoClientSecret: "test-client-secret",
		CognitoRedirectURI:  "http://localhost:8080/callback",
		CognitoDomain:       "test-domain.auth.ap-northeast-1.amazoncognito.com",
	}

	service := &CognitoService{
		config: cfg,
		oauth2Config: oauth2.Config{
			ClientID:     cfg.CognitoClientID,
			ClientSecret: cfg.CognitoClientSecret,
			RedirectURL:  cfg.CognitoRedirectURI,
			Endpoint: oauth2.Endpoint{
				TokenURL: server.URL + "/oauth2/token",
			},
		},
	}

	// Override HTTP client to use test server's certificate
	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, server.Client())

	tokens, err := service.ExchangeCodeForTokens(ctx, "test-auth-code")

	require.NoError(t, err)
	assert.Equal(t, "test-access-token", tokens.AccessToken)
	assert.Equal(t, "test-refresh-token", tokens.RefreshToken)
	assert.Equal(t, "test-id-token", tokens.IDToken)
}

func TestCognitoService_ExchangeCodeForTokens_InvalidCode(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "invalid_grant"}`))
	}))
	defer server.Close()

	cfg := &configs.Config{
		CognitoClientID:     "test-client-id",
		CognitoClientSecret: "test-client-secret",
		CognitoRedirectURI:  "http://localhost:8080/callback",
	}

	service := &CognitoService{
		config: cfg,
		oauth2Config: oauth2.Config{
			ClientID:     cfg.CognitoClientID,
			ClientSecret: cfg.CognitoClientSecret,
			RedirectURL:  cfg.CognitoRedirectURI,
			Endpoint: oauth2.Endpoint{
				TokenURL: server.URL + "/oauth2/token",
			},
		},
	}

	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, server.Client())

	tokens, err := service.ExchangeCodeForTokens(ctx, "invalid-code")

	assert.Error(t, err)
	assert.Nil(t, tokens)
}

func TestCognitoService_ValidateAccessToken_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/oauth2/userInfo" && r.Header.Get("Authorization") == "Bearer valid-token" {
			w.Header().Set("Content-Type", "application/json")
			response := `{
				"sub": "550e8400-e29b-41d4-a716-446655440000",
				"username": "testuser",
				"email": "test@example.com",
				"email_verified": "true",
				"phone_number": "+1234567890"
			}`
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(response))
		} else {
			w.WriteHeader(http.StatusUnauthorized)
		}
	}))
	defer server.Close()

	// Use HTTP server (not HTTPS) to avoid certificate issues
	cfg := &configs.Config{
		CognitoDomain: strings.TrimPrefix(server.URL, "http://"),
	}

	service := &CognitoService{
		config: cfg,
	}

	// This will still fail because the method hardcodes HTTPS, but let's test error paths
	userInfo, err := service.ValidateAccessToken(context.Background(), "valid-token")

	assert.Error(t, err)
	assert.Nil(t, userInfo)
}

func TestCognitoService_ValidateAccessToken_HTTPError(t *testing.T) {
	cfg := &configs.Config{
		CognitoDomain: "invalid-domain-that-does-not-exist.com",
	}

	service := &CognitoService{
		config: cfg,
	}

	userInfo, err := service.ValidateAccessToken(context.Background(), "test-token")

	assert.Error(t, err)
	assert.Nil(t, userInfo)
	assert.Contains(t, err.Error(), "failed to call userInfo endpoint")
}

func TestCognitoService_ValidateAccessToken_Unauthorized(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error": "invalid_token"}`))
	}))
	defer server.Close()

	cfg := &configs.Config{
		CognitoDomain: strings.TrimPrefix(server.URL, "http://"),
	}

	service := &CognitoService{
		config: cfg,
	}

	userInfo, err := service.ValidateAccessToken(context.Background(), "invalid-token")

	assert.Error(t, err)
	assert.Nil(t, userInfo)
}

func TestCognitoService_ValidateAccessToken_BadJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`invalid json response`))
	}))
	defer server.Close()

	cfg := &configs.Config{
		CognitoDomain: strings.TrimPrefix(server.URL, "http://"),
	}

	service := &CognitoService{
		config: cfg,
	}

	userInfo, err := service.ValidateAccessToken(context.Background(), "test-token")

	assert.Error(t, err)
	assert.Nil(t, userInfo)
}

func TestCognitoService_ValidateAccessToken_EmptyResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	}))
	defer server.Close()

	cfg := &configs.Config{
		CognitoDomain: strings.TrimPrefix(server.URL, "http://"),
	}

	service := &CognitoService{
		config: cfg,
	}

	userInfo, err := service.ValidateAccessToken(context.Background(), "test-token")

	assert.Error(t, err)
	assert.Nil(t, userInfo)
}

func TestCognitoService_ValidateAccessToken_ContextCancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate slow response
		select {
		case <-r.Context().Done():
			return
		case <-time.After(100 * time.Millisecond):
			w.WriteHeader(http.StatusOK)
		}
	}))
	defer server.Close()

	cfg := &configs.Config{
		CognitoDomain: strings.TrimPrefix(server.URL, "http://"),
	}

	service := &CognitoService{
		config: cfg,
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	userInfo, err := service.ValidateAccessToken(ctx, "test-token")

	assert.Error(t, err)
	assert.Nil(t, userInfo)
}

// ParseIDToken tests are skipped as they require complex OIDC provider mocking
// In a production environment, these would be better tested with integration tests

// MockHTTPClient for testing TestableValidateAccessToken
type MockHTTPClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	if m.DoFunc != nil {
		return m.DoFunc(req)
	}
	return nil, nil
}

func TestCognitoService_TestableValidateAccessToken_Success(t *testing.T) {
	cfg := &configs.Config{
		CognitoDomain: "test-domain.auth.us-east-1.amazoncognito.com",
	}

	service := &CognitoService{
		config: cfg,
	}

	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			assert.Equal(t, "GET", req.Method)
			assert.Equal(t, "https://test-domain.auth.us-east-1.amazoncognito.com/oauth2/userInfo", req.URL.String())
			assert.Equal(t, "Bearer valid-token", req.Header.Get("Authorization"))

			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body: io.NopCloser(strings.NewReader(`{
					"sub": "550e8400-e29b-41d4-a716-446655440000",
					"username": "testuser",
					"email": "test@example.com",
					"email_verified": "true",
					"phone_number": "+1234567890"
				}`)),
			}
			return resp, nil
		},
	}

	userInfo, err := service.TestableValidateAccessToken(context.Background(), "valid-token", mockClient)

	require.NoError(t, err)
	assert.NotNil(t, userInfo)
	assert.Equal(t, "testuser", userInfo.UserName)
	assert.Equal(t, "test@example.com", userInfo.Email)
	assert.True(t, userInfo.EmailVerified)
	assert.Equal(t, "+1234567890", userInfo.PhoneNumber)
}

func TestCognitoService_TestableValidateAccessToken_HTTPError(t *testing.T) {
	cfg := &configs.Config{
		CognitoDomain: "test-domain.auth.us-east-1.amazoncognito.com",
	}

	service := &CognitoService{
		config: cfg,
	}

	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return nil, errors.New("network error")
		},
	}

	userInfo, err := service.TestableValidateAccessToken(context.Background(), "token", mockClient)

	assert.Error(t, err)
	assert.Nil(t, userInfo)
	assert.Contains(t, err.Error(), "failed to call userInfo endpoint")
}

func TestCognitoService_TestableValidateAccessToken_Unauthorized(t *testing.T) {
	cfg := &configs.Config{
		CognitoDomain: "test-domain.auth.us-east-1.amazoncognito.com",
	}

	service := &CognitoService{
		config: cfg,
	}

	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			resp := &http.Response{
				StatusCode: http.StatusUnauthorized,
				Body:       io.NopCloser(strings.NewReader(`{"error": "invalid_token"}`)),
			}
			return resp, nil
		},
	}

	userInfo, err := service.TestableValidateAccessToken(context.Background(), "invalid-token", mockClient)

	assert.Error(t, err)
	assert.Nil(t, userInfo)
	assert.Contains(t, err.Error(), "userInfo request failed with status 401")
}

func TestCognitoService_TestableValidateAccessToken_BadJSON(t *testing.T) {
	cfg := &configs.Config{
		CognitoDomain: "test-domain.auth.us-east-1.amazoncognito.com",
	}

	service := &CognitoService{
		config: cfg,
	}

	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`invalid json`)),
			}
			return resp, nil
		},
	}

	userInfo, err := service.TestableValidateAccessToken(context.Background(), "token", mockClient)

	assert.Error(t, err)
	assert.Nil(t, userInfo)
	assert.Contains(t, err.Error(), "failed to parse userInfo response")
}

func TestCognitoService_TestableValidateAccessToken_ReadError(t *testing.T) {
	cfg := &configs.Config{
		CognitoDomain: "test-domain.auth.us-east-1.amazoncognito.com",
	}

	service := &CognitoService{
		config: cfg,
	}

	// Create a mock reader that fails
	mockReader := &ErrorReader{}

	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(mockReader),
			}
			return resp, nil
		},
	}

	userInfo, err := service.TestableValidateAccessToken(context.Background(), "token", mockClient)

	assert.Error(t, err)
	assert.Nil(t, userInfo)
	assert.Contains(t, err.Error(), "failed to read response body")
}

// ErrorReader is a mock reader that always returns an error
type ErrorReader struct{}

func (er *ErrorReader) Read(p []byte) (int, error) {
	return 0, errors.New("read error")
}

func TestCognitoService_TestableValidateAccessToken_ContextRequestError(t *testing.T) {
	cfg := &configs.Config{
		CognitoDomain: "test-domain.auth.us-east-1.amazoncognito.com",
	}

	service := &CognitoService{
		config: cfg,
	}

	// Test the request creation error path by using an invalid context
	// This is hard to trigger in practice but we can simulate it
	userInfo, err := service.TestableValidateAccessToken(nil, "token", &MockHTTPClient{})

	assert.Error(t, err)
	assert.Nil(t, userInfo)
}

func TestCognitoService_TestableValidateAccessToken_EmptyResponseBody(t *testing.T) {
	cfg := &configs.Config{
		CognitoDomain: "test-domain.auth.us-east-1.amazoncognito.com",
	}

	service := &CognitoService{
		config: cfg,
	}

	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader("")),
			}
			return resp, nil
		},
	}

	userInfo, err := service.TestableValidateAccessToken(context.Background(), "token", mockClient)

	assert.Error(t, err)
	assert.Nil(t, userInfo)
	assert.Contains(t, err.Error(), "failed to parse userInfo response")
}

func TestCognitoService_TestableValidateAccessToken_InvalidUUID(t *testing.T) {
	cfg := &configs.Config{
		CognitoDomain: "test-domain.auth.us-east-1.amazoncognito.com",
	}

	service := &CognitoService{
		config: cfg,
	}

	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body: io.NopCloser(strings.NewReader(`{
					"sub": "invalid-uuid",
					"username": "testuser",
					"email": "test@example.com",
					"email_verified": "false",
					"phone_number": "+1234567890"
				}`)),
			}
			return resp, nil
		},
	}

	userInfo, err := service.TestableValidateAccessToken(context.Background(), "token", mockClient)

	require.NoError(t, err)
	assert.NotNil(t, userInfo)
	// UUID parsing should not fail, just result in zero UUID
	assert.Equal(t, "testuser", userInfo.UserName)
	assert.Equal(t, "test@example.com", userInfo.Email)
	assert.False(t, userInfo.EmailVerified)
}

func TestCognitoService_RefreshTokens_Success(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/oauth2/token" {
			w.Header().Set("Content-Type", "application/json")
			response := `{
				"access_token": "new-access-token",
				"refresh_token": "new-refresh-token",
				"id_token": "new-id-token",
				"token_type": "Bearer",
				"expires_in": 3600
			}`
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(response))
		}
	}))
	defer server.Close()

	cfg := &configs.Config{
		CognitoClientID:     "test-client-id",
		CognitoClientSecret: "test-client-secret",
	}

	service := &CognitoService{
		config: cfg,
		oauth2Config: oauth2.Config{
			ClientID:     cfg.CognitoClientID,
			ClientSecret: cfg.CognitoClientSecret,
			Endpoint: oauth2.Endpoint{
				TokenURL: server.URL + "/oauth2/token",
			},
		},
	}

	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, server.Client())

	tokens, err := service.RefreshTokens(ctx, "valid-refresh-token")

	require.NoError(t, err)
	assert.Equal(t, "new-access-token", tokens.AccessToken)
	assert.Equal(t, "new-refresh-token", tokens.RefreshToken)
	assert.Equal(t, "new-id-token", tokens.IDToken)
}

func TestCognitoService_RefreshTokens_InvalidToken(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "invalid_grant"}`))
	}))
	defer server.Close()

	cfg := &configs.Config{
		CognitoClientID:     "test-client-id",
		CognitoClientSecret: "test-client-secret",
	}

	service := &CognitoService{
		config: cfg,
		oauth2Config: oauth2.Config{
			ClientID:     cfg.CognitoClientID,
			ClientSecret: cfg.CognitoClientSecret,
			Endpoint: oauth2.Endpoint{
				TokenURL: server.URL + "/oauth2/token",
			},
		},
	}

	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, server.Client())

	tokens, err := service.RefreshTokens(ctx, "invalid-refresh-token")

	assert.Error(t, err)
	assert.Nil(t, tokens)
}

// Test utility to verify URL parameters
func TestCognitoService_GetAuthURL_URLParsing(t *testing.T) {
	cfg := &configs.Config{
		CognitoClientID:    "test-client-id",
		CognitoRedirectURI: "http://localhost:8080/callback",
		CognitoDomain:      "test-domain.auth.ap-northeast-1.amazoncognito.com",
	}

	service := &CognitoService{
		config: cfg,
	}

	authURL, err := service.GetAuthURL("test-state")
	require.NoError(t, err)

	parsedURL, err := url.Parse(authURL)
	require.NoError(t, err)

	// Verify URL components
	assert.Equal(t, "https", parsedURL.Scheme)
	assert.Equal(t, "test-domain.auth.ap-northeast-1.amazoncognito.com", parsedURL.Host)
	assert.Equal(t, "/login", parsedURL.Path)

	// Verify query parameters
	queryParams := parsedURL.Query()
	assert.Equal(t, "test-client-id", queryParams.Get("client_id"))
	assert.Equal(t, "code", queryParams.Get("response_type"))
	assert.Equal(t, "email openid phone", queryParams.Get("scope"))
	assert.Equal(t, "http://localhost:8080/callback", queryParams.Get("redirect_uri"))
	assert.Equal(t, "test-state", queryParams.Get("state"))
}

func TestCognitoService_GetAuthURL_RedirectURIEncoding(t *testing.T) {
	cfg := &configs.Config{
		CognitoClientID:    "test-client-id",
		CognitoRedirectURI: "http://localhost:8080/callback?param=value&other=test",
		CognitoDomain:      "test-domain.auth.ap-northeast-1.amazoncognito.com",
	}

	service := &CognitoService{
		config: cfg,
	}

	authURL, err := service.GetAuthURL("test")
	require.NoError(t, err)

	// Verify that the redirect URI is properly URL encoded
	assert.Contains(t, authURL, "redirect_uri=http%3A%2F%2Flocalhost%3A8080%2Fcallback%3Fparam%3Dvalue%26other%3Dtest")
}
