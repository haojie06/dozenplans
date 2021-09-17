package controllers

import (
	"dozenplans/models"

	"github.com/gin-gonic/gin"
)

func returnResult(code int, msg string, data interface{}, context *gin.Context) {
	result := models.Result{
		Code:    code,
		Message: msg,
		Data:    data,
	}
	context.JSON(code, gin.H{"result": result})
}
