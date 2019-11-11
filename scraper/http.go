package scraper

import (
	"time"

	"github.com/ccin2p3/collectd-prometheus-plugin/transport"
)

var newHTTPClientFnc = func(timeout time.Duration,
	creds transport.HTTPBasicCreds) transport.HTTPDoer {
	return transport.NewHTTPClient(timeout, creds)
}
