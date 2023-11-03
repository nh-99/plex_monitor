package discord

// Stores functions that are used as commands for the Discord bot.

import (
	"plex_monitor/internal/database/models"

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
	}
}

// GetCommandHandlers returns a map of command handlers that are used to
// register commands
func GetCommandHandlers() map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"health":              healthHandler,
		"rescan-plex-library": refreshPlexLibraryHandler,
	}
}

// userHasAccessToCommand checks if the user has access to the command.
func userHasAccessToCommand(permissionType models.PermissionType, i *discordgo.InteractionCreate) bool {
	if i == nil {
		// There is no interaction to check permissions for, so no-op
		return false
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
		return false
	}

	user, err := models.GetUserWithFrontendUserID(userID)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":         err,
			"discordUserId": i.User.ID,
		}).Error("Failed to get user with frontend user ID")
		return false
	}

	if user.IsAnonymous() {
		return false
	}

	return user.CheckPermission(permissionType)
}

// handleAccessOrError checks if the user has access to the command and if not
// responds with an error message.
func handleAccessOrError(permissionType models.PermissionType, s *discordgo.Session, i *discordgo.InteractionCreate) bool {
	if !userHasAccessToCommand(permissionType, i) {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "You don't have permission to perform this command",
			},
		})
		return false
	}
	return true
}

func respondToError(s *discordgo.Session, i *discordgo.InteractionCreate, err error) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "An error occurred while executing this command",
		},
	})
	logrus.WithFields(logrus.Fields{
		"error": err,
	}).Error("An error occurred while executing this command")
}

func healthHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if !handleAccessOrError(models.PermissionTypeCheckHealth, s, i) {
		return
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "🟢 I'm healthy!",
		},
	})
}
