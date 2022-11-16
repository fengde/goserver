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

// 返回http请求id
func (c *Context) RequestId() string {
	requestId, _ := c.Get(RequestIdName)
	return fmt.Sprintf("%v", requestId)
}

// 返回日志ctx
func (c *Context) LogCtx() context.Context {
	return logx.NewCtx(c.RequestId())
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
				_out(ginc, http.StatusInternalServerError, "failed", "internal server error", nil)
			}
		}()

		fType := reflect.TypeOf(reflect.ValueOf(f).Interface())
		argNum := fType.NumIn()
		args := make([]reflect.Value, argNum)

		ctx := &Context{
			Context: ginc,
		}
		args[0] = reflect.ValueOf(ctx)

		for i := 1; i < argNum; i++ {
			paramPtr := fType.In(i)
			paramKind := paramPtr.Kind()
			if paramKind == reflect.Ptr {
				argi := reflect.New(paramPtr.Elem()).Interface()
				// 解析参数到入参结构体
				if err := _params(ginc, argi); err != nil {
					_out(ginc, http.StatusOK, "failed", err.Error(), nil)
					return
				}
				// 移除入参string的首尾空格
				_trim(argi)

				// 入参tag规则校验
				if err := _paramsValidate(argi); err != nil {
					_out(ginc, http.StatusOK, "failed", err.Error(), nil)
					return
				}

				args[i] = reflect.ValueOf(argi)
			}
		}

		// 对handler的返回参数进行封装返回：
		// 1) 返回 error
		// 2) 返回 *struct, error
		reponseArgs := reflect.ValueOf(f).Call(args)

		var responseData, responseErr *reflect.Value

		switch {
		case len(reponseArgs) == 1:
			responseErr = &reponseArgs[0]
		case len(reponseArgs) == 2:
			responseData, responseErr = &reponseArgs[0], &reponseArgs[1]
		}

		if responseData != nil && responseData.Kind() != reflect.Ptr {
			panic("handler return invalied")
		}

		if responseErr != nil && responseErr.Type().Name() != "error" {
			panic("handler return invalied")
		}

		if responseErr != nil {
			if responseErr.IsNil() {
				if responseData != nil {
					_out(ginc, http.StatusOK, "success", "", responseData.Interface())
				} else {
					_out(ginc, http.StatusOK, "success", "", nil)
				}
			} else {
				logx.ErrorWithCtx(ctx.LogCtx(), "error trace: ", errorx.GetStack(responseErr.Interface().(error)))
				_out(ginc, http.StatusOK, "failed", responseErr.MethodByName("Error").Call(nil)[0].String(), nil)
			}
		}
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
