// Package config parses command-line/environment/config file arguments
// and make available to other packages.
package config

import (
	"runtime"
	"strings"
	"time"

	"github.com/Akagi201/utilgo/conflag"
	flags "github.com/jessevdk/go-flags"
	log "github.com/sirupsen/logrus"
)

// Opts configs
var Opts struct {
	Conf         string `long:"conf" description:"esalert config file"`
	AlertFileDir string `long:"alerts" short:"a" required:"true" description:"A yaml file, or directory with yaml files, containing alert definitions"`
	LuaInit      string `long:"lua-init" description:"If set the given lua script file will be executed at the initialization of every lua vm"`
	LuaVMs       int    `long:"lua-vms" default:"1" description:"How many lua vms should be used. Each vm is completely independent of the other, and requests are executed on whatever vm is available at that moment. Allows lua scripts to not all be blocked on the same os thread"`
	SlackWebhook string `long:"slack-webhook" description:"Slack webhook url, required if using any Slack actions"`
	ForceRun     string `long:"force-run" description:"If set with the name of an alert, will immediately run that alert and exit. Useful for testing changes to alert definitions"`
	LogLevel     string `long:"log-level" default:"info" description:"Adjust the log level. Valid options are: error, warn, info, debug"`
}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func init() {
	parser := flags.NewParser(&Opts, flags.Default|flags.IgnoreUnknown)

	parser.Parse()

	if Opts.Conf != "" {
		conflag.LongHyphen = true
		conflag.BoolValue = false
		args, err := conflag.ArgsFrom(Opts.Conf)
		if err != nil {
			panic(err)
		}

		parser.ParseArgs(args)
	}

	log.Infof("alert opts: %+v", Opts)
}

func init() {
	if level, err := log.ParseLevel(strings.ToLower(Opts.LogLevel)); err != nil {
		log.SetLevel(level)
	}

	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: time.RFC3339,
	})
}
