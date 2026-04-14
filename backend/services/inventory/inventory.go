package inventory_service

import (
	"github.com/chibx/vuecom/backend/services/inventory/internal/global"
	"github.com/chibx/vuecom/backend/services/inventory/internal/pubsub"
	inventoryPr "github.com/chibx/vuecom/backend/shared/proto/go/inventory"
	"google.golang.org/grpc"
)

func Register(s *grpc.Server) {
	pubsub.InitPubSub()
	inventoryPr.RegisterInventoryServiceServer(s, &Service{})
}

func Destroy() {
	pubsub.DefPubSub.Close()
	global.Redis.Close()
}
