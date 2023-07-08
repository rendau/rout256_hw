package logger

import (
	"context"
	"log"
	"strings"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
)

var (
	sl *zap.SugaredLogger
)

func Init(level string, dev bool) {
	var cfg zap.Config

	if level == "" {
		level = "info"
	}

	if dev {
		cfg = zap.NewDevelopmentConfig()
	} else {
		cfg = zap.NewProductionConfig()
		cfg.Level.SetLevel(getZapLevel(level))
	}

	cfg.EncoderConfig.MessageKey = "message"
	cfg.EncoderConfig.LevelKey = "level"
	cfg.EncoderConfig.TimeKey = "ts"
	cfg.EncoderConfig.NameKey = "logger"
	cfg.EncoderConfig.CallerKey = "caller"
	cfg.EncoderConfig.StacktraceKey = "stacktrace"
	cfg.EncoderConfig.EncodeLevel = zapcore.LowercaseLevelEncoder
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.EncoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	cfg.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

	logger, err := cfg.Build()
	if err != nil {
		log.Fatal(err)
	}

	sl = logger.Sugar()
}

func getZapLevel(v string) zapcore.Level {
	switch strings.ToLower(v) {
	case "error":
		return zap.ErrorLevel
	case "warn":
		return zap.WarnLevel
	case "info":
		return zap.InfoLevel
	case "debug":
		return zap.DebugLevel
	default:
		return zap.InfoLevel
	}
}

func Fatalw(ctx context.Context, err error, msg string, args ...any) {
	withTraceID(ctx).Desugar().
		With(zap.String("error", err.Error())).Sugar().
		Fatalw(msg, args...)
}

func Errorw(ctx context.Context, err error, msg string, args ...any) {
	withTraceID(ctx).Desugar().
		With(zap.String("error", err.Error())).Sugar().
		Errorw(msg, args...)
}

func Warnw(ctx context.Context, msg string, args ...any) {
	withTraceID(ctx).Warnw(msg, args...)
}

func Infow(ctx context.Context, msg string, args ...any) {
	withTraceID(ctx).Infow(msg, args...)
}

func Debugw(ctx context.Context, msg string, args ...any) {
	withTraceID(ctx).Debugw(msg, args...)
}

func withTraceID(ctx context.Context) *zap.SugaredLogger {
	if ctx == nil {
		return sl
	}

	span := opentracing.SpanFromContext(ctx)
	if span == nil {
		return sl
	}

	if sc, ok := span.Context().(jaeger.SpanContext); ok {
		return sl.Desugar().With(
			zap.Stringer("trace_id", sc.TraceID()),
			zap.Stringer("span_id", sc.SpanID()),
		).Sugar()
	}

	return sl
}

func MiddlewareGRPC(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	h, err := handler(ctx, req)
	if err != nil {
		Errorw(ctx, err, "error while processing handler",
			"method", info.FullMethod,
		)
	}

	return h, err
}
