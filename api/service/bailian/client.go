package bailian

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type Client struct {
	baseURL    string
	apiKey     string
	appId      string
	httpClient *http.Client
}

func NewClientWithConfig(c *ClientConfig) *Client {
	var httpClient = &http.Client{}

	if c.Timeout != 0 {
		httpClient.Timeout = c.Timeout
	}
	if c.Transport != nil {
		httpClient.Transport = c.Transport
	}

	baseURL := c.BaseURL
	if baseURL == "" {
		baseURL = DefaultBaseURL
	}

	return &Client{
		baseURL:    baseURL,
		apiKey:     c.ApiKey,
		appId:      c.AppId,
		httpClient: httpClient,
	}
}

func NewClient(apiKey, appId string) *Client {
	return NewClientWithConfig(&ClientConfig{
		ApiKey: apiKey,
		AppId:  appId,
	})
}

func (c *Client) sendRequest(req *http.Request) (*http.Response, error) {
	return c.httpClient.Do(req)
}

func (c *Client) sendJSONRequest(req *http.Request, res interface{}) error {
	resp, err := c.sendRequest(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errBody struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		}
		err = json.NewDecoder(resp.Body).Decode(&errBody)
		if err != nil {
			return err
		}
		return fmt.Errorf("HTTP response error: [%v]%v", errBody.Code, errBody.Message)
	}

	err = json.NewDecoder(resp.Body).Decode(res)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) getBaseURL() string {
	return strings.TrimSuffix(c.baseURL, "/")
}

func (c *Client) getApiKey() string {
	return c.apiKey
}

func (c *Client) getAppId() string {
	return c.appId
}

func (c *Client) API() *API {
	return &API{
		c: c,
	}
}
