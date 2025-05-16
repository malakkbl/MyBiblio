package constants

import "strings"

// User roles
const (
	RoleAdmin    = "admin"
	RoleManager  = "manager"
	RoleEmployee = "employee"
	RoleUser     = "user"
)

// RoleDescription maps roles to their descriptions
var RoleDescription = map[string]string{
	RoleAdmin:    "Full system access and management capabilities",
	RoleManager:  "Manage books, authors, and generate reports",
	RoleEmployee: "Handle orders and customer service",
	RoleUser:     "Browse books and place orders",
}

// ValidRoles contains all valid role values
var ValidRoles = map[string]bool{
	RoleAdmin:    true,
	RoleManager:  true,
	RoleEmployee: true,
	RoleUser:     true,
}

// GetAvailableRoles returns a list of all available roles
func GetAvailableRoles() []string {
	roles := make([]string, 0, len(ValidRoles))
	for role := range ValidRoles {
		roles = append(roles, role)
	}
	return roles
}

// IsValidRole checks if a role is valid
func IsValidRole(role string) bool {
	return ValidRoles[strings.ToLower(role)]
}

// GetRoleDescription returns a user-friendly description of a role
func GetRoleDescription(role string) string {
	if desc, ok := RoleDescription[strings.ToLower(role)]; ok {
		return desc
	}
	return "Unknown role"
}
