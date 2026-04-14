package products

import (
	catReq "github.com/chibx/vuecom/backend/services/gateway/api/v1/request/catalog"
)

func normalizeProdReq(prod *catReq.CreateProductReq) {
	if prod.SalePrice > prod.BasePrice {
		prod.SalePrice = prod.BasePrice
	}

	if prod.DiscountEnd != nil && prod.DiscountStart != nil {
		if prod.DiscountEnd.Before(*prod.DiscountStart) {
			prod.DiscountEnd = prod.DiscountStart
		}
	}

	if prod.NewTo != nil && prod.NewFrom != nil {
		if prod.NewTo.Before(*prod.NewFrom) {
			prod.NewTo = prod.NewFrom
		}
	}

}
