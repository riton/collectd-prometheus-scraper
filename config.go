package main

import (
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// All this configuration stuff should be integrated in
// collectd complex configuration callback as soon as the
// Go bindings are available

const (
	configEnvVarPrefix = "COLLECTD_PROM"
	configBaseDir      = "/etc"
	configFilename     = "collectd-prometheus_scraper"
)

type config struct {
	Debug bool `yaml:"debug"`
}

func newConfig() (config, error) {
	var cfg config

	vcfg, err := newViperConfig()
	if err != nil {
		return cfg, errors.Wrap(err, "initializing viper")
	}

	if err = vcfg.Unmarshal(&cfg); err != nil {
		return cfg, errors.Wrap(err, "unmarshaling configuration")
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

	err := v.ReadInConfig()

	return v, err
}
