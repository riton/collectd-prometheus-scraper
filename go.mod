module gitlab.in2p3.fr/rferrand/collectd-prometheus-plugin

go 1.12

require (
	collectd.org v0.3.0
	github.com/Tigraine/go-timemilli v0.0.0-20181122103244-0ec4e3ceddb1
	github.com/prometheus/client_model v0.0.0-20190812154241-14fe0d1b01d4
	github.com/prometheus/common v0.7.0
	gitlab.in2p3.fr/rferrand/go-system-utils v0.1.0
)

replace collectd.org => gitlab.in2p3.fr/rferrand/go-collectd v0.3.1-0.20191107161421-7932bdd4c73c
