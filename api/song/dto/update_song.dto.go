package dto

import "time"

type UpdateSongDto struct {
	Name        string    `validate:"required"`
	ReleaseDate time.Time `validate:"required"`
	Text        string    `validate:"required"`
	Link        string    `validate:"required"`
}
