package main

import (
	"net/http"

	"github.com/exercise/beer/connection"
	"github.com/exercise/beer/handler"
	"github.com/exercise/beer/model"
	"github.com/exercise/beer/service"
	"github.com/gin-gonic/gin"
)

func main() {
	mongoDB := connection.NewMongoDB()
	mariaDB := connection.NewMariaDB()

	beerService := service.NewBeerService(mongoDB, mariaDB)
	beerHandler := handler.NewBeerHandler(beerService)

	mariaDB.AutoMigrate(&model.Beer{})

	r := gin.Default()

	r.GET("/beers", beerHandler.FindAll)
	r.POST("/beers", beerHandler.Create)
	r.POST("/beersimage", beerHandler.Upload)
	r.DELETE("/beers/:beerId", beerHandler.Delete)
	r.PUT("/beers/:beerId", beerHandler.Update)
	r.StaticFS("/file", http.Dir("images"))
	r.Run()

}
