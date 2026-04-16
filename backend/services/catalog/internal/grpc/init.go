package igrpc

import (
	gl "github.com/chibx/vuecom/backend/services/catalog/internal/global"
	inventory_service "github.com/chibx/vuecom/backend/services/inventory"
	order_service "github.com/chibx/vuecom/backend/services/orders"
	"github.com/chibx/vuecom/backend/shared/proto/go/inventory"
	"github.com/chibx/vuecom/backend/shared/proto/go/orders"
	"go.uber.org/zap"
)

var (
	OrderClient     orders.OrderServiceClient
	InventoryClient inventory.InventoryServiceClient
)

func InitClients() func() {

	// Register ALL services on the single in-memory server
	registerServices(
		inventory_service.Register,
		order_service.Register,
		// payment.Register,
		// notification.Register,
		// ...
	)

	conn := clientConn()

	OrderClient = orders.NewOrderServiceClient(conn)
	InventoryClient = inventory.NewInventoryServiceClient(conn)

	return func() {
		// payment_service.Destroy()
		order_service.Destroy()
		inventory_service.Destroy()
		err := conn.Close()
		if err != nil {
			gl.Logger.Error("Error closing gRPC listener", zap.Error(err))
		}
		shutdown()
	}
}
