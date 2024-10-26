package song

import (
	"fmt"
	"net/http"
	"song-lib/api/song/dto"
	"song-lib/common"
	coredto "song-lib/common/dto"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

type controller struct {
	common.BaseController
	service Service
}

func NewController(service Service) common.Controller {
	return &controller{
		BaseController: common.NewBaseController("/songs"),
		service:        service,
	}
}

func (c *controller) RegisterRoutes(r *chi.Router) {
	(*r).Post("/", common.MakeHTTPHandleFunc(c.create))
	(*r).Get("/{id}", common.MakeHTTPHandleFunc(c.findOne))
	(*r).Put("/{id}", common.MakeHTTPHandleFunc(c.update))
	(*r).Delete("/{id}", common.MakeHTTPHandleFunc(c.remove))
	(*r).Get("/", common.MakeHTTPHandleFunc(c.findAll))
}

func (c *controller) findOne(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")
	if err := common.ValidUUID(id); err != nil {
		return common.NewBadRequestError("Invalid uuid param", err)
	}

	versePagination, err := getPagination(r)
	if err != nil {
		return err
	}

	song, err := c.service.FindOne(id, *versePagination)
	if err != nil {
		return err
	}

	return common.WriteJSON(w, http.StatusOK, SongModelToDto(song))
}

func getPagination(r *http.Request) (*coredto.PaginationDto, error) {
	skip, _ := strconv.Atoi(r.URL.Query().Get("skip"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	pagination := coredto.PaginationDto{Skip: skip, Limit: limit}

	if err := common.ValidStruct(pagination); err != nil {
		return nil, err
	}
	return &pagination, nil
}

func getFindSongsDto(r *http.Request) (*dto.FindSongsDto, error) {

	dateStr := r.URL.Query().Get("releaseDate")
	var releaseDate time.Time
	var err error

	if dateStr != "" {
		releaseDate, err = time.Parse(time.DateOnly, dateStr)
		if err != nil {
			return nil, err
		}
	}

	findSongsDto := dto.FindSongsDto{
		Name:        r.URL.Query().Get("name"),
		ReleaseDate: releaseDate,
		GroupID:     r.URL.Query().Get("groupId"),
		GroupName:   r.URL.Query().Get("groupName"),
		Search:      r.URL.Query().Get("search"),
	}

	return &findSongsDto, nil
}

func (c *controller) findAll(w http.ResponseWriter, r *http.Request) error {
	pagination, err := getPagination(r)
	if err != nil {
		return common.NewBadRequestError("Invalid query params", err)
	}

	findSongsDto, err := getFindSongsDto(r)
	if err != nil {
		return common.NewBadRequestError("Invalid query params", err)
	}

	songs, err := c.service.FindAll(*pagination, *findSongsDto)
	if err != nil {
		return err
	}

	return common.WriteJSON(w, http.StatusOK, SongModelToDtoSlice(songs))
}

func (c *controller) create(w http.ResponseWriter, r *http.Request) error {
	createSongDto, err := common.ParseBody(r, &dto.CreateSongDto{})
	if err != nil {
		return common.NewBadRequestError(err.Error(), err)
	}

	song, err := c.service.Create(*createSongDto)
	if err != nil {
		return err
	}
	return common.WriteJSON(w, http.StatusOK, SongModelToDto(song))
}

func (c *controller) update(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")
	if err := common.ValidUUID(id); err != nil {
		return common.NewBadRequestError("Invalid uuid path variable", err)
	}

	updateSongDto, err := common.ParseBody(r, &dto.UpdateSongDto{})
	if err != nil {
		return common.NewBadRequestError(err.Error(), err)
	}

	song, err := c.service.Update(id, *updateSongDto)
	if err != nil {
		return err
	}
	return common.WriteJSON(w, http.StatusOK, SongModelToDto(song))
}

func (c *controller) remove(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")
	if err := common.ValidUUID(id); err != nil {
		return common.NewBadRequestError("Invalid uuid param", err)
	}

	if err := c.service.Remove(id); err != nil {
		return err
	}

	message := fmt.Sprintf("Song with id = %s successfully deleted", id)
	return common.WriteJSON(w, http.StatusOK, common.NewApiResponse(message))
}
