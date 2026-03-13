package grpc

import (
	catalog_service "github.com/chibx/vuecom/backend/services/catalog"
	order_service "github.com/chibx/vuecom/backend/services/orders"
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
		conn.Close()
		shutdown()
	}
}
