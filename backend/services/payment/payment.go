package payment_service

import (
	"github.com/chibx/vuecom/backend/services/payment/internal/global"
	"github.com/chibx/vuecom/backend/services/payment/internal/pubsub"
	ordersPr "github.com/chibx/vuecom/backend/shared/proto/go/orders"
	"google.golang.org/grpc"
)

func Register(s *grpc.Server) {
	pubsub.InitPubSub()
	ordersPr.RegisterOrderServiceServer(s, &Service{})
}

func Destroy() {
	pubsub.DefPubSub.Close()
	global.Redis.Close()
}
