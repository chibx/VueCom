// Shared types for propagated context data
package ctx

import "github.com/chibx/vuecom/backend/shared/models/db/users"

type BackendUser struct {
	ID     int
	ApiKey *users.ApiKey // To know if the user is acting through an api key
}

type Customer struct {
	ID int
}
