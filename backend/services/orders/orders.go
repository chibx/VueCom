package order_service

import (
	"github.com/chibx/vuecom/backend/services/orders/internal/global"
	"github.com/chibx/vuecom/backend/services/orders/internal/pubsub"
	ordersPr "github.com/chibx/vuecom/backend/shared/proto/go/orders"
	"google.golang.org/grpc"
)

type orderService struct{}

func NewOrderService() *orderService {
	return &orderService{}
}

func Register(s *grpc.Server) {
	pubsub.InitPubSub()
	ordersPr.RegisterOrderServiceServer(s, &Service{})
}

func Destroy() {
	global.AmqpChan.Close()
	global.AmqpConn.Close()
	global.Redis.Close()
}
