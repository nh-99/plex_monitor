package discord

// Stores functions that are used as commands for the Discord bot.

import (
	"plex_monitor/internal/database/models"
	"plex_monitor/internal/worker"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

// HandlerFunc is a struct that represents a command handler. It is used to
// register commands with the Discord bot. The Name field is used to match
// the command with the command handler.
type HandlerFunc struct {
	Name    string                                                     `json:"name" bson:"name"`
	Handler func(s *discordgo.Session, i *discordgo.InteractionCreate) `json:"-" bson:"-"`
}

// GetCommands returns a slice of commands that are used to register commands
func GetCommands() []*discordgo.ApplicationCommand {
	return []*discordgo.ApplicationCommand{
		{
			Name:        "health",
			Description: "Get the health of the services",
		},
		{
			Name:        "rescan-plex-library",
			Description: "Rescan a Plex library",
			Type:        discordgo.ChatApplicationCommand,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:         "library",
					Description:  "The library to rescan",
					Type:         discordgo.ApplicationCommandOptionInteger,
					Required:     true,
					Autocomplete: true,
				},
			},
		},
		{
			Name:        "request",
			Description: "🙏 add any content you want!",
			Type:        discordgo.ChatApplicationCommand,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "title",
					Description: "The title of the content you want",
					Type:        discordgo.ApplicationCommandOptionString,
					Required:    true,
				},
			},
		},
	}
}

// GetCommandHandlers returns a map of command handlers that are used to
// register commands
func GetCommandHandlers() map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"health":              healthHandler,
		"rescan-plex-library": refreshPlexLibraryHandler,
		"request":             requestMedia,
	}
}

// userHasAccessToCommand checks if the user has access to the command.
func userHasAccessToCommand(permissionType models.PermissionType, i *discordgo.InteractionCreate, l *logrus.Entry) (bool, *models.User) {
	if i == nil {
		// There is no interaction to check permissions for, so no-op
		return false, nil
	}

	var userID string
	if i.User != nil && i.Member == nil {
		// If we have a user but no member, then we are dealing with a DM
		userID = i.User.ID
	} else if i.User == nil && i.Member != nil && i.Member.User != nil {
		// If we have a member and a member user, then we are dealing with a guild
		userID = i.Member.User.ID
	} else {
		// There is no user to check permissions for
		return false, nil
	}

	user, err := models.GetUserWithFrontendUserID(userID)
	if err != nil {
		l.WithFields(logrus.Fields{
			"error":         err,
			"discordUserId": i.User.ID,
		}).Error("Failed to get user with frontend user ID")
		return false, nil
	}

	if user.IsAnonymous() {
		l.WithFields(logrus.Fields{
			"discordUserId": i.User.ID,
		}).Info("User is anonymous")
		return false, nil
	}

	l = l.WithFields(logrus.Fields{
		"discordUserId":  i.User.ID,
		"userId":         user.ID.Hex(),
		"commandName":    i.ApplicationCommandData().Name,
		"commandId":      i.ApplicationCommandData().ID,
		"permissionType": permissionType,
	})

	l.Infof("User %s is performing command %s", user.Email, i.ApplicationCommandData().Name)

	return user.CheckPermission(permissionType), user
}

// handleAccessOrError checks if the user has access to the command and if not
// responds with an error message.
func handleAccessOrError(permissionType models.PermissionType, s *discordgo.Session, i *discordgo.InteractionCreate, l *logrus.Entry) (bool, *models.User) {
	hasAccess, user := userHasAccessToCommand(permissionType, i, l)
	if !hasAccess {
		l.WithFields(logrus.Fields{
			"discordUserId":  i.User.ID,
			"commandName":    i.ApplicationCommandData().Name,
			"commandId":      i.ApplicationCommandData().ID,
			"permissionType": permissionType,
		}).Infof("User %s does not have permission to perform command %s", i.User.ID, i.ApplicationCommandData().Name)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "You don't have permission to perform this command",
			},
		})
		return false, nil
	}
	return true, user
}

func respondToError(s *discordgo.Session, i *discordgo.InteractionCreate, err error, l *logrus.Entry) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "An error occurred while executing this command",
		},
	})
	l.WithFields(logrus.Fields{
		"error": err,
	}).Error("An error occurred while executing this command")
}

func healthHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	l := logrus.NewEntry(logrus.StandardLogger())

	hasAccess, _ := handleAccessOrError(models.PermissionTypeCheckHealth, s, i, l)
	if !hasAccess {
		return
	}

	var embeds []*discordgo.MessageEmbed

	// Loop through all services and get their health
	latestHealth := worker.GetCronWorker(worker.NewHealthCronWorker().Name()).(*worker.HealthCronWorker).GetLatestHealthMap()
	for serviceName, serviceHealth := range latestHealth {
		// Get the service
		service, err := models.GetServiceByName(models.ServiceType(serviceName))
		if err != nil {
			l.WithFields(logrus.Fields{
				"error":       err,
				"serviceName": serviceName,
			}).Error("Failed to get service by name")
			continue
		}

		// Get the config as a standard config
		config, err := service.GetConfigAsStandardConfig()
		if err != nil {
			l.WithFields(logrus.Fields{
				"error":       err,
				"serviceName": serviceName,
			}).Error("Failed to get config as standard config")
			continue
		}

		healthColor := 0x00ff00
		if !serviceHealth.Healthy {
			healthColor = 0xff0000
		}

		status := "Healthy"
		if !serviceHealth.Healthy {
			status = "Unhealthy"
		}

		// Create the embed
		embed := &discordgo.MessageEmbed{
			URL:   config.Host,
			Title: strings.Title(strings.ToLower(string(service.ServiceName))) + " Health",
			Color: healthColor,
			Thumbnail: &discordgo.MessageEmbedThumbnail{
				URL: getImageForService(string(service.ServiceName)),
			},
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "Status",
					Value:  status,
					Inline: true,
				},
				{
					Name:   "Version",
					Value:  serviceHealth.Version,
					Inline: true,
				},
				{
					Name:   "Last Checked",
					Value:  serviceHealth.LastChecked.Format("2006-01-02 15:04:05"),
					Inline: true,
				},
			},
		}

		embeds = append(embeds, embed)
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: embeds,
			Flags:  discordgo.MessageFlagsEphemeral,
		},
	})
}

func getImageForService(serviceName string) string {
	switch serviceName {
	case "plex":
		return "https://www.plex.tv/wp-content/themes/plex/assets/img/favicons/plex-180.png"
	case "ombi":
		return "https://raw.githubusercontent.com/Ombi-app/Ombi/gh-pages/img/android-chrome-512x512.png"
	case "radarr":
		return "https://wiki.servarr.com/assets/radarr/logos/512.png"
	case "sonarr":
		return "https://wiki.servarr.com/assets/sonarr/logos/512.png"
	}

	return ""
}
