package integration

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	inventoryV1 "github.com/Artyom099/factory/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/Artyom099/factory/shared/pkg/proto/payment/v1"
)

var _ = Describe("OrderService", Ordered, func() {
	var (
		ctx             context.Context
		cancel          context.CancelFunc
		inventoryClient inventoryV1.InventoryServiceClient
		_               paymentV1.PaymentServiceClient
	)

	BeforeEach(func() {
		ctx, cancel = context.WithCancel(suiteCtx)

		conn, err := grpc.NewClient(
			env.App.Address(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		Expect(err).ToNot(HaveOccurred(), "ожидали успешное подключение к gRPC приложению")

		inventoryClient = inventoryV1.NewInventoryServiceClient(conn)
	})

	AfterEach(func() {
		err := env.ClearDB(ctx)
		Expect(err).ToNot(HaveOccurred(), "ожидали успешную очистку бд")

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
})
