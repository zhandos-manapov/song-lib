package model

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type SongModel struct {
	ID          pgtype.UUID `db:"id"`
	Name        pgtype.Text `db:"name"`
	ReleaseDate time.Time   `db:"release_date"`
	Text        pgtype.Text `db:"text"`
	Link        pgtype.Text `db:"link"`
	GroupID     pgtype.UUID `db:"group_id"`
	GroupName   pgtype.Text `db:"group_name"`
}
