package integration

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/Artyom099/factory/platform/pkg/logger"
	inventoryV1 "github.com/Artyom099/factory/shared/pkg/proto/inventory/v1"
)

var _ = Describe("InventoryService", Ordered, func() {
	var (
		ctx             context.Context
		cancel          context.CancelFunc
		inventoryClient inventoryV1.InventoryServiceClient
	)

	BeforeEach(func() {
		ctx, cancel = context.WithCancel(suiteCtx)

		// Создаём gRPC клиент
		conn, err := grpc.NewClient(
			env.App.Address(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		Expect(err).ToNot(HaveOccurred(), "ожидали успешное подключение к gRPC приложению")

		inventoryClient = inventoryV1.NewInventoryServiceClient(conn)
	})

	AfterEach(func() {
		// Чистим коллекцию после теста
		err := env.ClearPartsCollection(ctx)
		Expect(err).ToNot(HaveOccurred(), "ожидали успешную очистку коллекции parts")

		cancel()
	})

	Describe("Create", func() {
		It("должен успешно создавать новую деталь", func() {
			req := env.GetCreatePartRequest()

			resp, err := inventoryClient.CreatePart(ctx, req)

			Expect(err).ToNot(HaveOccurred())
			Expect(resp.GetUuid()).ToNot(BeEmpty())
			Expect(resp.GetUuid()).To(MatchRegexp(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`))
		})
	})

	Describe("Get", func() {
		var partUUID string

		BeforeEach(func() {
			var err error
			partUUID, err = env.InsertTestPart(ctx)
			Expect(err).ToNot(HaveOccurred(), "ожидали успешную вставку детали в MongoDB")
		})

		It("должен успешно возвращать деталь по UUID", func() {
			resp, err := inventoryClient.GetPart(ctx, &inventoryV1.GetPartRequest{
				Uuid: partUUID,
			})

			Expect(err).ToNot(HaveOccurred())

			Expect(resp.GetPart()).ToNot(BeNil())
			Expect(resp.GetPart().Uuid).To(Equal(partUUID))
			Expect(resp.GetPart().Name).ToNot(BeEmpty())
			Expect(resp.GetPart().Description).ToNot(BeEmpty())
			// Expect(resp.GetPart().Price).ToNot(BeEmpty())
			// Expect(resp.GetPart().StockQuantity).ToNot(BeEmpty())
			// Expect(resp.GetPart().Category).To(Equal(inventoryV1.Category_CATEGORY_ENGINE))

			Expect(resp.GetPart().GetDimensions()).ToNot(BeNil())
			// Expect(resp.GetPart().GetDimensions().Width).ToNot(BeEmpty())
			// Expect(resp.GetPart().GetDimensions().Height).ToNot(BeEmpty())
			// Expect(resp.GetPart().GetDimensions().Length).ToNot(BeEmpty())
			// Expect(resp.GetPart().GetDimensions().Weight).ToNot(BeEmpty())

			Expect(resp.GetPart().GetCreatedAt()).ToNot(BeNil())
		})
	})

	Describe("List", Ordered, func() {
		var partUUID1, partUUID2, partUUID3 string

		BeforeEach(func() {
			var err error
			partUUID1, err = env.InsertTestPart(ctx)
			partUUID2, err = env.InsertTestPart(ctx)
			partUUID3, err = env.InsertTestPart(ctx)
			Expect(err).ToNot(HaveOccurred(), "ожидали успешную вставку детали в MongoDB")
		})

		It("list all parts", func() {
			resp, err := inventoryClient.ListParts(ctx, &inventoryV1.ListPartsRequest{
				Filter: &inventoryV1.PartsFilter{},
			})

			Expect(err).ToNot(HaveOccurred())
			Expect(resp.GetParts()).ToNot(BeNil())
			Expect(len(resp.GetParts())).To(Equal(3))

			Expect(resp.GetParts()[0].Uuid).To(Equal(partUUID1))
			Expect(resp.GetParts()[1].Uuid).To(Equal(partUUID2))
			Expect(resp.GetParts()[2].Uuid).To(Equal(partUUID3))
		})

		XIt("list part by uuid", func() {
			resp, err := inventoryClient.ListParts(ctx, &inventoryV1.ListPartsRequest{
				Filter: &inventoryV1.PartsFilter{
					Uuids: []string{partUUID1, partUUID3},
				},
			})

			Expect(err).ToNot(HaveOccurred())
			parts := resp.GetParts()
			logger.Debug(ctx, "", zap.Any("Parts: ", parts))
			Expect(len(parts)).To(Equal(2))

			Expect(parts[0].Uuid).To(Or(Equal(partUUID1), Equal(partUUID3)))
			Expect(parts[1].Uuid).To(Or(Equal(partUUID1), Equal(partUUID3)))
		})
	})
})
