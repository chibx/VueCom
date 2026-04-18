package igrpc

import (
	catalog_service "github.com/chibx/vuecom/backend/services/catalog"
	gl "github.com/chibx/vuecom/backend/services/gateway/internal/global"
	inventory_service "github.com/chibx/vuecom/backend/services/inventory"
	order_service "github.com/chibx/vuecom/backend/services/orders"
	payment_service "github.com/chibx/vuecom/backend/services/payment"
	"github.com/chibx/vuecom/backend/shared/proto/go/catalog"
	"github.com/chibx/vuecom/backend/shared/proto/go/inventory"
	"github.com/chibx/vuecom/backend/shared/proto/go/orders"
	"go.uber.org/zap"
)

var (
	OrderClient     orders.OrderServiceClient
	CatalogClient   catalog.CatalogServiceClient
	InventoryClient inventory.InventoryServiceClient
)

func InitClients() func() {

	// Register ALL services on the single in-memory server
	registerServices(
		catalog_service.Register,
		inventory_service.Register,
		order_service.Register,
		// payment.Register,
		// notification.Register,
		// ...
	)

	conn := clientConn()

	OrderClient = orders.NewOrderServiceClient(conn)
	CatalogClient = catalog.NewCatalogServiceClient(conn)
	InventoryClient = inventory.NewInventoryServiceClient(conn)

	return func() {
		catalog_service.Destroy()
		order_service.Destroy()
		payment_service.Destroy()
		inventory_service.Destroy()
		err := conn.Close()
		if err != nil {
			gl.Logger.Error("Error closing gRPC listener", zap.Error(err))
		}
		shutdown()
	}
}
