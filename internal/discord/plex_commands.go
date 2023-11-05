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
		libraryID := int(data.Options[0].IntValue())

		// Get the service config from the database
		service, err := models.GetServiceByName(models.ServiceTypePlex)
		if err != nil {
			respondToError(s, i, err)
			return
		}

		// Get the config as a Plex config
		config, err := service.GetConfigAsStandardConfig()
		if err != nil {
			respondToError(s, i, err)
			return
		}

		// Decrypt the service key
		plexKey, err := service.GetAndDecryptKey()
		if err != nil {
			respondToError(s, i, err)
			return
		}

		// Create the Plex driver
		plexDriver := servicerestdriver.NewPlexRestDriver("plex", config.Host, plexKey, logrus.WithField("service", "plex"))

		// Get the libraries
		libraries, err := plexDriver.GetLibraries()
		if err != nil {
			respondToError(s, i, fmt.Errorf("failed to get libraries: %w", err))
			return
		}

		// Find the library
		var library *servicerestdriver.PlexLibrary
		for _, l := range libraries {
			if l.Key == libraryID {
				library = &l
				break
			}
		}
		if library == nil {
			respondToError(s, i, fmt.Errorf("failed to find library with ID %d", libraryID))
			return
		}

		// Refresh the library
		err = plexDriver.ScanLibrary(library.Key)
		if err != nil {
			respondToError(s, i, err)
			return
		}

		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf(
					"Starting a rescan on the %s library",
					library.Title,
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
		logrus.Errorf("Failed to get service by name: %v", err)
		return
	}

	// Get the config as a Plex config
	config, err := service.GetConfigAsStandardConfig()
	if err != nil {
		logrus.Errorf("Failed to get config as plex config: %v", err)
		return
	}

	// Decrypt the service key
	plexKey, err := service.GetAndDecryptKey()
	if err != nil {
		logrus.Errorf("Failed to decrypt service key: %v", err)
		return
	}

	// Refresh the library
	plexDriver := servicerestdriver.NewPlexRestDriver("plex", config.Host, plexKey, logrus.WithField("service", "plex"))

	// Get the libraries
	libraries, err := plexDriver.GetLibraries()
	if err != nil {
		logrus.Errorf("Failed to get libraries: %v", err)
		return
	}

	// Add the libraries to the choices
	for _, library := range libraries {
		choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
			Name:  library.Title,
			Value: library.Key,
		})
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionApplicationCommandAutocompleteResult,
		Data: &discordgo.InteractionResponseData{
			Choices: choices,
		},
	})
	if err != nil {
		logrus.Errorf("Failed to respond to interaction: %v", err)
		return
	}
}
