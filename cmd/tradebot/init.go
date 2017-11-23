package main

import (
	"runtime"

	"github.com/jessevdk/go-flags"
)

var opts struct {
	Exchanges []string `long:"exchanges" default:"binance" default:"zb" description:"Supported exchanges"`
	Conf      string   `long:"conf" description:"Config file"`
	LogLevel  string   `long:"log-level" default:"info" description:"Adjust the log level. Valid options are: error, warn, info, debug"`
}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

var parser = flags.NewParser(&opts, flags.Default|flags.IgnoreUnknown)
