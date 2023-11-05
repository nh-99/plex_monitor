package servicerestdriver

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

// PlexRestDriver is a REST driver to interact with Plex. It handles making all standard HTTP requests
type PlexRestDriver struct {
	ServiceRestDriver
}

// NewPlexRestDriver returns a new Plex rest driver.
func NewPlexRestDriver(name, host, key string, logger *logrus.Entry) *PlexRestDriver {
	return &PlexRestDriver{
		ServiceRestDriver: ServiceRestDriver{
			Name:    name,
			Host:    host,
			Key:     key,
			Client:  &http.Client{},
			Logger:  logger,
			Retries: 10,
			Backoff: 5 * time.Second,
		},
	}
}

// Do executes a request against Plex.
func (s *PlexRestDriver) Do(req *http.Request) (*http.Response, error) {
	// Set the X-Plex-Token in the URL query
	params := req.URL.Query()
	params.Set("X-Plex-Token", s.Key)
	req.URL.RawQuery = params.Encode()

	// Execute the request with the service rest driver
	return s.ExecuteRequestSafe(req)
}

// ScanLibrary scans the library with the given ID.
func (s *PlexRestDriver) ScanLibrary(libraryID int) error {
	// Create a new request
	req, err := http.NewRequest(http.MethodGet, s.Host+"/library/sections/"+strconv.Itoa(libraryID)+"/refresh", nil)
	if err != nil {
		return err
	}

	// Execute the request
	_, err = s.Do(req)
	if err != nil {
		return err
	}

	return nil
}

type plexLibraryResponse struct {
	XMLName   xml.Name      `xml:"MediaContainer"`
	Directory []PlexLibrary `xml:"Directory"`
}

// PlexLibrary is the struct that represents the Plex library.
type PlexLibrary struct {
	XMLName          xml.Name `xml:"Directory"`
	Key              int      `xml:"key,attr"`
	Title            string   `xml:"title,attr"`
	Type             string   `xml:"type,attr"`
	ScannedAt        string   `xml:"scannedAt,attr"`
	CreatedAt        string   `xml:"createdAt,attr"`
	UpdatedAt        string   `xml:"updatedAt,attr"`
	ContentChangedAt string   `xml:"contentChangedAt,attr"`
}

// GetLibraries returns all libraries.
func (s *PlexRestDriver) GetLibraries() ([]PlexLibrary, error) {
	// Create a new request
	req, err := http.NewRequest(http.MethodGet, s.Host+"/library/sections", nil)
	if err != nil {
		return nil, err
	}

	// Execute the request
	resp, err := s.Do(req)
	if err != nil {
		return nil, err
	}

	// Parse the response
	var response plexLibraryResponse
	err = parseXML(resp, &response)
	if err != nil {
		return nil, err
	}

	return response.Directory, nil
}

func parseXML(resp *http.Response, v interface{}) error {
	// Unmarshal the response
	err := xml.NewDecoder(resp.Body).Decode(v)
	if err != nil {
		return err
	}

	// Close the response body
	err = resp.Body.Close()

	return nil
}

// PlexServerCapabilitiesResponse has data about the Plex server.
type PlexServerCapabilitiesResponse struct {
	XMLName                       xml.Name `xml:"MediaContainer"`
	Size                          int      `xml:"size,attr"`
	AllowCameraUpload             int      `xml:"allowCameraUpload,attr"`
	AllowChannelAccess            int      `xml:"allowChannelAccess,attr"`
	AllowMediaDeletion            int      `xml:"allowMediaDeletion,attr"`
	AllowSharing                  int      `xml:"allowSharing,attr"`
	AllowSync                     int      `xml:"allowSync,attr"`
	AllowTuners                   int      `xml:"allowTuners,attr"`
	BackgroundProcessing          int      `xml:"backgroundProcessing,attr"`
	CompanionProxy                int      `xml:"companionProxy,attr"`
	CountryCode                   string   `xml:"countryCode,attr"`
	Diagnostics                   string   `xml:"diagnostics,attr"`
	EventStream                   int      `xml:"eventStream,attr"`
	FriendlyName                  string   `xml:"friendlyName,attr"`
	HubSearch                     int      `xml:"hubSearch,attr"`
	ItemClusters                  int      `xml:"itemClusters,attr"`
	LiveTV                        int      `xml:"livetv,attr"`
	MachineIdentifier             string   `xml:"machineIdentifier,attr"`
	MediaProviders                int      `xml:"mediaProviders,attr"`
	Multiuser                     int      `xml:"multiuser,attr"`
	MusicAnalysis                 int      `xml:"musicAnalysis,attr"`
	MyPlex                        int      `xml:"myPlex,attr"`
	MyPlexMappingState            string   `xml:"myPlexMappingState,attr"`
	MyPlexSigninState             string   `xml:"myPlexSigninState,attr"`
	MyPlexSubscription            int      `xml:"myPlexSubscription,attr"`
	MyPlexUsername                string   `xml:"myPlexUsername,attr"`
	OfflineTranscode              int      `xml:"offlineTranscode,attr"`
	OwnerFeatures                 string   `xml:"ownerFeatures,attr"`
	PhotoAutoTag                  int      `xml:"photoAutoTag,attr"`
	Platform                      string   `xml:"platform,attr"`
	PlatformVersion               string   `xml:"platformVersion,attr"`
	PluginHost                    int      `xml:"pluginHost,attr"`
	PushNotifications             int      `xml:"pushNotifications,attr"`
	ReadOnlyLibraries             int      `xml:"readOnlyLibraries,attr"`
	StreamingBrainABRVersion      int      `xml:"streamingBrainABRVersion,attr"`
	StreamingBrainVersion         int      `xml:"streamingBrainVersion,attr"`
	Sync                          int      `xml:"sync,attr"`
	TranscoderActiveVideoSessions int      `xml:"transcoderActiveVideoSessions,attr"`
	TranscoderAudio               int      `xml:"transcoderAudio,attr"`
	TranscoderLyrics              int      `xml:"transcoderLyrics,attr"`
	TranscoderPhoto               int      `xml:"transcoderPhoto,attr"`
	TranscoderSubtitles           int      `xml:"transcoderSubtitles,attr"`
	TranscoderVideo               int      `xml:"transcoderVideo,attr"`
	TranscoderVideoBitrates       string   `xml:"transcoderVideoBitrates,attr"`
	TranscoderVideoQualities      string   `xml:"transcoderVideoQualities,attr"`
	TranscoderVideoResolutions    string   `xml:"transcoderVideoResolutions,attr"`
	UpdatedAt                     int      `xml:"updatedAt,attr"`
	Updater                       int      `xml:"updater,attr"`
	Version                       string   `xml:"version,attr"`
	VoiceSearch                   int      `xml:"voiceSearch,attr"`
}

// GetServerCapabilities returns the server capabilities.
func (s *PlexRestDriver) GetServerCapabilities() (*PlexServerCapabilitiesResponse, error) {
	// Create a new request
	req, err := http.NewRequest(http.MethodGet, s.Host+"/", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create new request: %w", err)
	}

	// Execute the request
	resp, err := s.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}

	// Parse the response
	var response PlexServerCapabilitiesResponse
	err = parseXML(resp, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}

// GetHealth returns the health of the Plex server.
func (s *PlexRestDriver) GetHealth() (ServiceHealth, error) {
	// Get the server capabilities
	capabilities, err := s.GetServerCapabilities()
	if err != nil {
		return ServiceHealth{
			Healthy:     false,
			Version:     "N/A",
			LastChecked: time.Now(),
		}, fmt.Errorf("failed to get server capabilities: %w", err)
	}

	// Create the health struct
	health := ServiceHealth{
		Healthy:     true,
		Version:     capabilities.Version,
		LastChecked: time.Now(),
	}

	return health, nil
}
