package main

import (
	"context"
	"song-lib/api/group"
	"song-lib/api/song"
	"song-lib/common"
	"song-lib/config"
	"song-lib/db"
)

type module struct {
	SongService song.Service
}

func NewModule(ctx context.Context, db db.Database, env config.Config) *module {
	// Create dependency tree here

	// Stores
	songStore := song.NewStore(db.GetPool())
	groupStore := group.NewStore(db.GetPool())

	// Services
	songService := song.NewService(songStore, groupStore, env)

	return &module{
		SongService: songService,
	}
}

func (m *module) Controllers() []common.Controller {
	return []common.Controller{
		song.NewController(m.SongService),
	}
}
