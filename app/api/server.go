package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gkazioka/clientsapi/app/controllers"
	"github.com/gkazioka/clientsapi/app/middlewares"
	"github.com/gkazioka/clientsapi/app/services"
	"github.com/gkazioka/clientsapi/app/utils"
)

func GetServer() {
	userService := services.GetInstance("postgres")

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(middlewares.GetMetrics())
	registerRoutes(r, userService)
	r.Run()
	utils.Initialize()
}

func registerRoutes(r *gin.Engine, userService services.UserService) {
	userController := controllers.NewUserController(userService)
	r.GET("/status", userController.APIStatus)
	r.GET("/clients", userController.ListClients)
	r.GET("/clients/:code", userController.FindByCode)
	r.POST("/clients", userController.AddClients)
}
