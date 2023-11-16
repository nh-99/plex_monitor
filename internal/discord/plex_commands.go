package discord

import (
	"fmt"
	"plex_monitor/internal/database/models"
	servicerestdriver "plex_monitor/internal/service_rest_driver"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

func refreshPlexLibraryHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	l := logrus.NewEntry(logrus.StandardLogger())

	hasAccess, _ := handleAccessOrError(models.PermissionTypeCheckHealth, s, i, l)
	if !hasAccess {
		return
	}

	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		data := i.ApplicationCommandData()
		libraryID := int(data.Options[0].IntValue())

		// Get the service config from the database
		service, err := models.GetServiceByName(models.ServiceTypePlex)
		if err != nil {
			respondToError(s, i, err, l)
			return
		}

		// Get the config as a Plex config
		config, err := service.GetConfigAsStandardConfig()
		if err != nil {
			respondToError(s, i, err, l)
			return
		}

		// Decrypt the service key
		plexKey, err := service.GetAndDecryptKey()
		if err != nil {
			respondToError(s, i, err, l)
			return
		}

		// Create the Plex driver
		plexDriver := servicerestdriver.NewPlexRestDriver("plex", config.Host, plexKey, logrus.WithField("service", "plex"))

		// Get the libraries
		libraries, err := plexDriver.GetLibraries()
		if err != nil {
			respondToError(s, i, fmt.Errorf("failed to get libraries: %w", err), l)
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
			respondToError(s, i, fmt.Errorf("failed to find library with ID %d", libraryID), l)
			return
		}

		// Refresh the library
		err = plexDriver.ScanLibrary(library.Key)
		if err != nil {
			respondToError(s, i, err, l)
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
			l.WithError(err).Errorf("Failed to respond to rescan interaction, but rescan was triggered")
			return
		}

		// TODO: send the admin discord user a message that a scan has been triggered
	case discordgo.InteractionApplicationCommandAutocomplete:
		respondWithPlexLibraryAutocomplete(s, i, l)
	}
}

func respondWithPlexLibraryAutocomplete(s *discordgo.Session, i *discordgo.InteractionCreate, l *logrus.Entry) {
	choices := make([]*discordgo.ApplicationCommandOptionChoice, 0, 25)

	// Get the service config from the database
	service, err := models.GetServiceByName(models.ServiceTypePlex)
	if err != nil {
		l.WithError(err).Errorf("Failed to get service by name: %v", err)
		return
	}

	// Get the config as a Plex config
	config, err := service.GetConfigAsStandardConfig()
	if err != nil {
		l.WithError(err).Error("Failed to get config as plex config")
		return
	}

	// Decrypt the service key
	plexKey, err := service.GetAndDecryptKey()
	if err != nil {
		l.WithError(err).Error("Failed to decrypt service key")
		return
	}

	// Refresh the library
	plexDriver := servicerestdriver.NewPlexRestDriver("plex", config.Host, plexKey, logrus.WithField("service", "plex"))

	// Get the libraries
	libraries, err := plexDriver.GetLibraries()
	if err != nil {
		l.WithError(err).Error("Failed to get libraries")
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
		l.WithError(err).Error("Failed to respond to interaction")
		return
	}
}
