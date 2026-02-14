package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Result 统一返回结构体
type Result struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

// Success 成功返回
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Result{
		Code: 200,
		Msg:  "success",
		Data: data,
	})
}

// Error 错误返回
func Error(c *gin.Context, code int, msg string) {
	c.JSON(http.StatusOK, Result{
		Code: code,
		Msg:  msg,
		Data: nil,
	})
}
