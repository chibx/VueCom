package types

import (
	"errors"

	"github.com/chibx/vuecom/backend/services/gateway/internal/types/database"
)

type DatabaseErr error

var ErrDbNil DatabaseErr = errors.New("record not found")
var ErrDbUnimplemented DatabaseErr = errors.New("unimplemented")

type Database interface {
	BackendUsers() database.BackendUserRepository
	Customers() database.CustomerRepository
	Categories() database.CategoryRepository
	Products() database.ProductRepository
	Orders() database.OrderRepository
	Inventory() database.InventoryRepository
	AppData() database.AppDataRepository
	Migrate() error
}
