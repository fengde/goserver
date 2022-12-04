package middleware

import (
	"context"
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
	return func(c *gin.Context) {
		if global.Conf.Jaeger.Way == "http" {
			tracer, closer, err := getJaegerConfig().NewTracer(jaegerConfig.Logger(jaeger.StdLogger))
			if err != nil {
				logx.Error(err)
				c.Next()
				return
			}
			defer closer.Close()

			xRequestId := c.GetString(XRequestIdHeader)

			var parentSpan opentracing.Span

			header := http.Header{}
			header.Add(XRequestIdHeader, xRequestId)

			lastCtx, err := tracer.Extract(opentracing.HTTPHeaders, header)
			if err != nil {
				parentSpan = tracer.StartSpan(
					"recv request",
					opentracing.Tags{"http.url": c.Request.URL.Path, XRequestIdHeader: xRequestId},
				)
			} else {
				parentSpan = tracer.StartSpan(
					"recv request",
					opentracing.Tags{"url": c.Request.URL.Path, XRequestIdHeader: xRequestId},
					opentracing.ChildOf(lastCtx),
				)
			}

			defer parentSpan.Finish()

			oldCtx, exist := c.Get("ctx")
			if exist {
				c.Set("ctx", opentracing.ContextWithSpan(oldCtx.(context.Context), parentSpan))
			} else {
				c.Set("ctx", opentracing.ContextWithSpan(context.Background(), parentSpan))
			}

			c.Next()
		}
	}
}
