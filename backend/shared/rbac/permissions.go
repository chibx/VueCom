package rbac

const (
	CatUser       = "user_management"
	CatRole       = "role_management"
	CatSales      = "sales"
	CatInventory  = "inventory"
	CatOrder      = "orders"
	CatMetrics    = "metrics"
	CatDeployment = "deployment"
)

// User Management
const (
	PermUserView       Permission = "user:view"
	PermUserCreate     Permission = "user:create"
	PermUserEdit       Permission = "user:edit"
	PermUserDelete     Permission = "user:delete"
	PermUserDeleteAny  Permission = "user:delete:any" // Golden
	PermSignupTokenDel Permission = "user:sgn-tk-del"
	PermSignupTokenCrt Permission = "user:sgn-tk-crt"

// PermUserImpersonate Permission = "user:impersonate"
)

// Role Management
const (
	PermRoleView   Permission = "role:view"
	PermRoleCreate Permission = "role:create"
	PermRoleEdit   Permission = "role:edit"
	PermRoleDelete Permission = "role:delete"
)

// Sales
const (
	PermSalesView   Permission = "sales:view"
	PermSalesCreate Permission = "sales:create"
	PermSalesEdit   Permission = "sales:edit"
	PermSalesDelete Permission = "sales:delete"
	PermSalesExport Permission = "sales:export"
)

// Inventory
const (
	PermInventoryView   Permission = "inventory:view"
	PermInventoryCreate Permission = "inventory:create"
	PermInventoryEdit   Permission = "inventory:edit"
	PermInventoryDelete Permission = "inventory:delete"
	PermInventoryManage Permission = "inventory:manage"
	// All CRUD
)

// Orders
const (
	PermOrderView    Permission = "order:view"
	PermOrderProcess Permission = "order:process"
	PermOrderCancel  Permission = "order:cancel"
	PermOrderRefund  Permission = "order:refund"
	PermOrderManage  Permission = "order:manage"
)

// System
const (
	PermMetricsView   Permission = "metrics:view"
	PermDeployView    Permission = "deploy:view"
	PermDeployExecute Permission = "deploy:execute"
	PermSystemConfig  Permission = "system:config"
)

// AllPermissions registry for validation
var AllPermissions = []Permission{
	PermUserView, PermUserCreate, PermUserEdit, PermUserDelete, PermUserDeleteAny,
	PermRoleView, PermRoleCreate, PermRoleEdit, PermRoleDelete,
	PermSalesView, PermSalesCreate, PermSalesEdit, PermSalesDelete, PermSalesExport,
	PermInventoryView, PermInventoryCreate, PermInventoryEdit, PermInventoryDelete, PermInventoryManage,
	PermOrderView, PermOrderProcess, PermOrderCancel, PermOrderRefund, PermOrderManage,
	PermMetricsView, PermDeployView, PermDeployExecute, PermSystemConfig,
}

var DefaultPermissions = []Permission{
	// Will add permissions that would get added by default or be toggled on
}
