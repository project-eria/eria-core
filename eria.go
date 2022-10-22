package eria

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
	"time"

	configmanager "github.com/project-eria/eria-core/config-manager"
	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
)

var (
	// Version is a placeholder that will receive the git tag version during build time
	// go build -v -ldflags "-X github.com/project-eria/eria-core.AppVersion=vx.x.x
	AppVersion  = "-"
	CoreVersion = "-"
	_logLevel   = zerolog.InfoLevel
	_configPath *string
	_appName    string
)

// Init gets the app name and version and displays app version if requested
func Init(appName string) {
	_appName = appName
	//	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	showVersion := flag.Bool("v", false, "Display the version")
	logLevelStr := flag.String("log", "info", "log level [error, warn, info, debug, trace]")
	_configPath = flag.String("config", "config.yml", "config file path")
	flag.Parse()

	// Show version (-v)
	if *showVersion {
		fmt.Printf("%s\n", AppVersion)
		os.Exit(0)
	}

	zlog.Logger = zlog.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "02/01|15:04:05"})
	zerolog.TimestampFunc = func() time.Time {
		return time.Now().In(time.Local)
	}

	zlog.Info().Msgf("[eria:Init] Starting %s %s...", _appName, AppVersion)

	logLevel, err := zerolog.ParseLevel(*logLevelStr)
	if err == nil {
		_logLevel = logLevel
		zerolog.SetGlobalLevel(logLevel)
	}
	zlog.Info().Stringer("log level", _logLevel).Msg("[eria:Init] Set log level")

	// Get EriaCore version
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		zlog.Error().Msg("[core:NewServer] Getting build info failed (not in module mode?)!")
		return
	}

	for _, dep := range bi.Deps {
		if dep.Path == "github.com/project-eria/eria-core" {
			CoreVersion = dep.Version
		}
	}
}

// LoadConfig Loads the config file into a struct
func LoadConfig(config interface{}) *configmanager.ConfigManager {
	cm, err := configmanager.Init(*_configPath, config)
	if err != nil {
		if configmanager.IsFileMissing(err) {
			zlog.Fatal().Msg("[eria:loadconfig] Config file do not exists...")
		} else {
			zlog.Fatal().Str("filePath", *_configPath).Err(err).Msg("[eria:loadconfig]")
		}
	}

	if err := cm.Load(); err != nil {
		zlog.Fatal().Err(err).Msg("[eria:loadconfig]")
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
	zlog.Info().Msg("[eria:WaitForSignal] Keyboard interrupt received, Stopping...")
}
