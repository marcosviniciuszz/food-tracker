package main

import (
	"food-tracker/controller"
	"food-tracker/database"
	"food-tracker/router"
)

func main() {

	database.ConnectMongoDB()

	controller.InitializeRepositories()

	router.StartRouter()
}
