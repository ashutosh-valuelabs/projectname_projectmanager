package model

import "net/http"


// Mystruct : yaml configuration
type Mystruct struct {
	id func(http.ResponseWriter, *http.Request)
}

