package metrics

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"google.golang.org/grpc"
)

var (
	HistogramResponseTime *prometheus.HistogramVec

	CounterRequest *prometheus.CounterVec
)

func Init(service string) {
	HistogramResponseTime = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "ozon",
		Subsystem: "grpc",
		Name:      service + "_response_time_seconds",
		Buckets:   []float64{},
	},
		[]string{
			"status",
			"method",
		},
	)

	CounterRequest = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "ozon",
		Subsystem: "grpc",
		Name:      service + "_request_count",
	}, []string{"method"})
}

func MiddlewareGRPC(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	start := time.Now()

	h, err := handler(ctx, req)

	status := "ok"
	if err != nil {
		status = "error"
	}

	HistogramResponseTime.WithLabelValues(status, info.FullMethod).Observe(time.Since(start).Seconds())

	CounterRequest.WithLabelValues(info.FullMethod).Inc()

	return h, err
}
