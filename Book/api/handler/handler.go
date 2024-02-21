package handler

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"kitab/api/models"
	"kitab/service"
)

type Handler struct {
	services service.IServiceManager
}

func New(services service.IServiceManager) Handler {
	return Handler{
		services: services,
	}
}

func handleResponse(c *gin.Context, msg string, statusCode int, data interface{}) {
	resp := models.Response{}

	switch code := statusCode; {
	case code < 400:
		resp.Description = "OK"
		fmt.Println("~~~~> OK", "msg", "status", code)
	case code == 401:
		resp.Description = "Unauthorized"
	case code < 500:
		resp.Description = "Bad Request"
		fmt.Println("!!!!! BAD REQUEST", "msg", "status", code)
	default:
		resp.Description = "Internal Server Error"
		fmt.Println("!!!!! INTERNAL SERVER ERROR", "msg", msg, "status", code)
	}

	resp.StatusCode = statusCode
	resp.Data = data

	c.JSON(resp.StatusCode, resp)
}
