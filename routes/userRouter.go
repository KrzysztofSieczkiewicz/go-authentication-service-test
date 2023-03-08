package routes

import (
	controllers "github.com/KrzysztofSieczkiewicz/go-authentication-service-test/controllers"
	middleware "github.com/KrzysztofSieczkiewicz/go-authentication-service-test/middleware"
	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.Use(middleware.Authenticate())
	incomingRoutes.GET("users", controllers.GetUsers())
	incomingRoutes.GET("/users/:user_id", controllers.GetUser())
}
