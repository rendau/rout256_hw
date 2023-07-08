package tracer

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/uber/jaeger-client-go/config"
	"google.golang.org/grpc"
)

func InitGlobal(jaegerHostPort, service string) error {
	cfg := config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: jaegerHostPort,
		},
	}

	if _, err := cfg.InitGlobalTracer(service); err != nil {
		return err
	}

	return nil
}

func GetTracer() opentracing.Tracer {
	return opentracing.GlobalTracer()
}

func MiddlewareGRPC(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, info.FullMethod)
	defer span.Finish()

	h, err := handler(ctx, req)
	if err != nil {
		ext.Error.Set(span, true)
	}

	return h, err
}

func MarkSpanWithError(ctx context.Context, err error) error {
	if err == nil {
		return nil
	}

	span := opentracing.SpanFromContext(ctx)
	if span == nil {
		return err
	}

	ext.Error.Set(span, true)
	span.LogKV("error", err.Error())

	return err
}
