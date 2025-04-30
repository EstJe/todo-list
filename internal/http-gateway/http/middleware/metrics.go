package middleware

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// Общее количество запросов
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	// Время ответа с квантилями
	httpRequestDuration = promauto.NewSummaryVec(
		prometheus.SummaryOpts{
			Name:       "http_request_duration_seconds",
			Help:       "Duration of HTTP requests with quantiles",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.95: 0.005, 0.99: 0.001},
		},
		[]string{"method", "path", "status"},
	)
)

func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// ResponseWriter для записи статуса
		rw := &responseWriter{w, http.StatusOK}

		// Выполняем следующий обработчик
		next.ServeHTTP(rw, r)

		// Получаем параметры для метрик
		status := http.StatusText(rw.status)
		method := r.Method
		path := r.URL.Path
		duration := time.Since(start).Seconds()

		// Обновляем метрики
		httpRequestsTotal.WithLabelValues(method, path, status).Inc()
		httpRequestDuration.WithLabelValues(method, path, status).Observe(duration)
	})
}

// responseWriter запоминает статус код ответа
type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}
