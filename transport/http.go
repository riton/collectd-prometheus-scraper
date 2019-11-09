package transport

import (
	"net/http"
	"time"
)

type HTTPDoer interface {
	Do(*http.Request) (*http.Response, error)
}

type HTTPBasicCreds struct {
	User     string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

type httpClient struct {
	Credentials *HTTPBasicCreds
	client      *http.Client
}

func NewHTTPClient(timeout time.Duration, creds HTTPBasicCreds) *httpClient {
	hClient := http.Client{
		Timeout: timeout,
	}

	return &httpClient{
		Credentials: &creds,
		client:      &hClient,
	}
}

func (hc httpClient) Do(req *http.Request) (*http.Response, error) {
	if hc.Credentials != nil {
		req.SetBasicAuth(hc.Credentials.User, hc.Credentials.Password)
	}
	return hc.client.Do(req)
}
