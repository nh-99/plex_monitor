package firehose_test

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"plex_monitor/internal/controllers/api/firehose"
	"plex_monitor/internal/controllers/api/webhook"
	"plex_monitor/internal/database"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func setup() {
	// Init logging
	logrus.SetReportCaller(true)
	logrus.SetLevel(logrus.DebugLevel)
	// Initialize the database
	database.InitDB(os.Getenv("DATABASE_URL"), "plex_monitor_test")
}

func teardown() {
	// Drop the database
	database.DB.Drop(database.Ctx)

	// Close the database connection
	database.CloseDB()
}

func TestFirehose(t *testing.T) {
	setup()
	defer teardown()

	// BEGIN DATA SEEDING
	// Create a request to pass to our
	// handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("POST", "/webhook?service=sonarr", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Load json file and add it to the request body
	file, err := os.Open("../../../../../test/sonarr_webhook_response_sample__on_grab.json")
	assert.NoError(t, err)
	defer file.Close()
	// Read the file contents
	contents, err := io.ReadAll(file)
	assert.NoError(t, err)
	// Convert byte slice to string
	jsonString := string(contents)

	req.Body = io.NopCloser(bytes.NewBufferString(jsonString))

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(webhook.Entry)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Assert that the response was what we expected using testify
	assert.Equal(t, http.StatusOK, rr.Code)

	// END DATA SEEDING

	// Create a request to pass to our
	// handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err = http.NewRequest("GET", "/firehose", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	handler = http.HandlerFunc(firehose.Firehose)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Assert that the response was what we expected
	assert.Equal(t, http.StatusOK, rr.Code, "handler returned wrong status code")
	// Assert that the response body was what we expected
	assert.Contains(t, rr.Body.String(), `/tv/tv/Doctor Who (1963)`, "Did not find expected string in response body")

	// BEGIN DATA SEEDING AGAIN
	req, err = http.NewRequest("POST", "/webhook?service=radarr", nil)
	assert.NoError(t, err, "Creating radarr seed request should not error")

	// Load json file and add it to the request body
	file, err = os.Open("../../../../../test/radarr_webhook_response_sample__on_grab.json")
	assert.NoError(t, err)
	defer file.Close()
	// Read the file contents
	contents, err = io.ReadAll(file)
	assert.NoError(t, err)
	// Convert byte slice to string
	jsonString = string(contents)

	req.Body = io.NopCloser(bytes.NewBufferString(jsonString))

	handler = http.HandlerFunc(webhook.Entry)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Assert that the response was what we expected using testify
	assert.Equal(t, http.StatusOK, rr.Code)
	// END DATA SEEDING AGAIN

	// Create a request to pass to our
	// handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err = http.NewRequest("GET", "/firehose", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	handler = http.HandlerFunc(firehose.Firehose)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Assert that the response was what we expected
	assert.Equal(t, http.StatusOK, rr.Code, "handler returned wrong status code")
	// Assert that the order of responses looks correct
	sonnarIndex := strings.Index(rr.Body.String(), "sonarr")
	radarrIndex := strings.Index(rr.Body.String(), "radarr")
	assert.Greater(t, radarrIndex, sonnarIndex, "The ordering was incorrect on the response data")
}
