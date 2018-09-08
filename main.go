package main

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/alecthomas/kingpin"
	"github.com/mikkeloscar/flis/backend"
	"github.com/mikkeloscar/flis/compositor"
	"github.com/mikkeloscar/flis/config"
	"github.com/mikkeloscar/flis/layout/i3"
	wlc "github.com/mikkeloscar/go-wlc"
)

var (
	version = "unknown"
	flags   struct {
		Config   string
		Debug    bool
		Validate bool
	}
)

func init() {
	kingpin.
		Flag("config", "Path to config file").Short('c').
		StringVar(&flags.Config)
	kingpin.
		Flag("debug", "Enable debug logging").Short('d').
		BoolVar(&flags.Debug)
	kingpin.
		Flag("validate", "Validate config file and exit").Short('C').
		BoolVar(&flags.Validate)
	kingpin.Version(version)
}

func main() {
	kingpin.Parse()

	// configure log level
	if flags.Debug {
		log.SetLevel(log.DebugLevel)
	}

	conf, err := config.LoadConfig(flags.Config)
	if err != nil {
		log.Errorf("Failed to parse config '%s': %s", flags.Config, err)
		os.Exit(1)
	}

	if flags.Validate {
		log.Info("Config is valid")
		os.Exit(0)
	}

	// initialize compositor
	comp := compositor.New(conf, &backend.WLC{}, i3.New())

	wlc.SetOutputCreatedCb(comp.OutputCreated)
	wlc.SetViewCreatedCb(comp.ViewCreated)
	wlc.SetPointerMotionCb(comp.PointerMotion)
	wlc.SetKeyboardKeyCb(comp.KeyboardKey)

	if !wlc.Init() {
		os.Exit(1)
	}

	wlc.Run()
	os.Exit(0)
}
