package routes

import (
	"github.com/IainMosima/gomart/rest-server/handlers/auth"
	"github.com/gin-gonic/gin"
)

func SetupAuthRoutes(router *gin.Engine, authHandler auth.AuthHandlerInterface) {
	router.GET("/cognito/callback", authHandler.HandleCognitoCallback)

	authGroup := router.Group("/auth")
	{
		authGroup.POST("/validate", authHandler.ValidateToken)
		authGroup.POST("/refresh", authHandler.RefreshAccessToken)
		authGroup.GET("/login", authHandler.LoginHandler)
	}
}
