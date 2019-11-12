module gitlab.in2p3.fr/rferrand/collectd-prometheus-plugin

go 1.12

require (
	collectd.org v0.3.0
	github.com/Tigraine/go-timemilli v0.0.0-20181122103244-0ec4e3ceddb1
	github.com/pkg/errors v0.8.1
	github.com/prometheus/client_model v0.0.0-20190812154241-14fe0d1b01d4
	github.com/prometheus/common v0.7.0
	github.com/spf13/viper v1.5.0
	golang.org/x/crypto v0.0.0-20190308221718-c2843e01d9a2
)

replace collectd.org => gitlab.in2p3.fr/rferrand/go-collectd v0.3.1-0.20191112182549-47d192eee5d3
