package bailian

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type API struct {
	c      *Client
	apiKey string
	appId  string
}

func (api *API) WithApiKey(apiKey string) *API {
	api.apiKey = apiKey
	return api
}

func (api *API) WithAppId(appId string) *API {
	api.appId = appId
	return api
}

func (api *API) getApiKey() string {
	if api.apiKey != "" {
		return api.apiKey
	}
	return api.c.getApiKey()
}

func (api *API) getAppId() string {
	if api.appId != "" {
		return api.appId
	}
	return api.c.getAppId()
}

func (api *API) createBaseRequest(ctx context.Context, method, apiUrl string, body interface{}) (*http.Request, error) {
	var b io.Reader
	if body != nil {
		reqBytes, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		b = bytes.NewBuffer(reqBytes)
	} else {
		b = http.NoBody
	}

	fullUrl := fmt.Sprintf("%s/api/v1/apps/%s/completion", api.c.getBaseURL(), api.getAppId())
	req, err := http.NewRequestWithContext(ctx, method, fullUrl, b)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+api.getApiKey())
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("X-DashScope-SSE", "enable")
	return req, nil
}
