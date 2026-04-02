package catalog_service

import (
	"github.com/chibx/vuecom/backend/services/catalog/internal/global"
	"github.com/chibx/vuecom/backend/services/catalog/internal/pubsub"
	catalogPr "github.com/chibx/vuecom/backend/shared/proto/go/catalog"
	"google.golang.org/grpc"
)

type catalogService struct{}

func NewCatalogService() *catalogService {
	return &catalogService{}
}

func Register(s *grpc.Server) {
	pubsub.InitPubSub()
	catalogPr.RegisterCatalogServiceServer(s, &Service{})
}

func Destroy() {
	pubsub.DefPubSub.Close()
	global.Redis.Close()
}
