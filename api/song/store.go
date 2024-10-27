package song

import (
	"context"
	"fmt"
	"song-lib/api/song/dto"
	"song-lib/api/song/model"
	"song-lib/common"
	coredto "song-lib/common/dto"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store interface {
	Insert(
		name string,
		group_id string,
		text string,
		link string,
		release_date time.Time,
	) (*model.SongModel, error)
	UpdateOne(
		id string,
		name string,
		releaseDate time.Time,
		text string,
		link string,
	) (*model.SongModel, error)
	FindAll(pagination coredto.PaginationDto, findSongsDto dto.FindSongsDto) ([]*model.SongModel, error)
	FindOne(id string) (*model.SongModel, error)
	RemoveById(id string) error
}

type store struct {
	pool *pgxpool.Pool
}

func NewStore(pool *pgxpool.Pool) Store {
	return &store{
		pool: pool,
	}
}

func (s *store) Insert(
	name string,
	group_id string,
	text string,
	link string,
	release_date time.Time,
) (*model.SongModel, error) {
	query := `
  WITH song_temp AS (
  INSERT INTO "song" (name, group_id, text, link, release_date)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING *)
  
  SELECT song_temp.*, "group".name as group_name
  FROM song_temp
  INNER JOIN "group" ON song_temp.group_id = "group".id`

	var song model.SongModel
	if err := s.pool.QueryRow(
		context.Background(),
		query,
		name,
		group_id,
		text,
		link,
		release_date,
	).Scan(
		&song.ID,
		&song.Name,
		&song.GroupID,
		&song.Text,
		&song.Link,
		&song.ReleaseDate,
		&song.GroupName,
	); err != nil {
		return nil, err
	}
	return &song, nil
}

func (s *store) FindAll(pagination coredto.PaginationDto, findSongsDto dto.FindSongsDto) ([]*model.SongModel, error) {
	if pagination.Limit == 0 {
		pagination.Limit = common.DEFAULT_LIMIT
	}

	var query strings.Builder
	queryParams := []any{}

	query.WriteString(`
  SELECT "song".*, "group".name as group_name
	FROM "song"
	INNER JOIN "group" ON "song".group_id = "group".id`)

	if findSongsDto != (dto.FindSongsDto{}) {
		query.WriteString(" WHERE")

		if findSongsDto.Name != "" {
			if len(queryParams) > 0 {
				query.WriteString(" AND")
			}
			queryParams = append(queryParams, findSongsDto.Name)
			query.WriteString(fmt.Sprintf(` "song".name = $%d`, len(queryParams)))
		}

		if findSongsDto.GroupID != "" {
			if len(queryParams) > 0 {
				query.WriteString(" AND")
			}
			queryParams = append(queryParams, findSongsDto.GroupID)
			query.WriteString(fmt.Sprintf(` "song".group_id = $%d`, len(queryParams)))
		}

		if !findSongsDto.ReleaseDate.IsZero() {
			if len(queryParams) > 0 {
				query.WriteString(" AND")
			}
			queryParams = append(queryParams, findSongsDto.ReleaseDate)
			query.WriteString(fmt.Sprintf(` "song".release_date = $%d`, len(queryParams)))
		}

		if findSongsDto.GroupName != "" {
			if len(queryParams) > 0 {
				query.WriteString(" AND")
			}
			queryParams = append(queryParams, findSongsDto.GroupName)
			query.WriteString(fmt.Sprintf(` group_name = $%d`, len(queryParams)))
		}

		if findSongsDto.Search != "" {
			if len(queryParams) > 0 {
				query.WriteString(" AND")
			}
			queryParams = append(queryParams, fmt.Sprintf("%%%s%%", findSongsDto.Search))
			queryParamsLen := len(queryParams)
			query.WriteString(fmt.Sprintf(`
      (LOWER("song".name) LIKE LOWER($%d) OR 
      LOWER("song".text) LIKE LOWER($%d) OR 
      LOWER("group".name) LIKE LOWER($%d))`, queryParamsLen, queryParamsLen, queryParamsLen))
		}
	}

	queryParams = append(queryParams, pagination.Limit, pagination.Skip)
	query.WriteString(fmt.Sprintf(` LIMIT $%d OFFSET $%d`, len(queryParams)-1, len(queryParams)))

	rows, err := s.pool.Query(context.Background(), query.String(), queryParams...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	songs, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[model.SongModel])
	if err != nil {
		return nil, err
	}
	return songs, nil
}

func (s *store) FindOne(id string) (*model.SongModel, error) {
	query := `
	SELECT "song".*, "group".name as group_name 
	FROM "song" 
	INNER JOIN "group" ON "song".group_id = "group".id
	WHERE "song".id=$1`
	var song model.SongModel
	if err := s.pool.QueryRow(context.Background(), query, id).Scan(
		&song.ID,
		&song.Name,
		&song.GroupID,
		&song.Text,
		&song.Link,
		&song.ReleaseDate,
		&song.GroupName,
	); err != nil {
		return nil, err
	}
	return &song, nil
}

func (s *store) UpdateOne(
	id string,
	name string,
	releaseDate time.Time,
	text string,
	link string,
) (*model.SongModel, error) {
	query := `
  UPDATE "song" 
  SET name=$1, release_date=$2, text=$3, link=$4
  WHERE id=$5
  RETURNING *`
	var song model.SongModel
	if err := s.pool.QueryRow(
		context.Background(),
		query,
		name,
		releaseDate,
		text,
		link,
		id,
	).Scan(
		&song.ID,
		&song.Name,
		&song.GroupID,
		&song.Text,
		&song.Link,
		&song.ReleaseDate,
	); err != nil {
		return nil, err
	}
	return &song, nil
}

func (s *store) RemoveById(id string) error {
	query := `--sql
  DELETE FROM "song" WHERE id=$1`

	if tag, err := s.pool.Exec(context.Background(), query, id); err != nil {
		return err
	} else if tag.RowsAffected() < 1 {
		return fmt.Errorf("song with id = %s not found", id)
	}
	return nil
}
