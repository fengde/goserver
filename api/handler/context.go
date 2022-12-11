package handler

import (
	"context"
	"net/http"
	"reflect"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/fengde/gocommon/errorx"
	"github.com/fengde/gocommon/jsonx"
	"github.com/fengde/gocommon/logx"
	"github.com/gin-gonic/gin"
)

type Context struct {
	*gin.Context
	UserId int64
}

// 返回日志ctx
func (c *Context) GetCtx() context.Context {
	return GetCtx(c.Context)
}

// controller修饰器返回gin.HandlerFunc
func WrapF(f interface{}) gin.HandlerFunc {
	return func(ginc *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				logx.ErrorWithCtx(GetCtx(ginc), r)
				_out(ginc, http.StatusInternalServerError, "failed", "internal server error", nil)
			}
		}()

		fType := reflect.TypeOf(reflect.ValueOf(f).Interface())
		argNum := fType.NumIn()
		args := make([]reflect.Value, argNum)

		ctx := &Context{
			Context: ginc,
			UserId:  GetUserId(ginc),
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
				logx.ErrorWithCtx(ctx.GetCtx(), "error trace: ", errorx.GetStack(responseErr.Interface().(error)))
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
		"status":     status,
		"message":    message,
		"data":       data,
		"request_id": GetRequestId(ginc),
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

func GetCtx(ginc *gin.Context) context.Context {
	t, exist := ginc.Get("ctx")
	if !exist {
		return context.Background()
	}
	return t.(context.Context)
}

func GetUserId(ginc *gin.Context) int64 {
	return ginc.GetInt64("user_id")
}

func GetRequestId(ginc *gin.Context) string {
	return ginc.GetString("request_id")
}
