package model

import "github.com/jackc/pgx/v5/pgtype"

type GroupModel struct {
	ID   pgtype.UUID `db:"id"`
	Name pgtype.Text `db:"name"`
}
