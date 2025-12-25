package app

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"hr-api/pkg/e"
)

type Gin struct {
	C *gin.Context
}

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// Response setting gin.JSON
func (g *Gin) Response(httpCode, errCode int, data interface{}) {
	g.C.JSON(httpCode, Response{
		Code: errCode,
		Msg:  e.GetMsg(errCode),
		Data: data,
	})
	return
}

// 200 Response data in gin.Json
func (g *Gin) SuccessResponse(data interface{}) {
	g.C.JSON(http.StatusOK, Response{
		Code: 0,
		Msg:  "OK",
		Data: data,
	})
	return
}

// 400 Response data in gin.Json
func (g *Gin) FailResponse(msg string) {
	g.C.JSON(http.StatusBadRequest, Response{
		Code: 1,
		Msg:  msg,
	})
	return
}

// 401 Response data in gin.Json
func (g *Gin) UnauthorizedResponse(msg string, data interface{}) {
	g.C.JSON(http.StatusUnauthorized, Response{
		Code: 1,
		Msg:  "Unauthorized Request",
		Data: data,
	})
	return
}

// 403 Response data in gin.Json
func (g *Gin) PermDeniedResponse(msg string) {
	g.C.JSON(http.StatusForbidden, Response{
		Code: 1,
		Msg:  msg,
	})
	return
}

// 500 Response data in gin.Json
func (g *Gin) IntervalErrorResponse(msg string) {
	g.C.JSON(http.StatusInternalServerError, Response{
		Code: 1,
		Msg:  msg,
	})
	return
}
