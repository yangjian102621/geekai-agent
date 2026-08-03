package dify

import (
	"net/http"
	"time"
)

type ClientConfig struct {
	Host             string
	ApiSecretKey     string // deprecated: use DefaultAPISecret instead
	DefaultAPISecret string
	Timeout          time.Duration
	Transport        *http.Transport
}
