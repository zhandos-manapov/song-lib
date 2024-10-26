package main

import (
	"net/http"
	"song-lib/common"

	"github.com/go-chi/chi/v5"
)

type Router interface {
	LoadControllers(prefix string, controllers []common.Controller)
	Run(addr string) error
}

type router struct {
	mux *chi.Mux
}

func NewRouter() Router {
	return &router{
		mux: chi.NewRouter(),
	}
}

func (r *router) LoadControllers(prefix string, controllers []common.Controller) {
	r.mux.Route(prefix, func(r chi.Router) {
		for _, c := range controllers {
			r.Route(c.Path(), func(r chi.Router) {
				c.RegisterRoutes(&r)
			})
		}
	})
}

func (r *router) Run(addr string) error {
	return http.ListenAndServe(addr, r.mux)
}
