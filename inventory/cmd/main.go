package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"slices"
	"sync"
	"syscall"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/Artyom099/factory/inventory/internal/utils"
	inventoryV1 "github.com/Artyom099/factory/shared/pkg/proto/inventory/v1"
)

const grpcPort = 50051

type inventoryService struct {
	inventoryV1.UnimplementedInventoryServiceServer

	mu    sync.RWMutex
	parts map[string]*inventoryV1.Part
}

func (s *inventoryService) GetPart(_ context.Context, req *inventoryV1.GetPartRequest) (*inventoryV1.GetPartResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	partUuid := req.GetUuid()
	part, ok := s.parts[partUuid]

	if !ok {
		return nil, status.Errorf(codes.NotFound, "Part with UUID %s not found", partUuid)
	}

	return &inventoryV1.GetPartResponse{
		Part: part,
	}, nil
}

func (s *inventoryService) ListParts(_ context.Context, req *inventoryV1.ListPartsRequest) (*inventoryV1.ListPartsResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	filter := req.GetFilter()
	var result []*inventoryV1.Part

	for _, part := range s.parts {
		if matchPart(part, filter) {
			result = append(result, part)
		}
	}

	return &inventoryV1.ListPartsResponse{
		Parts: result,
	}, nil
}

func matchPart(p *inventoryV1.Part, f *inventoryV1.PartsFilter) bool {
	if f == nil {
		return true // фильтр пустой → вернуть всё
	}

	// фильтр по UUID
	if len(f.Uuids) > 0 && !slices.Contains(f.GetUuids(), p.GetUuid()) {
		return false
	}

	// фильтр по имени
	if len(f.Names) > 0 && !utils.ContainsInsensitive(f.GetNames(), p.GetName()) {
		return false
	}

	// фильтр по категории
	if len(f.Categories) > 0 && !slices.Contains(f.GetCategories(), p.GetCategory()) {
		return false
	}

	// фильтр по стране
	if len(f.ManufacturerCountries) > 0 && !utils.ContainsInsensitive(f.GetManufacturerCountries(), p.Manufacturer.GetCountry()) {
		return false
	}

	// фильтр по тегам (теги в Part должны пересекаться с фильтром)
	if len(f.Tags) > 0 && !utils.Intersects(f.GetTags(), p.GetTags()) {
		return false
	}

	return true
}

func (s *inventoryService) CreatePart(_ context.Context, req *inventoryV1.CreatePartRequest) (*inventoryV1.CreatePartResponse, error) {
	p := &inventoryV1.Part{
		Uuid:        uuid.New().String(),
		Name:        req.GetName(),
		Category:    req.GetCategory(),
		Description: req.GetDescription(),
		Dimensions: &inventoryV1.Dimensions{
			Length: req.GetDimensions().GetLength(),
			Width:  req.GetDimensions().GetWidth(),
			Height: req.GetDimensions().GetHeight(),
			Weight: req.GetDimensions().GetWeight(),
		},
		Price:         req.GetPrice(),
		StockQuantity: req.GetStockQuantity(),
		Manufacturer: &inventoryV1.Manufacturer{
			Name:    req.GetManufacturer().GetName(),
			Country: req.GetManufacturer().GetCountry(),
		},
		Tags:      req.GetTags(),
		CreatedAt: timestamppb.Now(),
		UpdatedAt: nil,
	}

	s.parts[p.GetUuid()] = p

	return &inventoryV1.CreatePartResponse{
		Uuid: p.GetUuid(),
	}, nil
}

func main() {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Printf("failed to listen: %v\n", err)
		return
	}
	defer func() {
		if cerr := listener.Close(); cerr != nil {
			log.Printf("failed to close listener: %v\n", cerr)
		}
	}()

	// Создаем gRPC сервер
	s := grpc.NewServer()

	// Регистрируем наш сервис
	service := &inventoryService{
		parts: make(map[string]*inventoryV1.Part),
	}

	inventoryV1.RegisterInventoryServiceServer(s, service)

	// Включаем рефлексию для отладки
	reflection.Register(s)

	go func() {
		log.Printf("🚀 gRPC server listening on %d\n", grpcPort)
		err = s.Serve(listener)
		if err != nil {
			log.Printf("failed to serve: %v\n", err)
			return
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑 Shutting down gRPC server...")
	s.GracefulStop()
	log.Println("✅ Server stopped")
}
