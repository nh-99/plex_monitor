package pipeline

import (
	"errors"
	"fmt"
	"plex_monitor/internal/database/models"

	"go.mongodb.org/mongo-driver/mongo"
)

const (
	// MediaRequestRequested is the event that is fired when a user requests a movie or TV show.
	MediaRequestRequested = "user_requested"

	// MediaRequestIngestedBySonarr is the event that is fired when a request is ingested by Sonarr.
	MediaRequestIngestedBySonarr = "request_ingested_by_sonarr"

	// MediaRequestIngestedByRadarr is the event that is fired when a request is ingested by Radarr.
	MediaRequestIngestedByRadarr = "request_ingested_by_radarr"

	// MediaRequestDownloading is the event that is fired when a request download is started.
	MediaRequestDownloading = "request_downloading"

	// MediaRequestDownloaded is the event that is fired when a request download is completed.
	MediaRequestDownloaded = "request_downloaded"

	// MediaRequestImported is the event that is fired when a request is imported into Plex.
	MediaRequestImported = "request_imported"

	ackKey = "acknowledgeWith"
)

// MediaRequestPipeline is the pipeline that all of the media requests flow through.
type MediaRequestPipeline struct {
	Pipeline
}

// MediaRequest is the pipeline that all of the media requests flow through.
type MediaRequest interface {
	// Requested is the event that is fired when a user requests a movie or TV show.
	Requested() error

	// RequestIngestedBySonarr is the event that is fired when a request is ingested by Sonarr.
	RequestIngestedBySonarr() error

	// RequestIngestedByRadarr is the event that is fired when a request is ingested by Radarr.
	RequestIngestedByRadarr() error

	// RequestDownloading is the event that is fired when a request download is started.
	RequestDownloading() error

	// RequestDownloaded is the event that is fired when a request download is completed.
	RequestDownloaded() error

	// RequestImported is the event that is fired when a request is imported into Plex.
	RequestImported() error
}

// NewMediaRequestPipeline returns a new media request pipeline.
func NewMediaRequestPipeline(id string) *MediaRequestPipeline {
	return &MediaRequestPipeline{
		Pipeline: Pipeline{
			ID:       id,
			Name:     "Media Request",
			Metadata: make(map[string]interface{}),
		},
	}
}

func getStepFunction(mediaRequestPipeline MediaRequestPipeline, key string) func() error {
	switch key {
	case MediaRequestRequested:
		return mediaRequestPipeline.Requested
	case MediaRequestIngestedBySonarr:
		return mediaRequestPipeline.RequestIngestedBySonarr
	case MediaRequestIngestedByRadarr:
		return mediaRequestPipeline.RequestIngestedByRadarr
	case MediaRequestDownloading:
		return mediaRequestPipeline.RequestDownloading
	case MediaRequestDownloaded:
		return mediaRequestPipeline.RequestDownloaded
	case MediaRequestImported:
		return mediaRequestPipeline.RequestImported
	}

	return nil
}

// CreateMediaRequestPipeline creates a new media request pipeline.
func CreateMediaRequestPipeline(id string) *MediaRequestPipeline {
	pipeline := NewMediaRequestPipeline(id)

	// Step 1: User Requested
	pipeline.Pipeline.AddStep("User Requested", MediaRequestRequested, getStepFunction(*pipeline, MediaRequestRequested))
	// Step 2: Request Ingested by Sonarr
	pipeline.Pipeline.AddStep("Request Ingested by Sonarr", MediaRequestIngestedBySonarr, getStepFunction(*pipeline, MediaRequestIngestedBySonarr))
	// Step 3: Request Ingested by Radarr
	pipeline.Pipeline.AddStep("Request Ingested by Radarr", MediaRequestIngestedByRadarr, getStepFunction(*pipeline, MediaRequestIngestedByRadarr))
	// Step 4: Request Downloading
	pipeline.Pipeline.AddStep("Request Downloading", MediaRequestDownloading, getStepFunction(*pipeline, MediaRequestDownloading))
	// Step 5: Request Downloaded
	pipeline.Pipeline.AddStep("Request Downloaded", MediaRequestDownloaded, getStepFunction(*pipeline, MediaRequestDownloaded))
	// Step 6: Request Imported
	pipeline.Pipeline.AddStep("Request Imported", MediaRequestImported, getStepFunction(*pipeline, MediaRequestImported))

	return pipeline
}

// GetMediaRequestPipelineByID gets a media request pipeline from the database by ID.
func GetMediaRequestPipelineByID(id string) (*MediaRequestPipeline, error) {
	pipeline, err := GetPipelineByID(id)
	if err != nil {
		return nil, err
	}

	mediaRequestPipeline := &MediaRequestPipeline{
		Pipeline: *pipeline,
	}

	// Add functions to the steps
	for i, step := range mediaRequestPipeline.Steps {
		stepFunc := getStepFunction(*mediaRequestPipeline, step.Key)
		if stepFunc == nil {
			return nil, fmt.Errorf("unable to find step function for step %s", step.Key)
		}
		mediaRequestPipeline.Steps[i].Function = stepFunc
	}

	return mediaRequestPipeline, nil
}

// GetOrCreateMediaRequestPipeline gets or creates a media request pipeline from the database.
func GetOrCreateMediaRequestPipeline(id string) (*Pipeline, error) {
	pipeline, err := GetMediaRequestPipelineByID(id)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, err
	}

	if pipeline == nil {
		pipeline = CreateMediaRequestPipeline(id)
		err = pipeline.Save()
		if err != nil {
			return nil, err
		}
	}

	return &pipeline.Pipeline, nil
}

// Requested is the event that is fired when a user requests a movie or TV show.
func (p *MediaRequestPipeline) Requested() error {
	// Validate that the metadata has the ombi key.
	if _, ok := p.Metadata["ombi"]; !ok {
		return errors.New("metadata does not have ombi key")
	}

	mediaRequestService := "ombi"
	// Retrieve the service from the database that holds data on the media request.
	service, err := models.GetServiceByName(mediaRequestService)
	if err != nil {
		return err
	}

	// Retrieve the users from the media request service.
	users, err := models.RetrieveUsersFromOmbi(service)
	// TODO: Catch error & mark pipeline as errored.

	// Filter the users down to the one that requested the media.
	ombiData := p.Metadata["ombi"].(map[string]interface{})
	userName := ombiData["userName"].(string)
	var user models.OmbiUser
	for _, u := range users {
		if u.UserName == userName {
			user = u
			break
		}
	}

	// Get the user's Discord ID from the Ombi data.
	discordID, err := models.RetrieveDiscordIDFromOmbi(service, user)
	if err != nil {
		return fmt.Errorf("unable to retrieve Discord ID from Ombi: %w", err)
	}

	// Set the Discord ID in the metadata.
	discordMetadata := map[string]interface{}{
		"discord": []map[string]interface{}{
			{
				"id": discordID,
			},
		},
	}
	p.AddMetadata(ackKey, discordMetadata)

	// Add additional Ombi metadata to the pipeline.
	ombiMetadata := p.GetMetadata("ombi").(map[string]interface{})
	ombiMetadata["userID"] = user.ID
	ombiMetadata["alias"] = user.Alias
	ombiMetadata["emailAddress"] = user.EmailAddress
	p.AddMetadata("ombi", ombiMetadata)

	return nil
}

// RequestIngestedBySonarr is the event that is fired when a request is ingested by Sonarr.
func (p *MediaRequestPipeline) RequestIngestedBySonarr() error {
	return nil
}

// RequestIngestedByRadarr is the event that is fired when a request is ingested by Radarr.
func (p *MediaRequestPipeline) RequestIngestedByRadarr() error {
	return nil
}

// RequestDownloading is the event that is fired when a request download is started.
func (p *MediaRequestPipeline) RequestDownloading() error {
	return nil
}

// RequestDownloaded is the event that is fired when a request download is completed.
func (p *MediaRequestPipeline) RequestDownloaded() error {
	return nil
}

// RequestImported is the event that is fired when a request is imported into Plex.
func (p *MediaRequestPipeline) RequestImported() error {
	return nil
}
