package products

import (
	"github.com/chibx/vuecom/backend/shared/proto/go/catalog"
	sharedReq "github.com/chibx/vuecom/backend/shared/types/request"
	"github.com/chibx/vuecom/backend/shared/utils"
)

func normalizeProdReq(prod *sharedReq.CreateProductReq) {
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

func createProductToRpc(s *sharedReq.CreateProductReq, parentId ...*uint32) (*catalog.CreateProductRequest, error) {
	// presetVal, err := structpb.NewStruct(s.PresetValues)
	// if err != nil {
	// 	return nil, err
	// }
	var p_id *uint32
	if len(parentId) > 0 {
		p_id = parentId[0]
	}

	return &catalog.CreateProductRequest{
		Name:             s.Name,
		Sku:              s.SKU,
		BasePrice:        s.BasePrice,
		SalePrice:        s.SalePrice,
		DiscountStart:    utils.NilTimeToRpc(s.DiscountStart),
		DiscountEnd:      utils.NilTimeToRpc(s.DiscountEnd),
		NewFrom:          utils.NilTimeToRpc(s.NewFrom),
		NewTo:            utils.NilTimeToRpc(s.NewTo),
		Weight:           s.Weight,
		Quantity:         s.Quantity,
		CountryOfManf:    s.CountryOfManf,
		Slug:             s.Slug,
		BrandId:          s.BrandId,
		ColorId:          s.ColorId,
		Medias:           s.Medias,
		MetaTitle:        s.MetaTitle,
		MetaDescription:  s.MetaDescription,
		SearchKeywords:   s.SearchKeywords,
		RelatedProducts:  s.RelatedProducts,
		UpSellProducts:   s.UpSellProducts,
		CrossSell:        s.CrossSell,
		PresetId:         s.PresetID,
		PresetValues:     s.PresetValues,
		InStock:          s.InStock,
		Categories:       s.Categories,
		IsNew:            s.IsNew,
		Enabled:          s.Enabled,
		ShortDescription: s.ShortDescription,
		FullDescription:  s.FullDescription,
		ParentId:         p_id,
	}, nil
}

func getProductFromRpc(s *catalog.GetProductResponse) (*sharedReq.CreateProductReq, error) {
	// var p_id *uint32
	// if len(parentId) > 0 {
	// 	p_id = parentId[0]
	// }

	return &sharedReq.CreateProductReq{
		Name:             s.Name,
		SKU:              s.Sku,
		BasePrice:        s.BasePrice,
		SalePrice:        s.SalePrice,
		DiscountStart:    utils.AsPointer(s.DiscountStart.AsTime()),
		DiscountEnd:      utils.AsPointer(s.DiscountEnd.AsTime()),
		NewFrom:          utils.AsPointer(s.NewFrom.AsTime()),
		NewTo:            utils.AsPointer(s.NewTo.AsTime()),
		Weight:           s.Weight,
		Quantity:         s.Quantity,
		CountryOfManf:    s.CountryOfManf,
		Slug:             s.Slug,
		BrandId:          s.BrandId,
		ColorId:          s.ColorId,
		MetaTitle:        s.MetaTitle,
		MetaDescription:  s.MetaDescription,
		SearchKeywords:   s.SearchKeywords,
		RelatedProducts:  s.RelatedProducts,
		UpSellProducts:   s.UpSellProducts,
		CrossSell:        s.CrossSell,
		PresetID:         s.PresetId,
		PresetValues:     s.PresetValues,
		InStock:          s.InStock,
		Categories:       s.Categories,
		IsNew:            s.IsNew,
		Enabled:          s.Enabled,
		ShortDescription: s.ShortDescription,
		FullDescription:  s.FullDescription,
	}, nil
}
