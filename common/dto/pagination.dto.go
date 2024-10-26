package coredto

type PaginationDto struct {
	Skip  int `validate:"number,min=0"`
	Limit int `validate:"number,min=0"`
}
