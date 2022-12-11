package middleware

import (
	"goserver/api/handler"
	"goserver/global"
	"net/http"

	"github.com/fengde/gocommon/logx"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegerConfig "github.com/uber/jaeger-client-go/config"
)

func getJaegerConfig() *jaegerConfig.Configuration {
	return &jaegerConfig.Configuration{
		ServiceName: global.Conf.Name,
		Sampler: &jaegerConfig.SamplerConfig{
			Type:  global.Conf.Jaeger.SamplerType,
			Param: global.Conf.Jaeger.SamplerParam,
		},
		Reporter: &jaegerConfig.ReporterConfig{
			LogSpans:          true,
			CollectorEndpoint: global.Conf.Jaeger.Endpoint,
		},
		Tags: []opentracing.Tag{
			{Key: "token", Value: global.Conf.Jaeger.Token},
		},
	}
}

// apm监控
func Jaeger() gin.HandlerFunc {
	return func(ginc *gin.Context) {
		if global.Conf.Jaeger.Way == "http" {
			tracer, closer, err := getJaegerConfig().NewTracer(jaegerConfig.Logger(jaeger.StdLogger))
			if err != nil {
				logx.Error(err)
				ginc.Next()
				return
			}
			defer closer.Close()

			xRequestId := handler.GetRequestId(ginc)

			var parentSpan opentracing.Span

			header := http.Header{}
			header.Add(RequestIdHeader, xRequestId)

			lastCtx, err := tracer.Extract(opentracing.HTTPHeaders, header)
			if err != nil {
				parentSpan = tracer.StartSpan(
					"recv request",
					opentracing.Tags{"http.url": ginc.Request.URL.Path, RequestIdHeader: xRequestId},
				)
			} else {
				parentSpan = tracer.StartSpan(
					"recv request",
					opentracing.Tags{"http.url": ginc.Request.URL.Path, RequestIdHeader: xRequestId},
					opentracing.ChildOf(lastCtx),
				)
			}

			defer parentSpan.Finish()

			ginc.Set("ctx", opentracing.ContextWithSpan(handler.GetCtx(ginc), parentSpan))

			ginc.Next()
		}
	}
}
