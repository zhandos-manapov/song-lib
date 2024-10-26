package song

import (
	"song-lib/api/song/dto"
	"song-lib/api/song/model"
)

func SongModelToDto(song *model.SongModel) *dto.SongResponseDto {
	songId, _ := song.ID.Value()
	groupId, _ := song.GroupID.Value()

	return &dto.SongResponseDto{
		ID:          songId.(string),
		Name:        song.Name.String,
		ReleaseDate: song.ReleaseDate,
		Text:        song.Text.String,
		Link:        song.Link.String,
		GroupId:     groupId.(string),
		GroupName:   song.GroupName.String,
	}
}

func SongModelToDtoSlice(songs []*model.SongModel) []*dto.SongResponseDto {
	songsDto := make([]*dto.SongResponseDto, len(songs))

	for i := range songs {
		songsDto[i] = SongModelToDto(songs[i])
	}
	return songsDto
}
