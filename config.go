package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/ccin2p3/collectd-prometheus-plugin/transport"
)

// All this configuration stuff should be integrated in
// collectd complex configuration callback as soon as the
// Go bindings are available

const (
	configEnvVarPrefix     = "COLLECTD_PROM"
	configBaseDir          = "/etc/"
	configFilename         = "collectd-prometheus_scraper"
	defaultScrapeTimeout   = 5 * time.Second
	defaultMetricsPath     = "/metrics"
	defaultHonorTimestamps = false
	defaultTransportScheme = "http"
	defaultPort            = 80
)

type scrapeCollectdConfig struct {
	PluginName                 string `mapstructure:"plugin_name"`
	MetadataPrefix             string `mapstructure:"metadata_prefix"`
	TypeInstanceOnlyHashedMeta bool   `mapstructure:"type_instance_only_hashed_meta"`
	HashLabelFunctionHashSize  int    `mapstructure:"hash_function_size"`
}

// scrapeConfig is heavily inspired by
// https://prometheus.io/docs/prometheus/latest/configuration/configuration/#scrape_config
type scraperConfig struct {
	JobName         string                   `mapstructure:"job_name"`
	Target          string                   `mapstructure:"target"`
	HonorTimestamps bool                     `mapstructure:"honor_timestamps"`
	ScrapeTimeout   time.Duration            `mapstructure:"scrape_timeout"`
	MetricsPath     string                   `mapstructure:"metrics_path"`
	Labels          map[string]interface{}   `mapstructure:"labels"`
	BasicAuth       transport.HTTPBasicCreds `mapstructure:"basic_auth"`
	Collectd        scrapeCollectdConfig     `mapstructure:"collectd"`
	Scheme          string                   `mapstructure:"scheme"`
	Port            int                      `mapstructure:"port"`
}

type config struct {
	Debug         bool            `mapstructure:"debug"`
	ScrapeTimeout time.Duration   `mapstructure:"scrape_timeout"`
	ScrapeConfigs []scraperConfig `mapstructure:"scrape_configs"`
}

func newConfig() (config, error) {
	var cfg config

	vcfg, err := newViperConfig()
	if err != nil {
		return cfg, errors.Wrap(err, "initializing viper")
	}

	cfg.Debug = vcfg.GetBool("debug")
	cfg.ScrapeTimeout = vcfg.GetDuration("scrape_timeout")

	if !vcfg.IsSet("scrape_configs") {
		return cfg, errors.New("missing scrape_configs section")
	}

	for scrapeKey := range vcfg.GetStringMap("scrape_configs") {

		sub := vcfg.Sub("scrape_configs." + scrapeKey)
		sub.SetDefault("scrape_timeout", cfg.ScrapeTimeout)
		sub.SetDefault("metrics_path", defaultMetricsPath)
		sub.SetDefault("honor_timestamps", defaultHonorTimestamps)
		sub.SetDefault("collectd.plugin_name", scrapeKey)
		sub.SetDefault("collectd.metadata_prefix", scrapeKey+".")
		sub.SetDefault("collectd.type_instance_only_hashed_meta", true)
		sub.SetDefault("collectd.hash_function_size", 32)
		sub.SetDefault("scheme", defaultTransportScheme)
		sub.SetDefault("port", defaultPort)

		var sconfig scraperConfig

		sconfig.JobName = scrapeKey
		if err := sub.Unmarshal(&sconfig); err != nil {
			return cfg, errors.Wrap(err, fmt.Sprintf("unmarsharling configuration for %s",
				scrapeKey))
		}

		cfg.ScrapeConfigs = append(cfg.ScrapeConfigs, sconfig)
	}

	return cfg, nil
}

func newViperConfig() (*viper.Viper, error) {
	v := viper.New()

	v.AddConfigPath(configBaseDir)
	v.SetConfigName(configFilename)

	v.SetEnvPrefix(configEnvVarPrefix)
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	v.AutomaticEnv()

	v.SetDefault("debug", false)
	v.SetDefault("scrape_timeout", defaultScrapeTimeout)

	err := v.ReadInConfig()

	return v, err
}
