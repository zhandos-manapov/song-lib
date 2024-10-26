package common

type BaseService interface {
}

type baseService struct {
}

func NewBaseService() BaseService {
	return &baseService{}
}
