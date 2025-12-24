package types

import (
	"context"
	"errors"
	dbModels "vuecom/shared/models/db"
)

type DatabaseErr error

var ErrDbNil DatabaseErr = errors.New("record not found")
var ErrDbUnimplemented DatabaseErr = errors.New("unimplemented")

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
	CreateUser(user *dbModels.BackendUser, ctx context.Context) error
	GetUserById(id int, ctx context.Context) (*dbModels.BackendUser, error)
	GetAdmin(ctx context.Context) (*dbModels.BackendUser, error)
	// The full api key is passed to the function
	//
	// Splitting should be done by the function
	GetUserByApiKey(apiKey string, ctx context.Context) (*dbModels.BackendUser, error)
	CreateSession(session *dbModels.BackendSession, ctx context.Context) error
	GetSessionByToken(token string, ctx context.Context) (*dbModels.BackendSession, error)
	GetSessions(userId int, ctx context.Context) ([]dbModels.BackendSession, error)
	DeleteSession(session *dbModels.BackendSession, ctx context.Context) error

	GetCountryIdByCode(code string, ctx context.Context) (uint, error)
}

type CustomerRepository interface {
	CreateUser(user *dbModels.Customer, ctx context.Context) error
	GetUserById(id int, ctx context.Context) (*dbModels.Customer, error)
	CreateSession(session *dbModels.CustomerSession, ctx context.Context) error
	GetSessionByToken(token string, ctx context.Context) (*dbModels.CustomerSession, error)
	GetSessions(userId int, ctx context.Context) ([]dbModels.CustomerSession, error)
	DeleteSession(session *dbModels.CustomerSession, ctx context.Context) error
}

type CategoryRepository interface {
	GetCategoryById(id int, ctx context.Context) (*dbModels.Category, error)
}

type ProductRepository interface {
	CreateProduct(product *dbModels.Product, ctx context.Context) error
	GetProductById(id int, ctx context.Context) (*dbModels.Product, error)
}

type OrderRepository interface {
	CreateOrder(order *dbModels.Order, ctx context.Context) error
	GetOrderById(id int, ctx context.Context) (*dbModels.Order, error)
}

type InventoryRepository interface {
	GetInventoryById(id int, ctx context.Context) (*dbModels.Inventory, error)
}

type AppDataRepository interface {
	CreateAppData(appData *dbModels.AppData, ctx context.Context) error
	GetAppData(ctx context.Context) (*dbModels.AppData, error)
	// api.Deps.DB.Model(&dbModels.BackendUser{}).Where("role = 'owner'").Count(&count)
	CountOwner(ctx context.Context) (int64, error)
}
