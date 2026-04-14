package utils

import (
	catReq "github.com/chibx/vuecom/backend/services/gateway/api/v1/request/catalog"
	catRes "github.com/chibx/vuecom/backend/services/gateway/api/v1/response/catalog"
	"github.com/chibx/vuecom/backend/shared/proto/go/catalog"
	"github.com/chibx/vuecom/backend/shared/utils"
)

func CreateProductToRpc(s *catReq.CreateProductReq, parentId ...*uint32) (*catalog.CreateProductRequest, error) {
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

func CreateProdToGetResp(s *catReq.CreateProductReq, productId uint32) *catRes.GetProductResp {
	return &catRes.GetProductResp{
		ID:               productId,
		Name:             s.Name,
		SKU:              s.SKU,
		BasePrice:        s.BasePrice,
		SalePrice:        s.SalePrice,
		DiscountStart:    s.DiscountStart,
		DiscountEnd:      s.DiscountEnd,
		NewFrom:          s.NewFrom,
		NewTo:            s.NewTo,
		Weight:           s.Weight,
		Quantity:         s.Quantity,
		CountryOfManf:    s.CountryOfManf,
		Slug:             s.Slug,
		Medias:           s.Medias,
		BrandId:          s.BrandId,
		ColorId:          s.ColorId,
		MetaTitle:        s.MetaTitle,
		MetaDescription:  s.MetaDescription,
		SearchKeywords:   s.SearchKeywords,
		RelatedProducts:  s.RelatedProducts,
		UpSellProducts:   s.UpSellProducts,
		CrossSell:        s.CrossSell,
		PresetID:         s.PresetID,
		PresetValues:     s.PresetValues,
		InStock:          s.InStock,
		Categories:       s.Categories,
		IsNew:            s.IsNew,
		Enabled:          s.Enabled,
		ShortDescription: s.ShortDescription,
		FullDescription:  s.FullDescription,
	}
}

func GetProductFromRpc(s *catalog.GetProductResponse) (*catRes.GetProductResp, error) {
	// var p_id *uint32
	// if len(parentId) > 0 {
	// 	p_id = parentId[0]
	// }

	return &catRes.GetProductResp{
		ID:               s.Id,
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
		Medias:           s.Medias,
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
