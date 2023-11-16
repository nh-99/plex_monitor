package mediarequest_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"plex_monitor/internal/config"
	mediarequest "plex_monitor/internal/controllers/api/media_request"
	"plex_monitor/internal/database"
	"plex_monitor/internal/database/models"

	"testing"

	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func setup(t *testing.T) models.User {
	// Init logging
	logrus.SetReportCaller(true)
	logrus.SetLevel(logrus.DebugLevel)
	// Get config
	conf := config.GetConfig()
	// Initialize the database
	database.InitDB(conf.Database.ConnectionString, "plex_monitor_test")

	// Create the test user
	user := models.User{
		Email: "foo@foo.foo",
		FrontendServices: []models.FrontendService{
			{
				Type: models.FrontendServiceTypeWeb,
			},
		},
		Permissions: []models.PermissionType{models.PermissionTypeGodMode}, // Test user has god mode permission
	}
	err := user.Save()
	if err != nil {
		t.Fatal(err)
	}

	err = user.Reload()
	if err != nil {
		t.Fatal(err)
	}

	return user
}

func teardown(t *testing.T) {
	// Drop the database
	database.DB.Drop(database.Ctx)

	// Close the database connection
	database.CloseDB()
}

func TestMediaRequest(t *testing.T) {
	user := setup(t)
	t.Cleanup(func() {
		teardown(t)
	})

	t.Run("Check that the request is successful", func(t *testing.T) {
		// Construct the request payload
		payload := map[string]interface{}{
			"name": "Predator",
		}
		req, err := http.NewRequest("POST", "/request/create", nil)
		if err != nil {
			t.Fatal(err)
		}

		// Create a token for the test user
		token, err := user.GetBearerToken()
		assert.NoError(t, err, "unexpected error while creating token")

		// Marshal payload into JSON
		jsonData, err := json.Marshal(payload)
		if err != nil {
			t.Fatal(err)
		}

		// Add the payload to the request body as JSON
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		req.Body = io.NopCloser(bytes.NewBuffer(jsonData))

		// Create a response recorder so you can inspect the response
		rr := httptest.NewRecorder()
		// Create a new chi router
		r := chi.NewRouter()
		r.Mount("/request", mediarequest.Routes())
		r.ServeHTTP(rr, req)

		// Check the status code is what we expect.
		assert.Equal(t, http.StatusOK, rr.Code, "handler returned wrong status code")

		// Unpack response body into a map[string]interface{}
		var response map[string]interface{}
		err = json.NewDecoder(rr.Body).Decode(&response)
		if err != nil {
			t.Fatal(err)
		}

		// Check that the response has the correct values
		assert.Equal(t, "success", response["status"], "response status is not correct")
		assert.Equal(t, "The request has been added to the queue", response["message"], "response message is not correct")

		// Check that the media request was created
		mediaRequest, err := models.GetMediaRequestForUser(user.ID.Hex(), response["id"].(string))
		if err != nil {
			t.Fatal(err)
		}

		// Check that the media request has the correct fields
		assert.Equal(t, mediaRequest.Name, "Predator", "media request name is not correct")
		assert.Equal(t, mediaRequest.CurrentStatus, models.MediaRequestStatusRequested, "media request status is not correct")
		assert.Equal(t, mediaRequest.RequestedBy, user.ID, "media request requested by is not correct")
	})
}
