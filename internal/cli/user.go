package cli

import (
	"fmt"
	"plex_monitor/internal/database"
	"plex_monitor/internal/database/models"
	"plex_monitor/internal/utils"
	"time"

	"github.com/urfave/cli/v2"
)

func getUserCreateCmd() *cli.Command {
	return &cli.Command{
		Name:    "user",
		Aliases: []string{"u"},
		Usage:   "Configure a new user in the system",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "email"},
		},
		Action: func(cCtx *cli.Context) error {
			email := cCtx.String("email")
			password := getPassword("Enter a password: ")
			hashBytes, _ := utils.HashString(password)
			s := string(hashBytes)
			_, err := database.DB.Collection("users").InsertOne(database.Ctx, models.User{
				Email:          email,
				HashedPassword: s,
				Activated:      true,
				CreatedAt:      time.Now(),
				CreatedBy:      models.SystemUserID,
				UpdatedAt:      time.Now(),
				UpdatedBy:      models.SystemUserID,
			})
			if err != nil {
				panic(err)
			}
			return nil
		},
	}
}

func getUserUpdatePermissionsCmd() *cli.Command {
	return &cli.Command{
		Name:    "permission",
		Aliases: []string{"p"},
		Usage:   "Update a user's permissions in the system",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "id"},
			&cli.StringFlag{Name: "email"},
			&cli.StringFlag{Name: "permission"},
			&cli.BoolFlag{Name: "remove"},
		},
		Action: func(cCtx *cli.Context) error {
			email := cCtx.String("email")
			id := cCtx.String("id")
			permission := cCtx.String("permission")
			remove := cCtx.Bool("remove")

			if email == "" && id == "" {
				fmt.Println("Must supply either a user ID or email")
				return nil
			}

			if permission == "" {
				// Show the valid permissions to the user
				fmt.Println("---------------------------------------------")
				fmt.Println("|             Valid Permissions             |")
				fmt.Println("---------------------------------------------")
				readablePerms := models.GetReadableUserPermissions()
				for _, p := range readablePerms {
					fmt.Printf("%s - %s, %s\n", p.PermissionType, p.Name, p.Description)
				}

				return nil
			}

			// Check if the permission is valid
			valid := false
			for _, p := range models.GetReadableUserPermissions() {
				if p.PermissionType == models.PermissionType(permission) {
					valid = true
					break
				}
			}

			if !valid {
				fmt.Printf("Invalid permission: %s\n", permission)
				return nil
			}

			// Lookup the user
			user, err := models.GetUser(id, email)
			if err != nil {
				panic(err)
			}

			if remove {
				// Find the permission in the user's permissions
				index := -1
				for i, p := range user.Permissions {
					if p == models.PermissionType(permission) {
						index = i
						break
					}
				}

				if index == -1 {
					fmt.Printf("User %s does not have permission %s\n", user.Email, permission)
					return nil
				}

				// Remove the permission from the user
				user.Permissions = append(user.Permissions[:index], user.Permissions[index+1:]...)

				// Update the user in the database
				err = user.Save()
				if err != nil {
					panic(err)
				}

				fmt.Printf("Removed permission %s from user %s\n", permission, user.Email)
			} else {
				// Check if the user already has the permission
				hasPermission := user.CheckPermission(models.PermissionType(permission))
				if hasPermission {
					fmt.Printf("User %s already has permission %s\n", user.Email, permission)
					return nil
				}

				// Add the permission to the user
				user.Permissions = append(user.Permissions, models.PermissionType(permission))

				// Update the user in the database
				err = user.Save()
				if err != nil {
					panic(err)
				}

				fmt.Printf("Added permission %s to user %s\n", permission, user.Email)
			}

			return nil
		},
	}
}

func getUserFrontendUpdateCmd() *cli.Command {
	return &cli.Command{
		Name:    "frontend",
		Aliases: []string{"f"},
		Usage:   "Update the frontend service for the supplied user",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "id"},
			&cli.StringFlag{Name: "email"},
			&cli.StringFlag{Name: "frontendType", Required: true},
			&cli.StringFlag{Name: "frontendUserId"},
			&cli.BoolFlag{Name: "remove"},
		},
		Action: func(cCtx *cli.Context) error {
			email := cCtx.String("email")
			id := cCtx.String("id")
			frontendType := cCtx.String("frontendType")
			frontendUserId := cCtx.String("frontendUserId")
			remove := cCtx.Bool("remove")

			if email == "" && id == "" {
				fmt.Println("Must supply either a user ID or email")
				return nil
			}

			if frontendType == "" {
				fmt.Println("Must supply a frontend type")
				return nil
			}

			// Lookup the user
			user, err := models.GetUser(id, email)
			if err != nil {
				panic(err)
			}

			if remove {
				// Find the frontend service in the user's frontend services
				index := -1
				for i, s := range user.FrontendServices {
					if s.Type == models.FrontendServiceType(frontendType) {
						index = i
						break
					}
				}

				if index == -1 {
					fmt.Printf("User %s does not have frontend service %s\n", user.Email, frontendType)
					return nil
				}

				// Remove the frontend service from the user
				user.FrontendServices = append(user.FrontendServices[:index], user.FrontendServices[index+1:]...)

				// Update the user in the database
				err = user.Save()
				if err != nil {
					panic(err)
				}

				fmt.Printf("Removed frontend service %s from user %s\n", frontendUserId, user.Email)
			} else {
				// Check if the user already has the frontend service
				hasFrontendService := false
				for _, s := range user.FrontendServices {
					if s.Type == models.FrontendServiceType(frontendType) && (s.UserID == "" || s.UserID == frontendUserId) {
						hasFrontendService = true
						break
					}
				}

				if hasFrontendService {
					fmt.Printf("User %s already has frontend service %s\n", user.Email, frontendUserId)
				}

				// Add the frontend service to the user
				user.FrontendServices = append(user.FrontendServices, models.FrontendService{
					Type:   models.FrontendServiceType(frontendType),
					UserID: frontendUserId,
				})

				// Update the user in the database
				err = user.Save()
				if err != nil {
					panic(err)
				}

				fmt.Printf("Added frontend service %s to user %s\n", frontendType, user.Email)
			}

			return nil
		},
	}
}
