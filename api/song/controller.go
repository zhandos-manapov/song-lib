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

// FindOne	Find One Song
//
//	@Summary		Find One Song
//	@Description	finds one song with verse pagination
//	@Tags			songs
//	@Accept			json
//	@Produce		json
//	@Param			id			path		string	true	"Id of a song"
//	@Param			skip		query		int		false	"Skips this many verses"
//	@Param			limit		query		int		false	"Gets this many verses. If omitted, the default is 10"
//	@Success		200			{object}	dto.SongResponseDto
//	@Failure		400			{object}	common.apiError
//	@Failure		404			{object}	common.apiError
//	@Failure		500			{object}	common.apiError
//	@Router			/songs/{id}																											[get]
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

// FindAll Find All Songs
//
//	@Summary		Find All Songs
//	@Description	finds all songs with pagination
//	@Tags			songs
//	@Accept			json
//	@Produce		json
//	@Param			skip	query		int	false	"Skips this many entries"
//	@Param			limit	query		int	false	"Gets this many entries. If omitted, the default is 10"
//	@Success		200		{array}		dto.SongResponseDto
//	@Failure		400		{object}	common.apiError
//	@Failure		500		{object}	common.apiError
//	@Router			/songs	[get]
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

// Create	Create a Song
//
//	@Summary		Create a Song
//	@Description	creates a song
//	@Tags			songs
//	@Accept			json
//	@Produce		json
//	@Param			createSongDto	body		dto.CreateSongDto	true	"Body for creating a song"
//	@Success		200				{object}	dto.SongResponseDto
//	@Failure		400				{object}	common.apiError
//	@Failure		500				{object}	common.apiError
//	@Router			/songs																																							[post]
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

// Update	Update a Song
//
//	@Summary		Update a Song
//	@Description	updates song details
//	@Tags			songs
//	@Accept			json
//	@Produce		json
//	@Param			id				path		string				true	"Id of a song"
//	@Param			updateSongDto	body		dto.UpdateSongDto	true	"Body for updating a song"
//	@Success		200				{object}	dto.SongResponseDto
//	@Failure		400				{object}	common.apiError
//	@Failure		404				{object}	common.apiError
//	@Failure		500				{object}	common.apiError
//	@Router			/songs/{id}																																[put]
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

// Remove	Remove a Song
//
//	@Summary		Remove a Song
//	@Description	removes a song
//	@Tags			songs
//	@Accept			json
//	@Produce		json
//	@Param			id			path		string	true	"Id of a song"
//	@Success		200			{object}	common.apiResponse
//	@Failure		400			{object}	common.apiError
//	@Failure		404			{object}	common.apiError
//	@Failure		500			{object}	common.apiError
//	@Router			/songs/{id}																											[delete]
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
