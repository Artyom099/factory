package interceptor

import (
	"log"
	"net/http"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"

	orderMetrics "github.com/Artyom099/factory/order/internal/metrics"
)

// responseWriter wraps http.ResponseWriter to capture status code.
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// MetricsInterceptor создает HTTP middleware для записи метрик. Возвращает функцию,
// которая принимает следующий http.Handler и возвращает http.Handler с логикой метрик.
func MetricsInterceptor() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Засекаем время начала запроса
			start := time.Now()

			// Обёртка для захвата статуса ответа
			wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

			// Выполняем запрос
			next.ServeHTTP(wrapped, r)

			// Записываем время выполнения
			duration := time.Since(start)
			durationSeconds := duration.Seconds()
			method := r.Method + " " + r.URL.Path

			log.Printf("🕐 Request duration: %v (%f seconds) for method: %s", duration, durationSeconds, method)

			orderMetrics.RequestDuration.Record(r.Context(), durationSeconds,
				metric.WithAttributes(
					attribute.String("method", method),
				),
			)

			// Определяем статус ответа по коду HTTP
			statusLabel := "success"
			if wrapped.statusCode >= 400 {
				statusLabel = "error"
			}

			// Записываем метрику запроса
			orderMetrics.RequestsTotal.Add(r.Context(), 1,
				metric.WithAttributes(
					attribute.String("method", method),
					attribute.String("status", statusLabel),
				),
			)
		})
	}
}
