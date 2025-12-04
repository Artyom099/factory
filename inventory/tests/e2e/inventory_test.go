package integration

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

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

		md := metadata.Pairs("session-uuid", "test-session")
		ctx = metadata.NewOutgoingContext(ctx, md)

		conn, err := grpc.NewClient(
			env.App.Address(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		Expect(err).ToNot(HaveOccurred(), "ожидали успешное подключение к gRPC приложению")

		inventoryClient = inventoryV1.NewInventoryServiceClient(conn)
	})

	AfterEach(func() {
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
		var insertedPartPtr *InsertedPart

		BeforeEach(func() {
			var err error
			insertedPartPtr, err = env.InsertTestPart(ctx)
			Expect(err).ToNot(HaveOccurred(), "ожидали успешную вставку детали в MongoDB")
		})

		It("должен успешно возвращать деталь по UUID", func() {
			resp, err := inventoryClient.GetPart(ctx, &inventoryV1.GetPartRequest{
				Uuid: insertedPartPtr.Uuid,
			})

			Expect(err).ToNot(HaveOccurred())
			Expect(resp.GetPart()).ToNot(BeNil())

			AssertPartsEqual(insertedPartPtr, resp.GetPart())
		})
	})

	Describe("List", Ordered, func() {
		var insertedPartPtr1, insertedPartPtr2, insertedPartPtr3 *InsertedPart

		BeforeEach(func() {
			var err error

			insertedPartPtr1, err = env.InsertTestPart(ctx)
			Expect(err).ToNot(HaveOccurred(), "ожидали успешную вставку детали 1 в MongoDB")

			insertedPartPtr2, err = env.InsertTestPart(ctx)
			Expect(err).ToNot(HaveOccurred(), "ожидали успешную вставку детали 2 в MongoDB")

			insertedPartPtr3, err = env.InsertTestPart(ctx)
			Expect(err).ToNot(HaveOccurred(), "ожидали успешную вставку детали 3 в MongoDB")
		})

		It("list all parts", func() {
			resp, err := inventoryClient.ListParts(ctx, &inventoryV1.ListPartsRequest{
				Filter: &inventoryV1.PartsFilter{},
			})

			Expect(err).ToNot(HaveOccurred())
			parts := resp.GetParts()
			Expect(parts).ToNot(BeNil())
			Expect(parts).To(HaveLen(3))

			AssertPartsEqual(insertedPartPtr1, parts[0])
			AssertPartsEqual(insertedPartPtr2, parts[1])
			AssertPartsEqual(insertedPartPtr3, parts[2])
		})

		It("list parts by uuids", func() {
			resp, err := inventoryClient.ListParts(ctx, &inventoryV1.ListPartsRequest{
				Filter: &inventoryV1.PartsFilter{
					Uuids: []string{insertedPartPtr1.Uuid, insertedPartPtr2.Uuid},
				},
			})

			Expect(err).ToNot(HaveOccurred())

			parts := resp.GetParts()
			Expect(parts).To(HaveLen(2))

			receivedUuids := []string{parts[0].Uuid, parts[1].Uuid}
			Expect(receivedUuids).To(ConsistOf(
				insertedPartPtr1.Uuid,
				insertedPartPtr2.Uuid,
			))
		})

		It("list parts by names", func() {
			resp, err := inventoryClient.ListParts(ctx, &inventoryV1.ListPartsRequest{
				Filter: &inventoryV1.PartsFilter{
					Names: []string{insertedPartPtr1.Name, insertedPartPtr3.Name},
				},
			})

			Expect(err).ToNot(HaveOccurred())

			parts := resp.GetParts()
			Expect(parts).To(HaveLen(2))

			receivedNames := []string{parts[0].Name, parts[1].Name}
			Expect(receivedNames).To(ConsistOf(
				insertedPartPtr1.Name,
				insertedPartPtr3.Name,
			))
		})

		It("list parts by countries", func() {
			resp, err := inventoryClient.ListParts(ctx, &inventoryV1.ListPartsRequest{
				Filter: &inventoryV1.PartsFilter{
					ManufacturerCountries: []string{
						insertedPartPtr2.Manufacturer.Country,
						insertedPartPtr3.Manufacturer.Country,
					},
				},
			})

			Expect(err).ToNot(HaveOccurred())

			parts := resp.GetParts()
			Expect(parts).To(HaveLen(2))

			receivedCountries := []string{parts[0].Manufacturer.Country, parts[1].Manufacturer.Country}
			Expect(receivedCountries).To(ConsistOf(
				insertedPartPtr2.Manufacturer.Country,
				insertedPartPtr3.Manufacturer.Country,
			))
		})
	})
})
