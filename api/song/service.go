package song

import (
	"fmt"
	"math"
	"net/http"
	"net/url"
	"song-lib/api/group"
	"song-lib/api/song/dto"
	"song-lib/api/song/model"
	"song-lib/common"
	coredto "song-lib/common/dto"
	"song-lib/config"
	"strings"

	"github.com/jackc/pgx/v5"
)

type Service interface {
	Create(createSongDto dto.CreateSongDto) (*model.SongModel, error)
	Update(id string, updateSongDto dto.UpdateSongDto) (*model.SongModel, error)
	FindAll(pagination coredto.PaginationDto, findSongsDto dto.FindSongsDto) ([]*model.SongModel, error)
	FindOne(id string, versePagination coredto.PaginationDto) (*model.SongModel, error)
	Remove(id string) error
}

type service struct {
	common.BaseService
	store      Store
	groupStore group.Store
	env        config.Config
}

func NewService(store Store, groupStore group.Store, env config.Config) Service {
	return &service{
		BaseService: common.NewBaseService(),
		store:       store,
		groupStore:  groupStore,
		env:         env,
	}
}

func (s *service) getSongInfo(group string, song string) (*dto.SongApiResponseDto, error) {
	base, err := url.Parse(s.env.API_URL + "/info")
	if err != nil {
		return nil, err
	}

	params := url.Values{}
	params.Add("group", group)
	params.Add("song", song)
	base.RawQuery = params.Encode()

	resp, err := http.Get(base.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var songRespDto dto.SongApiResponseDto
	common.ParseJSON(resp.Body, &songRespDto)

	return &songRespDto, nil
}

func (s *service) Create(createSongDto dto.CreateSongDto) (*model.SongModel, error) {
	songRespDto, err := s.getSongInfo(createSongDto.Group, createSongDto.Song)
	if err != nil {
		return nil, err
	}

	group, err := s.groupStore.FindOneByName(createSongDto.Group)
	if err != nil && err == pgx.ErrNoRows {
		group, err = s.groupStore.Insert(createSongDto.Group)
		if err != nil {
			return nil, err
		}
	}

	group_id, _ := group.ID.Value()
	song, err := s.store.Insert(createSongDto.Song, group_id.(string), songRespDto.Text, songRespDto.Link, songRespDto.ReleaseDate)
	if err != nil {
		return nil, err
	}

	return song, nil
}

func (s *service) Update(id string, updateSongDto dto.UpdateSongDto) (*model.SongModel, error) {
	song, err := s.store.UpdateOne(
		id,
		updateSongDto.Name,
		updateSongDto.ReleaseDate,
		updateSongDto.Text,
		updateSongDto.Link,
	)
	if err == pgx.ErrNoRows {
		message := fmt.Sprintf("Song with id = %s not found", id)
		return nil, common.NewNotFoundError(message, err)
	} else if err != nil {
		return nil, err
	}
	return song, nil
}

func (s *service) FindAll(pagination coredto.PaginationDto, findSongsDto dto.FindSongsDto) ([]*model.SongModel, error) {
	songs, err := s.store.FindAll(pagination, findSongsDto)
	if err != nil {
		return nil, err
	}

	return songs, nil
}

func (s *service) FindOne(id string, versePagination coredto.PaginationDto) (*model.SongModel, error) {
	song, err := s.store.FindOne(id)
	if err == pgx.ErrNoRows {
		message := fmt.Sprintf("Song with id = %s not found", id)
		return nil, common.NewNotFoundError(message, err)
	} else if err != nil {
		return nil, err
	}

	// Pagination Handling
	if versePagination.Limit == 0 {
		versePagination.Limit = common.DEFAULT_LIMIT
	}

	verses := strings.Split(song.Text.String, "\n\n")
	versesLen := len(verses)
	minHigh := int(math.Min(float64(versePagination.Limit+versePagination.Skip), float64(versesLen)))
	textPaginatedSlice := verses[versePagination.Skip:minHigh]
	textPaginatedString := strings.Join(textPaginatedSlice, "\n\n")

	if err := song.Text.Scan(textPaginatedString); err != nil {
		return nil, err
	}
	// Pagination Handling End

	return song, nil
}

func (s *service) Remove(id string) error {
	if err := s.store.RemoveById(id); err != nil {
		return common.NewNotFoundError(err.Error(), err)
	}
	return nil
}
