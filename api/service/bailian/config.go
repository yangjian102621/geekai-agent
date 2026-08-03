package bailian

import (
	"net/http"
	"time"
)

const DefaultBaseURL = "https://dashscope.aliyuncs.com"

type ClientConfig struct {
	BaseURL   string
	ApiKey    string
	AppId     string
	Timeout   time.Duration
	Transport *http.Transport
}
