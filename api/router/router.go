package router

import (
	"food-tracker/controller"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func StartRouter() {

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Permitir todas as origens
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	router.GET("/orders", controller.GetOrders)
	router.POST("/orders/:id/confirm", controller.ConfirmOrder)
	router.POST("/orders/:id/startPrepare", controller.StartPreparation)
	router.POST("/orders/:id/readyToPickup", controller.ReadyToPickup)

	router.Run("localhost:8080")
}
