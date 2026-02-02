package metrics

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

const (
	serviceName = "order-service"
)

// Важно: meter должен быть создан один раз и переиспользоваться в рамках компонента
var meter = otel.Meter(serviceName)

var (
	// RequestsTotal - COUNTER для подсчета общего количества запросов
	// Тип: Int64Counter (монотонно возрастающий)
	// Использование: подсчет всех gRPC запросов с разбивкой по методам и статусам
	// Лейблы: method (название метода), status (success/error)
	RequestsTotal metric.Int64Counter

	// RequestDuration - HISTOGRAM для измерения времени выполнения запросов
	// Тип: Float64Histogram (распределение значений)
	// Использование: SLA мониторинг - отслеживание времени ответа API
	// Позволяет строить percentile (p50, p95, p99) для анализа производительности
	RequestDuration metric.Float64Histogram
)

func InitMetrics() error {
	var err error

	RequestsTotal, err = meter.Int64Counter(
		"assembly_requests_total", // id: 1, 2 - ??
		metric.WithDescription("Total number of Order service requests"),
	)
	if err != nil {
		return err
	}

	// Создаем гистограмму времени запросов с правильными bucket'ами для gRPC
	// Bucket'ы оптимизированы для времени отклика в диапазоне от микросекунд до секунд
	RequestDuration, err = meter.Float64Histogram(
		"assembly_request_duration_seconds", // id: 3 - ??
		metric.WithDescription("Duration of gRPC requests"),
		metric.WithUnit("s"),
		// Добавляем explicit bucket boundaries для более точного измерения gRPC запросов
		// 1ms, 2ms, 5ms, 10ms, 25ms, 50ms, 100ms, 250ms, 500ms, 1s, 2s, 5s
		metric.WithExplicitBucketBoundaries(
			0.001, 0.002, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1.0, 2.0, 5.0,
		),
	)
	if err != nil {
		return err
	}

	return nil
}
