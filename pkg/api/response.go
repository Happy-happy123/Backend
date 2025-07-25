package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    http.StatusOK,
		Message: "success",
		Data:    data,
	})
}

func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, Response{
		Code:    http.StatusCreated,
		Message: "created",
		Data:    data,
	})
}

func BadRequest(c *gin.Context, msg string) {
	c.JSON(http.StatusBadRequest, Response{
		Code:    http.StatusBadRequest,
		Message: msg,
	})
}

func NotFound(c *gin.Context, msg string) {
	c.JSON(http.StatusNotFound, Response{
		Code:    http.StatusNotFound,
		Message: msg,
	})
}

func ServerError(c *gin.Context, msg string) {
	c.JSON(http.StatusInternalServerError, Response{
		Code:    http.StatusInternalServerError,
		Message: msg,
	})
}