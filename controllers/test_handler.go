package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func DemoHandler(context *gin.Context) {
	context.String(http.StatusOK, "Hello visitor!")
}
