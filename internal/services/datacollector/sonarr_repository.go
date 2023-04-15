package datacollector

import "fmt"

type sonarrTvShow struct {
	absoluteEpisodeNumber      int
	airDate                    string
	airDateUtc                 string
	episodeFile                []string
	episodeFileId              int
	episodeNumber              int
	hasFile                    bool
	id                         int
	lastSearchTime             string
	monitored                  bool
	overview                   string
	sceneAbsoluteEpisodeNumber int
	sceneEpisodeNumber         int
	sceneSeasonNumber          int
	seasonNumber               int
	series                     []string
	seriesId                   int
	title                      string
	unverifiedSceneNumbering   bool
}

type SonarrCalendar struct {
	tvShows []sonarrTvShow
}

type SonarrQueue struct {
	tvShows []sonarrTvShow
}

func (s SonarrCalendar) collect() error {
	// Get calendar
	fmt.Println("calendar collect")
	return nil
}

func (s SonarrCalendar) store(db Database) error {
	fmt.Println("calendar store")
	return nil
}

func (s SonarrQueue) collect() error {
	// Get queue
	fmt.Println("queue collect")
	return nil
}

func (s SonarrQueue) store(db Database) error {
	// Store queue
	fmt.Println("queue store")
	return nil
}
