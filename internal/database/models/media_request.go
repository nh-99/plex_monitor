package models

import (
	"plex_monitor/internal/database"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MediaRequestStatus is the status of a media request
type MediaRequestStatus string

const (
	// MediaRequestCollection is the name of the collection for media requests
	MediaRequestCollection string = "media_requests"
	// MediaRequestStatusRequested is the state for when a user submits a request for new media
	MediaRequestStatusRequested MediaRequestStatus = "requested"
	// MediaRequestStatusApproved is used when the user request is approved
	MediaRequestStatusApproved MediaRequestStatus = "approved"
	// MediaRequestStatusDeclined is used when the request is declined
	MediaRequestStatusDeclined MediaRequestStatus = "declined"
	// MediaRequestStatusClassified is used when the request has been classified into a library/type (e.g. movies, tv, anime, etc.)
	MediaRequestStatusClassified MediaRequestStatus = "classified"
	// MediaRequestStatusSearched is the status for a media request that has been searched for
	MediaRequestStatusSearched MediaRequestStatus = "searched"
	// MediaRequestStatusGrabbed is used when a search has been sent to the download client
	MediaRequestStatusGrabbed MediaRequestStatus = "grabbed"
	// MediaRequestStatusDownloaded is the status for a media request that has been downloaded
	MediaRequestStatusDownloaded MediaRequestStatus = "downloaded"
	// MediaRequestStatusAdded is the status for a media request that is currently added
	MediaRequestStatusAdded MediaRequestStatus = "added"
	// MediaRequestStatusIngested is the status for a media request that has been ingested
	MediaRequestStatusIngested MediaRequestStatus = "ingested"
	// MediaRequestStatusScanned is used when a request has been scanned into the library
	MediaRequestStatusScanned MediaRequestStatus = "scanned"
	// MediaRequestStatusNotified is used when a request has been notified
	MediaRequestStatusNotified MediaRequestStatus = "notified"
)

// MediaRequest is the model for a media request
type MediaRequest struct {
	ID            primitive.ObjectID        `json:"id,omitempty" bson:"_id,omitempty"`
	Name          string                    `json:"name" bson:"name"`
	StatusHistory []MediaRequestStatusState `json:"statusHistory" bson:"statusHistory"`
	CurrentStatus MediaRequestStatus        `json:"currentStatus" bson:"currentStatus"`
	RequestedBy   primitive.ObjectID        `json:"requestedBy" bson:"requestedBy"`
}

// MediaRequestStatusState is used to represent a status that the request had at one point.
// It can be used to construct a history of the status of a request.
type MediaRequestStatusState struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Status    MediaRequestStatus `json:"status" bson:"status"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
}

// GetMediaRequestForUser gets a specified media request for a user
func GetMediaRequestForUser(userID string, requestID string) (MediaRequest, error) {
	var mediaRequest MediaRequest

	// Convert the requestID to an ObjectID
	requestObjectID, err := primitive.ObjectIDFromHex(requestID)
	if err != nil {
		return MediaRequest{}, err
	}

	// Convert the userID to an ObjectID
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return MediaRequest{}, err
	}

	// Get the media request
	err = database.DB.
		Collection(MediaRequestCollection).
		FindOne(
			database.Ctx,
			bson.M{
				"_id":         requestObjectID,
				"requestedBy": userObjectID,
			},
		).
		Decode(&mediaRequest)
	if err != nil {
		return MediaRequest{}, err
	}

	return mediaRequest, nil
}

// Save persists all of the fields on the struct to the database
func (mr *MediaRequest) Save() error {
	if mr.ID.IsZero() {
		mr.ID = primitive.NewObjectID()
	}

	// Has the state changed? If so, add a new state to the history
	if len(mr.StatusHistory) == 0 || mr.StatusHistory[len(mr.StatusHistory)-1].Status != mr.CurrentStatus {
		mr.StatusHistory = append(mr.StatusHistory, MediaRequestStatusState{
			Status:    mr.CurrentStatus,
			CreatedAt: time.Now(),
		})
	}

	// Upsert the media request
	opts := options.Update().SetUpsert(true)
	_, err := database.DB.Collection(MediaRequestCollection).UpdateOne(database.Ctx, bson.M{"_id": mr.ID}, bson.M{"$set": mr}, opts)
	if err != nil {
		return err
	}

	return nil
}

// Reload gets all of the fields from the DB and populates the struct
func (mr MediaRequest) Reload() error {
	err := database.DB.
		Collection(MediaRequestCollection).
		FindOne(
			database.Ctx,
			bson.M{"_id": mr.ID},
		).
		Decode(&mr)
	if err != nil {
		return err
	}

	return nil
}
