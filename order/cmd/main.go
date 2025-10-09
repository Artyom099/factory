package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"

	partV1 "github.com/Artyom099/shared/pkg/openapi/part/v1"
)

const (
	httpPort     = "8080"
	urlParamCity = "city"
	// Таймауты для HTTP-сервера
	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 10 * time.Second
)

// PartStorage представляет потокобезопасное хранилище данных о деталях
type PartStorage struct {
	mu     sync.RWMutex
	orders map[string]*partV1.Weather
}

// NewPartStorage создает новое хранилище данных о деталях
func NewPartStorage() *PartStorage {
	return &PartStorage{
		orders: make(map[string]*partV1.Weather),
	}
}

func (s *PartStorage) GetOrder(uuid string) *partV1.Weather {
	s.mu.RLock()
	defer s.mu.RUnlock()

	order, ok := s.orders[uuid]
	if !ok {
		return nil
	}

	return order
}

// UpdateWeather обновляет данные о погоде для указанного города
func (s *PartStorage) CreateOrder(city string, weather *weatherV1.Weather) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.weathers[city] = weather
}

func main() {
	storage := NewPartStorage()

	// Создаем обработчик API погоды
	weatherHandler := NewWeatherHandler(storage)

	// Создаем OpenAPI сервер
	weatherServer, err := weatherV1.NewServer(weatherHandler)
	if err != nil {
		log.Fatalf("ошибка создания сервера OpenAPI: %v", err)
	}

	// Инициализируем роутер Chi
	r := chi.NewRouter()

	// Добавляем middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))
	r.Use(render.SetContentType(render.ContentTypeJSON))

	// Определяем маршруты
	r.Route("/api/v1/orders", func(r chi.Router) {
		r.Post("/", CreateOrder(storage))                   // создание заказа
		r.Get("/{order_uuid}/pay", PayOrder(storage))       // оплата заказа
		r.Put("/{order_uuid}", GetOrder(storage))           // получить заказ по UUID
		r.Put("/{order_uuid}/cancel", CancelOrder(storage)) // отменить заказ
	})

	// Запускаем HTTP-сервер
	server := &http.Server{
		Addr:              net.JoinHostPort("localhost", httpPort),
		Handler:           r,
		ReadHeaderTimeout: readHeaderTimeout, // Защита от Slowloris атак - тип DDoS-атаки
	}

	// Запускаем сервер в отдельной горутине
	go func() {
		log.Printf("🚀 HTTP-сервер запущен на порту %s\n", httpPort)
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("❌ Ошибка запуска сервера: %v\n", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Завершение работы сервера...")

	// Создаем контекст с таймаутом для остановки сервера
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		log.Printf("❌ Ошибка при остановке сервера: %v\n", err)
	}

	log.Println("✅ Сервер остановлен")
}
