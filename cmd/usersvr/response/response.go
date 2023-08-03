package response

import (
	"github.com/gin-gonic/gin"
	"github.com/shixiaocaia/tiktok/cmd/gatewaysvr/log"
	"reflect"
)

const (
	successCode = 0
	errorCode   = 1
)

type response struct {
	StatusCode int32
	StatusMsg  string
}

func Response(ctx *gin.Context, httpStatus int, v interface{}) {
	ctx.JSON(httpStatus, v)
}

func Success(ctx *gin.Context, msg string, v interface{}) {
	// 记录日志，返回token信息
	if v == nil {
		Response(ctx, 200, response{successCode, msg})
	} else {
		setResponse(ctx, successCode, msg, v)
		Response(ctx, 200, v)
	}
}

func Fail(ctx *gin.Context, msg string, v interface{}) {
	if v == nil {
		Response(ctx, 200, response{errorCode, msg})
	} else {
		setResponse(ctx, errorCode, msg, v)
		Response(ctx, 200, v)
		ctx.Abort()
	}

}

// setResponse 写入日志错误信息
func setResponse(ctx *gin.Context, StatusCode int64, StatusMsg string, v interface{}) {
	getValue := reflect.ValueOf(v)
	field := getValue.Elem().FieldByName("StatusMsg")
	if field.CanSet() {
		field.SetString(StatusMsg)
	} else {
		log.Debug("cant set StatusMsg")
	}
	fieldCode := getValue.Elem().FieldByName("StatusCode")
	if fieldCode.CanSet() {
		fieldCode.SetInt(StatusCode)
	} else {
		log.Debug("cant set StatusMsg")
	}
}
