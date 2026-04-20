package config

import (
	"errors"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/kreatip/kreatip-backend/internal/model"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	httpRequestsTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total HTTP requests by status, method, and path.",
	}, []string{"status_code", "method", "path"})

	httpRequestDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "HTTP request duration in seconds.",
		Buckets: prometheus.DefBuckets,
	}, []string{"status_code", "method", "path"})
)

func init() {
	prometheus.MustRegister(httpRequestsTotal, httpRequestDuration)
}

// RegisterMetrics mounts /metrics and adds the instrumentation middleware.
func RegisterMetrics(app *fiber.App) {
	app.Get("/metrics", adaptor.HTTPHandler(promhttp.Handler()))

	app.Use(metricsMiddleware)
}

func metricsMiddleware(c *fiber.Ctx) error {
	if c.Path() == "/metrics" {
		return c.Next()
	}

	start := time.Now()
	method := c.Method()
	path := c.Path()

	err := c.Next()

	status := resolveStatus(c, err)
	elapsed := time.Since(start).Seconds()
	code := strconv.Itoa(status)

	httpRequestsTotal.WithLabelValues(code, method, path).Inc()
	httpRequestDuration.WithLabelValues(code, method, path).Observe(elapsed)

	return err
}

// resolveStatus mirrors ErrorHandler priority to get the correct HTTP status code.
func resolveStatus(c *fiber.Ctx, err error) int {
	if err == nil {
		return c.Response().StatusCode()
	}

	var appErr *model.AppError
	if errors.As(err, &appErr) {
		return appErr.Code
	}

	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		return fiber.StatusUnprocessableEntity
	}

	var fe *fiber.Error
	if errors.As(err, &fe) {
		return fe.Code
	}

	return fiber.StatusInternalServerError
}
