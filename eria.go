package eria

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
	"time"

	"github.com/project-eria/go-wot/consumer"
	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
)

var (
	// Version is a placeholder that will receive the git tag version during build time
	// go build -v -ldflags "-X github.com/project-eria/eria-core.AppVersion=vx.x.x
	AppVersion      = "-"
	BuildDate       = "-"
	CoreVersion     = "-"
	_logLevel       = zerolog.InfoLevel
	_configPath     *string
	_logPath        *string
	_logFormat      *string
	_logOutput      *os.File
	_logNoColor     = false
	_appName        string
	_consumedThings map[string]consumer.ConsumedThing
	_location       *time.Location
)

// Init gets the app name and version and displays app version if requested
func Init(appName string, config interface{}) {
	_appName = appName
	//	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	showVersion := flag.Bool("v", false, "Display the version")
	logLevelStr := flag.String("log", "info", "log level [error, warn, info, debug, trace]")
	_configPath = flag.String("config-path", "config.yml", "config file path")
	_logPath = flag.String("log-path", "", "log file path")
	_logFormat = flag.String("log-format", "pretty", "log output format [pretty, json]")
	flag.Parse()
	var err error
	// Show version (-v)
	if *showVersion {
		fmt.Printf("%s (%s)\n", AppVersion, BuildDate)
		os.Exit(0)
	}

	if *_logPath != "" {
		_logOutput, err = os.OpenFile(*_logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			fmt.Printf("Can't set log file '%s': %v\n", *_logPath, err)
			os.Exit(1)
		}
		_logNoColor = true
		zlog.Trace().Str("file", *_logPath).Msg("[core:Init] Log on file selected")
	} else {
		_logOutput = os.Stdout
		zlog.Trace().Msg("[core:Init] Log on console selected")
	}

	if *_logFormat == "pretty" {
		zlog.Logger = zlog.Output(zerolog.ConsoleWriter{Out: _logOutput, NoColor: _logNoColor, TimeFormat: "02/01|15:04:05"})
	} else if *_logFormat == "json" {
		zlog.Logger = zlog.Output(_logOutput)
	} else {
		fmt.Printf("unknown log format '%s'\n", *_logFormat)
		os.Exit(1)
	}

	zerolog.TimestampFunc = func() time.Time {
		return time.Now().In(time.Local)
	}

	zlog.Info().Msgf("[core:Init] Starting %s %s...", _appName, AppVersion)

	logLevel, err := zerolog.ParseLevel(*logLevelStr)
	if err == nil {
		_logLevel = logLevel
		zerolog.SetGlobalLevel(logLevel)
	}
	zlog.Info().Stringer("log level", _logLevel).Msg("[core:Init] Set log level")

	// Get EriaCore version
	// Based on https://stackoverflow.com/questions/54890161/how-to-get-go-detailed-build-logs-with-all-used-packages-in-gopath-and-go-modu/54890460#54890460
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		zlog.Error().Msg("[core:Init] Getting build info failed (not in module mode?)!")
		return
	}

	for _, dep := range bi.Deps {
		if dep.Path == "github.com/project-eria/eria-core" {
			CoreVersion = dep.Version
		}
	}

	// Load the config file
	loadConfig(config)

	// Set the location
	_location, err = time.LoadLocation(eriaConfig.Location)
	if err != nil {
		zlog.Error().Err(err).Msg("[core:Init] Can't load location")
		return
	}
}

// Return the location
func Location() *time.Location {
	return _location
}

// Stop and close all services/files
func Close() {
	zlog.Debug().Msg("[core:Close] Closing...")
	_logOutput.Close()
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
	zlog.Info().Msg("[core:WaitForSignal] Keyboard interrupt received")
}

func Start(instance string) {
	ConnectThings()
	// Automations are present in the config file
	startAutomations(instance)
	startCronScheduler()
	Producer(instance).StartServer()
	WaitForSignal()
	Producer(instance).StopServer()
}
