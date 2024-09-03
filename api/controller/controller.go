package controller

import (
	"food-tracker/repositories"
	"food-tracker/services"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

var orderRepo *repositories.OrderRepository

func InitializeRepositories() {
	var err error
	orderRepo, err = repositories.NewOrderRepository()

	if err != nil {
		panic("Failed to initialize order repository: " + err.Error())
	}

	fetcherService := services.NewFetcherService(*orderRepo)
	fetcherService.Start()

	log.Println("FetcherService started")
}

func GetOrders(c *gin.Context) {
	ctx := c.Request.Context()

	orders, err := orderRepo.GetOrders(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, orders)
}

func ConfirmOrder(c *gin.Context) {
	id := c.Param("id")

	ctx := c.Request.Context()
	err := orderRepo.ConfirmOrder(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order confirmed"})
}

func StartPreparation(c *gin.Context) {
	id := c.Param("id")

	ctx := c.Request.Context()
	err := orderRepo.StartPreparation(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Start prepare"})
}

func ReadyToPickup(c *gin.Context) {
	id := c.Param("id")

	ctx := c.Request.Context()
	err := orderRepo.ReadyToPickup(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Ready To Pickup"})
}
