package models

type SonarrWebhookData struct {
	Series struct {
		ID       int    `bson:"id"`
		Title    string `bson:"title"`
		Path     string `bson:"path"`
		TVDBID   int    `bson:"tvdbId"`
		TVMazeID int    `bson:"tvMazeId"`
		IMDBID   string `bson:"imdbId"`
		Type     string `bson:"type"`
	} `bson:"series"`
	Episodes []struct {
		ID            int    `bson:"id"`
		EpisodeNumber int    `bson:"episodeNumber"`
		SeasonNumber  int    `bson:"seasonNumber"`
		Title         string `bson:"title"`
		AirDate       string `bson:"airDate"`
		AirDateUtc    string `bson:"airDateUtc"`
	} `bson:"episodes"`
	Release struct {
		Quality        string `bson:"quality"`
		QualityVersion int    `bson:"qualityVersion"`
		ReleaseGroup   string `bson:"releaseGroup"`
		ReleaseTitle   string `bson:"releaseTitle"`
		Indexer        string `bson:"indexer"`
		Size           int    `bson:"size"`
	} `bson:"release"`
	DownloadClient     string `bson:"downloadClient"`
	DownloadClientType string `bson:"downloadClientType"`
	DownloadID         string `bson:"downloadId"`
	EventType          string `bson:"eventType"`
}
