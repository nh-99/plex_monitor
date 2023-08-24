package web

import "plex_monitor/internal/database/models"

// AppData is the data passed to the template.
type AppData struct {
	// Name is the name of the application
	Name string
	// User is the user that is logged in
	User models.User
}

// GetAppData returns the AppData struct.
func GetAppData(user *models.User) AppData {
	if user == nil {
		return AppData{
			Name: "Plex Monitor",
		}
	}

	return AppData{
		Name: "Plex Monitor",
		User: *user,
	}
}
