package repository

import (
	"net/http"

	"github.com/godemo/model"
	"github.com/godemo/model/cmodel"
	"github.com/godemo/util"
	
)

type ColorRepository struct {}

// GetALL
func (colorRepository *ColorRepository) GetAll() cmodel.ResponseDto {
	var output cmodel.ResponseDto

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