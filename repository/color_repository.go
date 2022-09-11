package repository

import (
	"net/http"

	"github.com/godemo/model"
	"github.com/godemo/model/custommodel"
	"github.com/godemo/util"
	
)

type ColorRepository struct {}

// GetALL
func (colorRepository *ColorRepository) GetAll() custommodel.ResponseDto {
	var output custommodel.ResponseDto

	db := util.CreateConnectionUsingGormToCommonSchema()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	
	var color []model.Color
	result := db.Order("color_id").Find(&color)

	if result.RowsAffected == 0 {
		output.Message = "Srver Error So no data found"
		output.IsSuccess = false
		output.Payload = nil
		output.StatusCode = http.StatusInternalServerError
		return output
	}

	type tempOutput struct {
		Output      []model.Color `json:"output"`
		OutputCount int             `json:"outputCount"`
	}
	var tOutput tempOutput
	tOutput.Output = color
	tOutput.OutputCount = len(color)
	output.Message = "Successfully Get All Colors"
	output.IsSuccess = true
	output.Payload = tOutput
	output.StatusCode = http.StatusOK
	return output
}

// GetByID
func (colorRepository *ColorRepository) GetByID(color model.Color) custommodel.ResponseDto {
	var output custommodel.ResponseDto
	if color.Color_id == 0 {
		output.Message = "Color_id is required"
		output.IsSuccess = false
		output.Payload = nil
		output.StatusCode = http.StatusBadRequest
		return output
	}

	db := util.CreateConnectionUsingGormToCommonSchema()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	result := db.Where(&model.Color{Color_id: color.Color_id}).First(&color)
	if result.RowsAffected == 0 {
		output.Message = "No data found"
		output.IsSuccess = false
		output.Payload = nil
		output.StatusCode = http.StatusNotFound
		return output
	}
	type tempOutput struct {
		Output model.Color `json:"output"`
	}
	var tOutput tempOutput
	tOutput.Output = color
	output.Message = "Color info details found for given criteria"
	output.IsSuccess = true
	output.Payload = tOutput
	output.StatusCode = http.StatusOK
	return output
}