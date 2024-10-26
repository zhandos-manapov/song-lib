package common

import (
	"github.com/go-chi/chi/v5"
)

type BaseController interface {
	Path() string
}

type Controller interface {
	BaseController
	RegisterRoutes(r *chi.Router)
}

type baseController struct {
	basePath string
}

func NewBaseController(basePath string) BaseController {
	return &baseController{
		basePath: basePath,
	}
}

func (c *baseController) Path() string {
	return c.basePath
}
