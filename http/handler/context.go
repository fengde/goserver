package handler

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/fengde/gocommon/errorx"
	"github.com/fengde/gocommon/jsonx"
	"github.com/fengde/gocommon/logx"
	"github.com/fengde/gocommon/timex"
	"github.com/fengde/gocommon/toolx"
	"github.com/gin-gonic/gin"
)

const RequestIdName = "request_id"

type Context struct {
	*gin.Context
}

// 返回请求id
func GetReqeustId(ginc *gin.Context) string {
	requestId := ginc.GetString(RequestIdName)
	if requestId == "" {
		requestId = fmt.Sprintf("%v%s", timex.NowUnixNano(), toolx.NewNumberCode(4))
		ginc.Set(RequestIdName, requestId)
	}
	return requestId
}

// controller修饰器返回gin.HandlerFunc
func WrapF(f interface{}) gin.HandlerFunc {
	return func(ginc *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				logx.ErrorWithCtx(logx.NewCtx(GetReqeustId(ginc)), r)
				_out(ginc, http.StatusInternalServerError, "fail", "internal server error", nil)
			}
		}()

		fType := reflect.TypeOf(reflect.ValueOf(f).Interface())
		argNum := fType.NumIn()
		args := make([]reflect.Value, argNum)

		args[0] = reflect.ValueOf(&Context{
			Context: ginc,
		})

		for i := 1; i < argNum; i++ {
			paramPtr := fType.In(i)
			paramKind := paramPtr.Kind()
			if paramKind == reflect.Ptr {
				argi := reflect.New(paramPtr.Elem()).Interface()
				// 解析参数到入参结构体
				if err := _params(ginc, argi); err != nil {
					_out(ginc, http.StatusOK, "fail", err.Error(), nil)
					return
				}
				// 移除入参string的首尾空格
				_trim(argi)

				// 入参tag规则校验
				if err := _paramsValidate(argi); err != nil {
					_out(ginc, http.StatusOK, "fail", err.Error(), nil)
					return
				}

				args[i] = reflect.ValueOf(argi)
			}
		}

		reflect.ValueOf(f).Call(args)
	}
}

// 协议规范通用返回
func _out(ginc *gin.Context, code int, status string, message string, data any) {
	if data == nil {
		data = map[string]any{}
	}
	ginH := gin.H{
		"status":      status,
		"message":     message,
		"data":        data,
		RequestIdName: GetReqeustId(ginc),
	}
	ginc.Set("out", jsonx.MarshalToStringNoErr(ginH))
	ginc.JSON(code, ginH)
}

// 解析json、form、uri数据，日常参看demo.go
func _params(ginc *gin.Context, r any) error {
	if err := ginc.ShouldBindQuery(r); err != nil {
		return err
	}
	if err := ginc.ShouldBind(r); err != nil {
		return err
	}
	return nil
}

// 移除入参变量，string参数的首尾空格
func _trim(r any) {
	ti := reflect.TypeOf(r)
	vi := reflect.ValueOf(r)

	if ti.Kind() == reflect.Ptr {
		realti := ti.Elem()
		realvi := vi.Elem()
		if realti.Kind() == reflect.Struct {
			for j := 0; j < realti.NumField(); j++ {
				switch realti.Field(j).Type.Kind() {
				case reflect.String:
					// 移除首尾空格
					value := realvi.Field(j).String()
					realvi.Field(j).Set(reflect.ValueOf(strings.TrimSpace(value)))
				}
			}
		}
	}
}

// 关于govalidator https://github.com/asaskevich/govalidator
func _paramsValidate(r any) error {
	_, err := govalidator.ValidateStruct(r)
	return err
}

// 返回普通文本
func (c *Context) OutString(code int, text string) {
	c.Set("out", text)
	c.String(http.StatusOK, text)
}

// 成功返回
func (c *Context) OutSuccess(data any) {
	c.Out("success", "", data)
}

// 失败返回
func (c *Context) OutFail(err error) {
	c.Out("fail", err.Error(), map[string]any{})
	// 集中打印错误堆栈
	logx.ErrorWithCtx(c.LogCtx(), "OutFail. ", errorx.GetStack(err))
}

// 提示重新登录返回
func (c *Context) OutRelogin() {
	c.Out("login", "need login", map[string]any{})
}

// 通用返回
func (c *Context) Out(status string, message string, data any) {
	_out(c.Context, http.StatusOK, status, message, data)
}

// 返回http请求id
func (c *Context) RequestId() string {
	requestId, _ := c.Get(RequestIdName)
	return fmt.Sprintf("%v", requestId)
}

// 返回日志ctx
func (c *Context) LogCtx() context.Context {
	return logx.NewCtx(c.RequestId())
}
