package auth

import (
	"net/http"

	"github.com/IainMosima/gomart/domains/auth/schema"
	"github.com/IainMosima/gomart/domains/auth/service"
	"github.com/IainMosima/gomart/rest-server/dtos"
	"github.com/gin-gonic/gin"
)

type AuthHandlerImpl struct {
	authService service.AuthService
}

func NewAuthHandlerImpl(authSvc service.AuthService) *AuthHandlerImpl {
	return &AuthHandlerImpl{
		authService: authSvc,
	}
}

func (a *AuthHandlerImpl) HandleCognitoCallback(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")

	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "code is required"})
		return
	}

	tokenResponse, err := a.authService.HandleCallback(c, &schema.HandleCallbackRequest{
		Code:  &code,
		State: state,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tokenResponse)

}

func (a *AuthHandlerImpl) ValidateToken(c *gin.Context) {
	var req dtos.ValidateTokenRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userInfo, err := a.authService.ValidateToken(c, req.AccessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, userInfo)
}

func (a *AuthHandlerImpl) RefreshAccessToken(c *gin.Context) {
	var req dtos.RefreshTokenRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	refreshTokenResponse, err := a.authService.RefreshAccessToken(c, &schema.RefreshTokenRequest{RefreshToken: req.RefreshToken})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, refreshTokenResponse)
}

func (a *AuthHandlerImpl) LoginHandler(c *gin.Context) {
	state := c.Query("state")
	loginUrl, err := a.authService.GetAuthURL(state)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusTemporaryRedirect, dtos.LoginResponseDTO{
		LoginUrl: loginUrl,
		State:    state,
		Message:  "Use the login url link to login and get the access token to be used in the app",
	})

}
