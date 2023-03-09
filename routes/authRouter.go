package routes

import (
	controller "github.com/KrzysztofSieczkiewicz/go-authentication-service-test/controllers"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("users/signup", controller.SignUp())
	incomingRoutes.POST("users/login", controller.LogIn())
}
