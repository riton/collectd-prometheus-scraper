package scraper

import (
	"time"

	"gitlab.in2p3.fr/rferrand/collectd-prometheus-plugin/transport"
)

var newHTTPClientFnc = func(timeout time.Duration,
	creds transport.HTTPBasicCreds) transport.HTTPDoer {
	return transport.NewHTTPClient(timeout, creds)
}
