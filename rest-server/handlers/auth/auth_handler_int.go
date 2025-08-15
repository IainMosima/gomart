package auth

import "github.com/gin-gonic/gin"

type AuthHandler interface {
	HandleCognitoCallback(c *gin.Context)
	ValidateToken(c *gin.Context)
	RefreshAccessToken(c *gin.Context)
	LoginHandler(c *gin.Context)
}
