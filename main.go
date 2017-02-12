package main

import (
	"os"

	log "github.com/Sirupsen/logrus"
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
	wlc.SetOutputResolutionCb(comp.OutputResolution)
	wlc.SetViewCreatedCb(comp.ViewCreated)
	wlc.SetViewRequestGeometryCb(comp.ViewRequestGeometry)
	wlc.SetPointerMotionCb(comp.PointerMotion)
	wlc.SetKeyboardKeyCb(comp.KeyboardKey)
	wlc.LogSetHandler(wlcLogHandler)

	if !wlc.Init() {
		os.Exit(1)
	}

	wlc.Run()
	os.Exit(0)
}

func wlcLogHandler(typ wlc.LogType, msg string) {
	format := "[WLC] %s"
	switch typ {
	case wlc.LogInfo:
		log.Debugf(format, msg)
	case wlc.LogWarn:
		log.Warnf(format, msg)
	case wlc.LogError:
		log.Errorf(format, msg)
	case wlc.LogWayland:
		log.Debugf("[WLC - Wayland] %s", msg)
	}
}
