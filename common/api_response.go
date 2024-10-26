package common

type apiResponse struct {
	Message string `json:"message"`
}

func NewApiResponse(message string) *apiResponse {
	return &apiResponse{
		Message: message,
	}
}
