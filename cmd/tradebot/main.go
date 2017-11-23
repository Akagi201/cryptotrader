package main

import (
	"strings"
	"time"

	"github.com/Akagi201/utilgo/conflag"
	log "github.com/sirupsen/logrus"
)

func main() {
	parser.Parse()

	if opts.Conf != "" {
		conflag.LongHyphen = true
		conflag.BoolValue = false
		args, err := conflag.ArgsFrom(opts.Conf)
		if err != nil {
			panic(err)
		}

		parser.ParseArgs(args)
	}

	log.Debugf("opts: %+v", opts)

	if opts.LogLevel == "" {
		return
	}

	level, err := log.ParseLevel(strings.ToLower(opts.LogLevel))
	if err != nil {
		log.Fatalf("log level error: %v", err)
	}

	log.SetLevel(level)

	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: time.RFC3339,
	})
}
