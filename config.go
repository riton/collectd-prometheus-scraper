package main

import (
	"strconv"

	ros "gitlab.in2p3.fr/rferrand/go-system-utils/os"
)

type config struct {
	Debug bool
}

func newConfig() *config {
	debug, _ := strconv.ParseBool(ros.GetEnv("DEBUG", "false"))
	return &config{
		Debug: debug,
	}
}
