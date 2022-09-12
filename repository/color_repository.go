package repository

import (
	"net/http"
	"strings"
	"time"
	
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
	result := db.Order("color_id desc").Find(&color)

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
		output.Message = "Color ID is required"
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

// Insert Color
func (colorRepository *ColorRepository) Insert(color model.Color) custommodel.ResponseDto {
	var output custommodel.ResponseDto
	if color.Color_name == "" {
		output.Message = "Color name is required"
		output.IsSuccess = false
		output.Payload = nil
		output.StatusCode = http.StatusBadRequest
		return output
	}

	db := util.CreateConnectionUsingGormToCommonSchema()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	var output1 model.Color
	result1 := db.Where("lower(color_name) = ?", strings.ToLower(color.Color_name)).First(&output1)
	if result1.RowsAffected > 0 {
		output.Message = color.Color_name+" Color already exists"
		output.IsSuccess = false
		output.Payload = nil
		output.StatusCode = http.StatusConflict
		return output
	}
	// ID Autoincrement
	_ = db.Raw("select coalesce ((max(color_id) + 1), 1) from public.color").First(&color.Color_id)

	
	result := db.Create(&color)
	if result.RowsAffected == 0 {
		output.Message = "Color not inserted for Internal Server Error"
		output.IsSuccess = false
		output.Payload = nil
		output.StatusCode = http.StatusInternalServerError
		return output
	}
	type tempOutput struct {
		Output model.Color `json:"output"`
	}
	var tOutput tempOutput
	tOutput.Output = color
	output.Message = "Color inserted successfully"
	output.IsSuccess = true
	output.Payload = tOutput
	output.StatusCode = http.StatusOK
	return output
}

// update color
func (colorRepository *ColorRepository) Update(color model.Color) custommodel.ResponseDto {
	var output custommodel.ResponseDto
	if color.Color_id == 0 {
		output.Message = "Color ID is required"
		output.IsSuccess = false
		output.Payload = nil
		output.StatusCode = http.StatusBadRequest
		return output
	}

	if color.Color_name == "" {
		output.Message = "Color name is required"
		output.IsSuccess = false
		output.Payload = nil
		output.StatusCode = http.StatusBadRequest
		return output
	}

	db := util.CreateConnectionUsingGormToCommonSchema()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	tx := db.Begin()
	tx.SavePoint("savepoint")

	var output1 model.Color
	result1 := db.Where(&model.Color{Color_id: color.Color_id}).First(&output1)
	if result1.RowsAffected == 0 {
		output.Message = "No data found"
		output.IsSuccess = false
		output.Payload = nil
		output.StatusCode = http.StatusNotFound
		tx.SavePoint("savepoint")
		return output
	}

	var archColor model.Color_archive
	var output2 model.Color_archive
	archColor.Color_id = output1.Color_id
	archColor.Color_name = output1.Color_name
	dt := time.Now()
	archColor.Changedate = dt.Format("2006-01-02 15:04:05")
	archColor.Changeflag = "Update"

	_ = db.Raw("select coalesce ((max(trackid) + 1), 1) from public.color").First(&output2.Trackid)
	archColor.Trackid = output2.Trackid
	archColor.Changeuser = "Admin"
	result2 := tx.Create(&archColor)
	if result2.RowsAffected == 0 {
		output.Message = "Internal Server Error"
		output.IsSuccess = false
		output.Payload = nil
		output.StatusCode = http.StatusInternalServerError
		tx.RollbackTo("savepoint")
		return output
	}

	
	result3 := db.Where("lower(color_name) = ?", strings.ToLower(color.Color_name)).Where("color_id = ?", color.Color_id).First(&output1)
	if result3.RowsAffected > 0 {
		output.Message = color.Color_name+" Color already exists"
		output.IsSuccess = false
		output.Payload = nil
		output.StatusCode = http.StatusConflict
		tx.RollbackTo("savepoint")
		return output
	}

	result := db.Model(&model.Color{Color_id: color.Color_id}).Updates(&output1)
	if result.RowsAffected == 0 {
		output.Message = "Color not updated for Internal Server Error"
		output.IsSuccess = false
		output.Payload = nil
		output.StatusCode = http.StatusInternalServerError
		tx.RollbackTo("savepoint")
		return output
	}

	tx.Commit()
	type tempOutput struct {
		Output model.Color `json:"output"`
	}
	var tOutput tempOutput
	tOutput.Output = output1
	output.Message = "Color updated successfully"
	output.IsSuccess = true
	output.Payload = tOutput
	output.StatusCode = http.StatusOK
	return output
}