package grpc

import (
	catalog_service "github.com/chibx/vuecom/backend/services/catalog"
	inventory_service "github.com/chibx/vuecom/backend/services/inventory"
	order_service "github.com/chibx/vuecom/backend/services/orders"
	payment_service "github.com/chibx/vuecom/backend/services/payment"
	"github.com/chibx/vuecom/backend/shared/proto/go/catalog"
	"github.com/chibx/vuecom/backend/shared/proto/go/orders"
)

var (
	OrderClient   orders.OrderServiceClient
	CatalogClient catalog.CatalogServiceClient
)

func InitClients() func() {

	// Register ALL services on the single in-memory server
	registerServices(
		catalog_service.Register,
		order_service.Register,
		// payment.Register,
		// notification.Register,
		// ...
	)

	conn := clientConn()

	OrderClient = orders.NewOrderServiceClient(conn)
	CatalogClient = catalog.NewCatalogServiceClient(conn)

	return func() {
		catalog_service.Destroy()
		order_service.Destroy()
		payment_service.Destroy()
		inventory_service.Destroy()
		conn.Close()
		shutdown()
	}
}
