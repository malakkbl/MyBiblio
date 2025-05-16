package constants

// User roles
const (
	RoleAdmin    = "admin"
	RoleUser     = "user"
	RoleManager  = "manager"
	RoleEmployee = "employee"
)

// ValidRoles contains all valid role values
var ValidRoles = map[string]bool{
	RoleAdmin:    true,
	RoleUser:     true,
	RoleManager:  true,
	RoleEmployee: true,
}
