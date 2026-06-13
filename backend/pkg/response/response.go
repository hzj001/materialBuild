package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Body struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func OK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Body{Code: 0, Message: "success", Data: data})
}

func OKMsg(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, Body{Code: 0, Message: msg})
}

func Fail(c *gin.Context, code int, msg string) {
	c.JSON(http.StatusOK, Body{Code: code, Message: msg})
}

func Unauthorized(c *gin.Context, msg string) {
	c.JSON(http.StatusUnauthorized, Body{Code: 401, Message: msg})
}

func PageData(list interface{}, total int64, page, pageSize int) gin.H {
	return gin.H{
		"list":      list,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	}
}
