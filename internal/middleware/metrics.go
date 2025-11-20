package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus"
	"reference-service-go/internal/metrics"
)

// HttpResponseTimeMetric wraps a Prometheus histogram for tracking HTTP response times.
type HttpResponseTimeMetric struct {
	*prometheus.HistogramVec
}

// NewHttpResponseTimeHistogramMetric creates a new HTTP response time metric with default labels for method, route, and status code.
func NewHttpResponseTimeHistogramMetric() *HttpResponseTimeMetric {
	responseTimeHistogram := metrics.NewHttpResponseTimeHistogram(
		metrics.HttpResponseTimeOpts{
			Namespace:  "app",
			LabelNames: []string{"method", "route", "code"},
		},
	)

	return &HttpResponseTimeMetric{
		HistogramVec: responseTimeHistogram,
	}
}

// ResponseTimes is a middleware that records HTTP response times and status codes to Prometheus metrics.
func (httpResponseTimeMetric *HttpResponseTimeMetric) ResponseTimes(
	next http.Handler,
) http.Handler {
	return http.HandlerFunc(
		func(responseWriter http.ResponseWriter, request *http.Request) {
			startTime := time.Now()

			responseWriterContainer := newResponseWriterWrapper(responseWriter)

			next.ServeHTTP(responseWriterContainer, request)

			statusCode := strconv.Itoa(responseWriterContainer.statusCode)
			route := getRoutePattern(request)
			duration := time.Since(startTime)
			httpResponseTimeMetric.WithLabelValues(request.Method, route, statusCode).
				Observe(duration.Seconds())
		},
	)
}

func getRoutePattern(request *http.Request) string {
	routeContext := chi.RouteContext(request.Context())
	routePattern := routeContext.RoutePattern()
	if routePattern == "" {
		return "undefined"
	}
	return routePattern
}
