package authz

// Permission codes stored in DB (permissions.permission_code) and JWT-side role checks.
const (
	PermOrdersViewAny         = "orders.view_any"
	PermOrdersManageStatus    = "orders.manage_status"
	PermConfigurationsViewAny = "configurations.view_any"
	PermConfigurationsManage  = "configurations.manage"
	PermAppointmentsViewAny   = "appointments.view_any"
	PermGarageViewAny         = "garage.view_any"
	PermDocumentsViewAny      = "documents.view_any"
	PermNewsManage            = "news.manage"
	PermAdminOrderStatuses    = "admin.order_statuses"
	PermAdminRolesView        = "admin.roles_view"
)

// AllPermissionCodes lists every defined permission (for admin role seed and tests).
var AllPermissionCodes = []string{
	PermOrdersViewAny,
	PermOrdersManageStatus,
	PermConfigurationsViewAny,
	PermConfigurationsManage,
	PermAppointmentsViewAny,
	PermGarageViewAny,
	PermDocumentsViewAny,
	PermNewsManage,
	PermAdminOrderStatuses,
	PermAdminRolesView,
}

// DefaultRolePermissions is used when the DB has no role_permissions rows (bootstrap / tests).
func DefaultRolePermissions() map[string][]string {
	admin := append([]string(nil), AllPermissionCodes...)

	manager := []string{
		PermOrdersViewAny,
		PermOrdersManageStatus,
		PermConfigurationsViewAny,
		PermConfigurationsManage,
		PermAppointmentsViewAny,
		PermGarageViewAny,
		PermDocumentsViewAny,
		PermNewsManage,
	}

	serviceAdvisor := []string{
		PermOrdersViewAny,
		PermOrdersManageStatus,
		PermConfigurationsViewAny,
		PermAppointmentsViewAny,
		PermGarageViewAny,
		PermDocumentsViewAny,
	}

	return map[string][]string{
		"customer":        {},
		"manager":         manager,
		"service_advisor": serviceAdvisor,
		"admin":           admin,
	}
}
