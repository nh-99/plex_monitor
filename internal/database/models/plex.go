package models

type PlexWebhookData struct {
	Event   string `bson:"event"`
	User    bool   `bson:"user"`
	Owner   bool   `bson:"owner"`
	Account struct {
		ID    int    `bson:"id"`
		Thumb string `bson:"thumb"`
		Title string `bson:"title"`
	} `bson:"Account"`
	Server struct {
		Title string `bson:"title"`
		UUID  string `bson:"uuid"`
	} `bson:"Server"`
	Player struct {
		Local         bool   `bson:"local"`
		PublicAddress string `bson:"publicAddress"`
		Title         string `bson:"title"`
		UUID          string `bson:"uuid"`
	} `bson:"Player"`
	Metadata struct {
		LibrarySectionType   string `bson:"librarySectionType"`
		RatingKey            string `bson:"ratingKey"`
		Key                  string `bson:"key"`
		ParentRatingKey      string `bson:"parentRatingKey"`
		GrandparentRatingKey string `bson:"grandparentRatingKey"`
		GUID                 string `bson:"guid"`
		LibrarySectionID     int    `bson:"librarySectionID"`
		Type                 string `bson:"type"`
		Title                string `bson:"title"`
		GrandparentKey       string `bson:"grandparentKey"`
		ParentKey            string `bson:"parentKey"`
		GrandparentTitle     string `bson:"grandparentTitle"`
		ParentTitle          string `bson:"parentTitle"`
		Summary              string `bson:"summary"`
		Index                int    `bson:"index"`
		ParentIndex          int    `bson:"parentIndex"`
		RatingCount          int    `bson:"ratingCount"`
		Thumb                string `bson:"thumb"`
		Art                  string `bson:"art"`
		ParentThumb          string `bson:"parentThumb"`
		GrandparentThumb     string `bson:"grandparentThumb"`
		GrandparentArt       string `bson:"grandparentArt"`
		AddedAt              int    `bson:"addedAt"`
		UpdatedAt            int    `bson:"updatedAt"`
	} `bson:"Metadata"`
}
