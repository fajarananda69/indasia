package main

import (
	"indasia/registerLogin/controllers"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/indasia/login", controllers.Login)
	router.GET("/indasia/forgot", controllers.ForgotPass)
	router.GET("/indasia/newpass", controllers.SetNewPass)
	router.GET("/indasia/setregister", controllers.SetRegister)
	// router.GET("/indasia/getregister/:key", controllers.GetRegister)
	router.GET("/indasia/validateToken/:key", controllers.ValidateToken)

	router.Run(":3000")

}
