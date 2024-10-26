package common

import "net/http"

type apiFunc func(http.ResponseWriter, *http.Request) error 
