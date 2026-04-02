package utils

import (
	"github.com/chibx/vuecom/backend/services/catalog/internal/global"
	"github.com/chibx/vuecom/backend/shared/models/db/catalog"
	catalogPr "github.com/chibx/vuecom/backend/shared/proto/go/catalog"
	"github.com/chibx/vuecom/backend/shared/utils"
	"go.uber.org/zap"
)

func FailOnError(err error, msg string) {
	if err != nil {
		global.Logger.Fatal(msg, zap.Error(err))
	}
}

func CreateProdRpcToDBModel(req *catalogPr.CreateProductRequest) *catalog.Product {
	return &catalog.Product{
		Name:             req.Name,
		SKU:              req.Sku,
		BasePrice:        req.BasePrice,
		SalePrice:        req.SalePrice,
		DiscountStart:    utils.AsPointer(req.DiscountStart.AsTime()),
		DiscountEnd:      utils.AsPointer(req.DiscountEnd.AsTime()),
		IsNew:            req.IsNew,
		NewFrom:          utils.AsPointer(req.NewFrom.AsTime()),
		NewTo:            utils.AsPointer(req.NewTo.AsTime()),
		CountryOfManf:    req.CountryOfManf,
		Enabled:          req.Enabled,
		ShortDescription: req.ShortDescription,
		FullDescription:  req.FullDescription,
		Slug:             req.Slug,
		Weight:           req.Weight,
		BrandId:          req.BrandId,
		ColorId:          req.ColorId,
		MetaTitle:        &req.MetaTitle,
		MetaDescription:  &req.MetaDescription,
		SearchKeywords:   req.SearchKeywords,
		PresetID:         req.PresetId,
	}
}
