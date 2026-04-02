// Shared types for propagated context data
package ctx

type BackendUser struct {
	ID int
	// To know if the user is acting through an api key
	UsedApi bool
}

type Customer struct {
	ID int
}
