package models

// SitewiseCustomMeta is the standard metadata
type SitewiseCustomMeta struct {
	NextToken  string   `json:"nextToken,omitempty"`
	EntryId    string   `json:"entryId,omitempty"`
	Resolution string   `json:"resolution,omitempty"`
	Aggregates []string `json:"aggregates,omitempty"`
}
