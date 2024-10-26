package dto

import "time"

type SongApiResponseDto struct {
	ReleaseDate time.Time
	Text        string
	Link        string
}
