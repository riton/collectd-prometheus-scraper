package scraper

import (
	"net/http"
	"time"
)

type httpDoer interface {
	Do(*http.Request) (*http.Response, error)
}

type HTTPBasicCreds struct {
	User     string
	Password string
}

type httpClient struct {
	Credentials *HTTPBasicCreds
	client      *http.Client
}

func newHTTPClient(timeout time.Duration, creds *HTTPBasicCreds) *httpClient {
	hClient := http.Client{
		Timeout: timeout,
	}

	return &httpClient{
		Credentials: creds,
		client:      &hClient,
	}
}

func (hc httpClient) Do(req *http.Request) (*http.Response, error) {
	if hc.Credentials != nil {
		req.SetBasicAuth(hc.Credentials.User, hc.Credentials.Password)
	}
	return hc.client.Do(req)
}
