package service

import (
	"github.com/gin-gonic/gin"
	"github.com/godemo/model"
	"github.com/godemo/repository"
)

type ColorService struct {}

func (colorService *ColorService) GetAllColorService(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Content-Type", "application/json")

	colorRepository := new(repository.ColorRepository)
	myOutput := colorRepository.GetAll()

	c.JSON(myOutput.StatusCode, myOutput)
}

func (colorService *ColorService) GetSingleColorService(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Content-Type", "application/json")
	var color model.Color
	c.ShouldBind(&color)
	colorRepository := new(repository.ColorRepository)
	myOutput := colorRepository.GetByID(color)
	c.JSON(myOutput.StatusCode, myOutput)
}



func (colorService *ColorService) AddRouters(router *gin.Engine){
	router.GET("/getallcolor", colorService.GetAllColorService)
	router.POST("/getsinglecolor", colorService.GetSingleColorService)
}