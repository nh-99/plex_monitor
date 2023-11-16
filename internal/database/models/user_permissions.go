package models

// PermissionType is the type of the permission.
type PermissionType string

// ReadableUserPermission is the struct that represents the human readable user permission.
// It can be used for pulling the available permissions in the system
type ReadableUserPermission struct {
	Name           string         `json:"name"`
	Description    string         `json:"description"`
	PermissionType PermissionType `json:"permission_type"`
}

// userPermissions is the map of the user permissions.
var userPermissions map[string]ReadableUserPermission

const (
	// PermissionTypeGodMode is the application god mode.
	PermissionTypeGodMode PermissionType = "pm.GOD"
	// PermissionTypeCheckHealth is the type of the check health permission.
	PermissionTypeCheckHealth PermissionType = "pm.check_health"
	// PermissionTypeManageUsers is the type of the manage users permission.
	PermissionTypeManageUsers PermissionType = "pm.manage_users"
	// PermissionTypeScanLibrary is the type of the scan library permission.
	PermissionTypeScanLibrary PermissionType = "pm.scan_library"
	// PermissionTypeRequestMedia allows a user to request new media.
	PermissionTypeRequestMedia PermissionType = "pm.request_media"
)

func init() {
	RegisterReadableUserPermission(ReadableUserPermission{
		Name:           "Check Health",
		Description:    "Allows the user to check the health of the system",
		PermissionType: PermissionTypeCheckHealth,
	})
	RegisterReadableUserPermission(ReadableUserPermission{
		Name:           "Manage Users",
		Description:    "Allows the user to manage users in the system",
		PermissionType: PermissionTypeManageUsers,
	})
	RegisterReadableUserPermission(ReadableUserPermission{
		Name:           "Scan Library",
		Description:    "Allows the user to scan the Plex libraries",
		PermissionType: PermissionTypeScanLibrary,
	})
	RegisterReadableUserPermission(ReadableUserPermission{
		Name:           "Request Media",
		Description:    "Allows the user to request new media",
		PermissionType: PermissionTypeRequestMedia,
	})
}

// CheckPermission checks if the user has the supplied permission.
func (u User) CheckPermission(permission PermissionType) bool {
	for _, p := range u.Permissions {
		godMode := p == PermissionTypeGodMode // Override for "god mode"
		if p == permission || godMode {
			return true
		}
	}

	return false
}

// RegisterReadableUserPermission registers a new readable user permission.
func RegisterReadableUserPermission(permission ReadableUserPermission) {
	if userPermissions == nil {
		userPermissions = make(map[string]ReadableUserPermission)
	}
	userPermissions[string(permission.PermissionType)] = permission
}

// GetReadableUserPermissions returns the readable user permissions.
func GetReadableUserPermissions() map[string]ReadableUserPermission {
	return userPermissions
}
