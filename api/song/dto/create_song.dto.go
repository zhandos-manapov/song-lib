package dto

type CreateSongDto struct {
	Group string `validate:"required"`
	Song  string `validate:"required"`
}
