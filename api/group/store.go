package group

import (
	"context"
	"song-lib/api/group/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Store interface {
	Insert(name string) (*model.GroupModel, error)
	FindOneByName(name string) (*model.GroupModel, error)
}

type store struct {
	pool *pgxpool.Pool
}

func NewStore(pool *pgxpool.Pool) Store {
	return &store{
		pool: pool,
	}
}

func (s *store) Insert(name string) (*model.GroupModel, error) {
	query := `
	INSERT INTO "group" (name)
	VALUES ($1)
	RETURNING id, name`

	var group model.GroupModel
	if err := s.pool.QueryRow(context.Background(), query, name).Scan(&group.ID, &group.Name); err != nil {
		return nil, err
	}

	return &group, nil
}

func (s *store) FindOneByName(name string) (*model.GroupModel, error) {
	query := `
  SELECT id, name 
  FROM "group"
  WHERE name=$1`

	var group model.GroupModel
	if err := s.pool.QueryRow(context.Background(), query, name).Scan(&group.ID, &group.Name); err != nil {
		return nil, err
	}

	return &group, nil
}
