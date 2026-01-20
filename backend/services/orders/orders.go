package order_service

type catalogService struct{}

func NewCatalogService() *catalogService {
	return &catalogService{}
}
