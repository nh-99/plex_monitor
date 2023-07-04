package models

type OmbiWebhookData struct {
	RequestId                        string `json:"requestId"`
	RequestedUser                    string `json:"requestedUser"`
	Title                            string `json:"title"`
	RequestedDate                    string `json:"requestedDate"`
	Type                             string `json:"type"`
	AdditionalInformation            string `json:"additionalInformation"`
	LongDate                         string `json:"longDate"`
	ShortDate                        string `json:"shortDate"`
	LongTime                         string `json:"longTime"`
	ShortTime                        string `json:"shortTime"`
	Overview                         string `json:"overview"`
	Year                             string `json:"year"`
	EpisodesList                     string `json:"episodesList"`
	SeasonsList                      string `json:"seasonsList"`
	PosterImage                      string `json:"posterImage"`
	ApplicationName                  string `json:"applicationName"`
	ApplicationUrl                   string `json:"applicationUrl"`
	IssueDescription                 string `json:"issueDescription"`
	IssueCategory                    string `json:"issueCategory"`
	IssueStatus                      string `json:"issueStatus"`
	IssueSubject                     string `json:"issueSubject"`
	NewIssueComment                  string `json:"newIssueComment"`
	IssueUser                        string `json:"issueUser"`
	UserName                         string `json:"userName"`
	Alias                            string `json:"alias"`
	RequestedByAlias                 string `json:"requestedByAlias"`
	UserPreference                   string `json:"userPreference"`
	DenyReason                       string `json:"denyReason"`
	AvailableDate                    string `json:"availableDate"`
	RequestStatus                    string `json:"requestStatus"`
	ProviderId                       string `json:"providerId"`
	PartiallyAvailableEpisodeNumbers string `json:"partiallyAvailableEpisodeNumbers"`
	PartiallyAvailableSeasonNumber   string `json:"partiallyAvailableSeasonNumber"`
	PartiallyAvailableEpisodesList   string `json:"partiallyAvailableEpisodesList"`
	PartiallyAvailableEpisodeCount   string `json:"partiallyAvailableEpisodeCount"`
	NotificationType                 string `json:"notificationType"`
}
