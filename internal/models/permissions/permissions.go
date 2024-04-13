package permissions

// ============================================================================
// Constants
// ============================================================================

const (
	PermissionAdmin      = "admin"
	PermissionSuperAdmin = "superadmin"
)

// ============================================================================
// Permission Type
// ============================================================================

// Type to hold user permission codes
type Perms []string

// Checks if a user has a specific permission
func (p Perms) Include(code string) bool {
	for i := range p {
		if code == p[i] {
			return true
		}
	}
	return false
}
