package service

import (
	"github.com/gin-gonic/gin"
	"github.com/godemo/repository"
)

type ColorService struct {}

func (colorService *ColorService) GetAll(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Content-Type", "application/json")

	colorRepository := new(repository.ColorRepository)
	myOutput := colorRepository.GetAll()

	c.JSON(myOutput.StatusCode, myOutput)
}


func (colorService *ColorService) AddRouters(router *gin.Engine){
	router.GET("/getall", colorService.GetAll)
}