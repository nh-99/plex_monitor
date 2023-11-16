package webhook

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"plex_monitor/internal/config"
	"plex_monitor/internal/database"
	"plex_monitor/internal/database/models"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func setup() {
	// Init logging
	logrus.SetReportCaller(true)
	logrus.SetLevel(logrus.DebugLevel)
	// Get config
	conf := config.GetConfig()
	// Initialize the database
	database.InitDB(conf.Database.ConnectionString, "plex_monitor_test")
}

func teardown() {
	// Drop the database
	database.DB.Drop(database.Ctx)

	// Close the database connection
	database.CloseDB()
}

func TestWebhookWithInvalidService(t *testing.T) {
	setup()
	defer teardown()

	// Create a request to pass to our
	// handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("POST", "/webhook?service=invalid", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Add the request body
	var jsonStr = []byte(`{"test": "test"}`)
	req.Body = io.NopCloser(bytes.NewBuffer(jsonStr))

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Entry)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Assert that the response was what we expected
	assert.Equal(t, http.StatusBadRequest, rr.Code, "handler returned wrong status code")
	// Assert that the response body was what we expected
	assert.Equal(t, "{\"status\":\"error\",\"message\":\"Invalid service\",\"success\":false}\n", rr.Body.String(), "handler returned unexpected body")
}

func TestWebhookWithPlexService(t *testing.T) {
	setup()
	defer teardown()

	// Create a request to pass to our
	// handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("POST", "/webhook?service=plex", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Construct multipart form request and add the json file to the "payload" field
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	f, err := os.Open("../../../../../test/plex_webhook_response_sample.json")
	assert.NoError(t, err)
	defer f.Close()

	// Parse the file data into a JSON object
	json, err := io.ReadAll(f)
	assert.NoError(t, err)

	// Create the form field
	fw, err := w.CreateFormField("payload")
	assert.NoError(t, err)

	// Write the json to the form field
	_, err = fw.Write(json)
	assert.NoError(t, err)

	// Close the writer
	w.Close()
	req.Body = io.NopCloser(&b)

	// Set the content type header
	req.Header.Set("Content-Type", w.FormDataContentType())

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Entry)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Assert that the response was what we expected
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}

	// Assert that we stored the event in the database
	evt, err := database.DB.Collection(database.WebhookCollectionName).CountDocuments(database.Ctx, bson.M{"event": "media.pause"})
	assert.NoError(t, err)
	assert.Equal(t, int64(1), evt)

	// Assert that we captured the raw data
	test := bson.M{"metadata.service": "plex"}
	raw, err := models.CountFilesInBucket(models.RawRequestWiresBucket, test)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), raw)
}

func TestWebhookWithSonarrService(t *testing.T) {
	setup()
	defer teardown()

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
	handler := http.HandlerFunc(Entry)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Assert that the response was what we expected using testify
	assert.Equal(t, http.StatusOK, rr.Code)

	// Assert that we stored the event in the database
	count, err := database.DB.Collection(database.WebhookCollectionName).CountDocuments(database.Ctx, bson.M{"series.id": 73})
	assert.NoError(t, err)
	assert.Equal(t, int64(1), count)

	// Assert that we captured the raw data
	test := bson.M{"metadata.service": "sonarr"}
	raw, err := models.CountFilesInBucket(models.RawRequestWiresBucket, test)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), raw)
}

func TestWebhookWithSonarrServiceHealth(t *testing.T) {
	setup()
	defer teardown()

	// Create a request to pass to our
	// handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("POST", "/webhook?service=sonarr", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Load json file and add it to the request body
	file, err := os.Open("../../../../../test/sonarr_webhook_response_sample_health_status.json")
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
	handler := http.HandlerFunc(Entry)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Assert that the response was what we expected using testify
	assert.Equal(t, http.StatusOK, rr.Code)

	// Assert that we stored the event in the database
	count, err := database.DB.Collection(database.WebhookCollectionName).CountDocuments(database.Ctx, bson.M{"message": "Indexers unavailable due to failures: indexerName"})
	assert.NoError(t, err)
	assert.Equal(t, int64(1), count)

	// Assert that we captured the raw data
	test := bson.M{"metadata.service": "sonarr"}
	raw, err := models.CountFilesInBucket(models.RawRequestWiresBucket, test)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), raw)
}

func TestWebhookWithRadarrService(t *testing.T) {
	setup()
	defer teardown()

	// Create a request to pass to our
	// handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("POST", "/webhook?service=radarr", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Load json file and add it to the request body
	file, err := os.Open("../../../../../test/radarr_webhook_response_sample__on_grab.json")
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
	handler := http.HandlerFunc(Entry)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Assert that the response was what we expected using testify
	assert.Equal(t, http.StatusOK, rr.Code)

	// Assert that we stored the event in the database
	count, err := database.DB.Collection(database.WebhookCollectionName).CountDocuments(database.Ctx, bson.M{"movie.id": 686})
	assert.NoError(t, err)
	assert.Equal(t, int64(1), count)

	// Assert that we captured the raw data
	test := bson.M{"metadata.service": "radarr"}
	raw, err := models.CountFilesInBucket(models.RawRequestWiresBucket, test)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), raw)
}

func TestWebhookWithRadarrServiceHealth(t *testing.T) {
	setup()
	defer teardown()

	// Create a request to pass to our
	// handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("POST", "/webhook?service=radarr", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Load json file and add it to the request body
	file, err := os.Open("../../../../../test/radarr_webhook_response_sample_health_status.json")
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
	handler := http.HandlerFunc(Entry)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Assert that the response was what we expected using testify
	assert.Equal(t, http.StatusOK, rr.Code)

	// Assert that we stored the event in the database
	count, err := database.DB.Collection(database.WebhookCollectionName).CountDocuments(database.Ctx, bson.M{"message": "Indexers unavailable due to failures: indexerName"})
	assert.NoError(t, err)
	assert.Equal(t, int64(1), count)

	// Assert that we captured the raw data
	test := bson.M{"metadata.service": "radarr"}
	raw, err := models.CountFilesInBucket(models.RawRequestWiresBucket, test)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), raw)
}

func TestWebhookWithOmbiService(t *testing.T) {
	setup()
	defer teardown()

	// Create a request to pass to our
	// handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("POST", "/webhook?service=ombi", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Load json file and add it to the request body
	file, err := os.Open("../../../../../test/ombi_webhook_response_sample.json")
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
	handler := http.HandlerFunc(Entry)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Assert that the response was what we expected using testify
	assert.Equal(t, http.StatusOK, rr.Code)

	// Assert that we stored the event in the database
	count, err := database.DB.Collection(database.WebhookCollectionName).CountDocuments(database.Ctx, bson.M{"requestId": "1234"})
	assert.NoError(t, err)
	assert.Equal(t, int64(1), count)

	// Assert that we captured the raw data
	test := bson.M{"metadata.service": "ombi"}
	raw, err := models.CountFilesInBucket(models.RawRequestWiresBucket, test)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), raw)
}
