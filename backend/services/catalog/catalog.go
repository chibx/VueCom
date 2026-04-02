package catalog_service

import (
	"github.com/chibx/vuecom/backend/services/catalog/internal/global"
	catalogPr "github.com/chibx/vuecom/backend/shared/proto/go/catalog"
	"google.golang.org/grpc"
)

type catalogService struct{}

func NewCatalogService() *catalogService {
	return &catalogService{}
}

func Register(s *grpc.Server) {
	catalogPr.RegisterCatalogServiceServer(s, &Service{})
}

func Destroy() {
	global.AmqpChan.Close()
	global.AmqpConn.Close()
	global.Redis.Close()
}
