package eria

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	configmanager "github.com/project-eria/eria-core/config-manager"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	_showVersion *bool
	_logLevel    *string
	_configPath  *string
	_version     string
	_appName     string
)

// Init gets the app name and version and displays app version if requested
func Init(appName string, version string) {
	_version = version
	_appName = appName
	//	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	_showVersion = flag.Bool("v", false, "Display the version")
	_logLevel = flag.String("log", "info", "log level [error, warn, info, debug, trace]")
	_configPath = flag.String("config", "config.yml", "config file path")
	flag.Parse()
	// Show version (-v)
	if *_showVersion {
		fmt.Printf("%s\n", _version)
		os.Exit(0)
	}

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "02/01|15:04:05"})
	zerolog.TimestampFunc = func() time.Time {
		return time.Now().In(time.Local)
	}
	log.Info().Msgf("[eria:Init] Starting %s %s...", _appName, _version)

	level, err := zerolog.ParseLevel(*_logLevel)
	if err != nil {
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)
	log.Log().Stringer("level", level).Msg("[eria:Init] Set log level")
}

// LoadConfig Loads the config file into a struct
func LoadConfig(config interface{}) *configmanager.ConfigManager {
	cm, err := configmanager.Init(*_configPath, config)
	if err != nil {
		if configmanager.IsFileMissing(err) {
			log.Fatal().Msg("[eria:loadconfig] Config file do not exists...")
		} else {
			log.Fatal().Str("filePath", *_configPath).Err(err).Msg("[eria:loadconfig]")
		}
	}

	if err := cm.Load(); err != nil {
		log.Fatal().Err(err).Msg("[eria:loadconfig]")
	}
	return cm
}

// WaitForSignal Wait for any signal and runs all the defer
func WaitForSignal() {
	// Set up channel on which to send signal notifications.
	// We must use a buffered channel or risk missing the signal
	// if we're not ready to receive when the signal is sent.
	c := make(chan os.Signal, 1)
	signal.Notify(c,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	// Block until keyboard interrupt is received.
	<-c
	log.Info().Msg("[eria:WaitForSignal] Keyboard interrupt received, Stopping...")
}
