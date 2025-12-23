package types

import (
	"context"
	"errors"
	dbModels "vuecom/shared/models/db"
)

var ErrDbNil = errors.New("record not found")
var ErrDbUnimplemented = errors.New("unimplemented")

type Database interface {
	BackendUsers() BackendUserRepository
	Customers() CustomerRepository
	Categories() CategoryRepository
	Products() ProductRepository
	Orders() OrderRepository
	Inventory() InventoryRepository
	AppData() AppDataRepository
	Migrate() error
}

type BackendUserRepository interface {
	GetBackendUserById(id int, ctx context.Context) (*dbModels.BackendUser, error)
	// The full api key is passed to the function
	//
	// Splitting should be done by the function
	GetBackendUserByApiKey(apiKey string, ctx context.Context) (*dbModels.BackendUser, error)
}

type CustomerRepository interface {
	GetCustomerById(id int, ctx context.Context) (*dbModels.Customer, error)
}

type CategoryRepository interface {
	GetCategoryById(id int, ctx context.Context) (*dbModels.Category, error)
}

type ProductRepository interface {
	GetProductById(id int, ctx context.Context) (*dbModels.Product, error)
}

type OrderRepository interface {
	GetOrderById(id int, ctx context.Context) (*dbModels.Order, error)
}

type InventoryRepository interface {
	GetInventoryById(id int, ctx context.Context) (*dbModels.Inventory, error)
}

type AppDataRepository interface {
	GetAppData(ctx context.Context) (*dbModels.AppData, error)
}
