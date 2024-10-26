package dto

import "time"

type FindSongsDto struct {
	Name        string
	GroupID     string `validate:"uuid4"`
	ReleaseDate time.Time

	Search    string
	GroupName string
}
