package dto

import "time"

type SongResponseDto struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	ReleaseDate time.Time `json:"releaseDate"`
	Text        string    `json:"text"`
	Link        string    `json:"link"`
	GroupId     string    `json:"groupId"`
	GroupName   string    `json:"groupName"`
}
