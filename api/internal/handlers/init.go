package handlers

import (
	"net/http"
	"time"
)

var Chttp *http.Client

func init() {
	Chttp = &http.Client{Timeout: 5 * time.Second}
}
