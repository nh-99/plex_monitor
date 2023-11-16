package discord

import (
	"plex_monitor/internal/database/models"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

// requestMedia is a huge command. This allows a user to request new media through the bot.
// Essentially, a user can request new media via a simple command:
// /request param:content
// On the backend, we handle _all_ the magic - looking up the content, determining
// the type of media, and suggesting the most relevant content first back to the user.
// After that, we can send the request off to the relevant services. However, this process
// happens asynchronnously. From there the whole pipeline is monitored and feedback is given
// where neccessary. The user can also get the status of their request at any time - see
// getMediaRequestStatus
func requestMedia(s *discordgo.Session, i *discordgo.InteractionCreate) {
	l := logrus.NewEntry(logrus.StandardLogger())

	hasAccess, user := handleAccessOrError(models.PermissionTypeCheckHealth, s, i, l)
	if !hasAccess {
		return
	}

	data := i.ApplicationCommandData()
	requestTitle := data.Options[0].StringValue()

	// Create the media request object
	mediaRequest := models.MediaRequest{
		Name:          requestTitle,
		CurrentStatus: models.MediaRequestStatusRequested, // We always start in the requested state
		RequestedBy:   user.ID,
	}

	// Save the media request to the database
	err := mediaRequest.Save()
	if err != nil {
		respondToError(s, i, err, l)
		return
	}

	// TODO: lookup the content via TMDB and plex movie agent?
	// TODO: determine an ordering for the list - most relevant to least relevant
	// TODO: populate the list below with the options

	// Respond to the interaction with a prompt to choose the correct media
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.SelectMenu{
							// Select menu, as other components, must have a customID, so we set it to this value.
							CustomID:    "media-request-select",
							Placeholder: "Choose the correct item 👇",
							Options: []discordgo.SelectMenuOption{
								// TODO: populate based on backend query to external services
								{
									Label: "JS",
									Value: "js",
									Emoji: discordgo.ComponentEmoji{
										Name: "🟨",
									},
									Description: "JavaScript programming language",
								},
							},
						},
					},
				},
			},
			Flags: discordgo.MessageFlagsFailedToMentionSomeRolesInThread,
		},
	})

	// TODO: kick off a worker that will lookup plex libraries & map the request to the correct library
	// TODO: if we do the above, we will be cutting out ombi; if not:
	//
	// TODO: send the request to Ombi
	//
}

// getMediaRequestStatus returns the status of a request to the user.
func getMediaRequestStatus(s *discordgo.Session, i *discordgo.InteractionCreate) {
}
