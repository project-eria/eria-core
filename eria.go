package eria

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	configmanager "github.com/project-eria/eria-core/config-manager"
	"github.com/project-eria/go-wot/consumer"
	"github.com/project-eria/go-wot/producer"
	"github.com/project-eria/go-wot/protocolHttp"
	"github.com/project-eria/go-wot/protocolWebSocket"
	"github.com/project-eria/go-wot/thing"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	_showVersion *bool
	_logLevel    *string
	_configPath  *string
	_version     string
	_appName     string
	_cancel      context.CancelFunc
	_ctx         context.Context
	_wait        sync.WaitGroup
)

// Init gets the app name and version and displays app version if requested
func Init(appName string, version string) {
	_ctx, _cancel = context.WithCancel(context.Background())
	_version = version
	_appName = appName
	//	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	_showVersion = flag.Bool("v", false, "Display the version")
	_logLevel = flag.String("log", "info", "log level [error, warn, info, debug, trace]")
	_configPath = flag.String("config", "config.yml", "config file path")
	flag.Parse()
	// Show version (-v)
	if *_showVersion {
		fmt.Printf("%s\n", version)
		os.Exit(0)
	}

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "02/01|15:04:05"})
	zerolog.TimestampFunc = func() time.Time {
		return time.Now().In(time.Local)
	}
	log.Info().Msgf("[eria:Init] Starting %s %s...", appName, version)

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
	_cancel()
	// Wait for the child goroutine to finish, which will only occur when
	// the child process has stopped and the call to cmd.Wait has returned.
	// This prevents main() exiting prematurely.
	_wait.Wait()
}

// RunSingleThing lauch a server for HTTP and WS requests
func RunSingleThing(t *thing.Thing, port uint) *producer.ExposedThing {
	//	_wait.Add(1)
	myProducer := producer.New(&_wait)
	exposedThing := myProducer.Produce(t)
	httpServer := protocolHttp.NewServer("127.0.0.1", port)
	myProducer.AddServer(httpServer)
	wsServer := protocolWebSocket.NewServer(httpServer)
	myProducer.AddServer(wsServer)

	// for key, property := range t.Properties {
	// 	property := property // Copy https://go.dev/doc/faq#closures_and_goroutines
	// 	handler := func() (interface{}, error) {
	// 		return property.GetValue(), nil
	// 	}
	// 	exposedThing.SetPropertyReadHandler(key, handler)
	// }
	myProducer.Expose()
	go func() {
		<-_ctx.Done()
		myProducer.Stop()
		//		_wait.Done()
	}()
	return exposedThing
}

// ConsumeThing connect a remote thing WS server
func ConsumeThing(url string) (*consumer.ConsumedThing, error) {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal().Err(err).Msg("[eria:ConsumeThing]")
		return nil, errors.New("can't open thing url")
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Fatal().Str("status", resp.Status).Str("url", url).Msg("[eria:ConsumeThing] incorrect response")
		return nil, errors.New("incorrect response")
	}

	var td thing.Thing
	if err := json.NewDecoder(resp.Body).Decode(&td); err != nil {
		log.Fatal().Str("url", url).Err(err).Msg("[eria:ConsumeThing]")
		return nil, errors.New("can't decode json")
	}

	myConsumer := consumer.New()
	httpClient := protocolHttp.NewClient()
	myConsumer.AddClient(httpClient)
	wsClient := protocolWebSocket.NewClient()
	myConsumer.AddClient(wsClient)
	consumedThing := myConsumer.Consume(&td)

	if err != nil {
		return nil, err
	}
	return consumedThing, nil
}
