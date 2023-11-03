package discord

import (
	"fmt"
	"plex_monitor/internal/database/models"
	servicerestdriver "plex_monitor/internal/service_rest_driver"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

func refreshPlexLibraryHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if !handleAccessOrError(models.PermissionTypeScanLibrary, s, i) {
		return
	}

	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		data := i.ApplicationCommandData()

		// Get the service config from the database
		service, err := models.GetServiceByName(models.ServiceTypePlex)
		if err != nil {
			respondToError(s, i, err)
			return
		}

		// Get the config as a Plex config
		config, err := service.GetConfigAsPlexConfig()
		if err != nil {
			respondToError(s, i, err)
			return
		}

		// Refresh the library
		libraryID := int(data.Options[0].IntValue())
		plexDriver := servicerestdriver.NewPlexRestDriver("plex", config.Host, config.Key, logrus.WithField("service", "plex"))
		err = plexDriver.ScanLibrary(libraryID)
		if err != nil {
			respondToError(s, i, err)
			return
		}

		// Get the library name
		library, err := config.GetLibraryByID(libraryID)
		if err != nil {
			respondToError(s, i, err)
			return
		}

		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf(
					"Starting a rescan on the %s library",
					library.Name,
				),
			},
		})
		if err != nil {
			panic(err)
		}

		// TODO: send the admin discord user a message that a scan has been triggered
	case discordgo.InteractionApplicationCommandAutocomplete:
		respondWithPlexLibraryAutocomplete(s, i)
	}
}

func respondWithPlexLibraryAutocomplete(s *discordgo.Session, i *discordgo.InteractionCreate) {
	choices := make([]*discordgo.ApplicationCommandOptionChoice, 0, 25)

	// Get the service config from the database
	service, err := models.GetServiceByName(models.ServiceTypePlex)
	if err != nil {
		respondToError(s, i, err)
		return
	}

	// Get the config as a Plex config
	config, err := service.GetConfigAsPlexConfig()
	if err != nil {
		respondToError(s, i, err)
		return
	}

	// Get the libraries from the Plex config
	libraries := config.Libraries

	// Add the libraries to the choices
	for _, library := range libraries {
		choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
			Name:  library.Name,
			Value: fmt.Sprintf("%d", library.ID),
		})
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionApplicationCommandAutocompleteResult,
		Data: &discordgo.InteractionResponseData{
			Choices: choices,
		},
	})
	if err != nil {
		respondToError(s, i, err)
		return
	}
}
