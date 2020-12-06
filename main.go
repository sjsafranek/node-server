package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/sjsafranek/logger"
	"github.com/sjsafranek/node-server/config"
	"github.com/sjsafranek/node-server/server"
)

const (
	PROJECT             string = "node-server"
	VERSION             string = "0.0.1"
	DEFAULT_CONFIG_FILE string = "config.toml"
)

var (
	CONFIG_FILE string = DEFAULT_CONFIG_FILE
)

func init() {
	var print_version bool
	flag.StringVar(&CONFIG_FILE, "c", DEFAULT_CONFIG_FILE, "Config file")
	flag.BoolVar(&print_version, "V", false, "Print version and exit")
	flag.Parse()

	if print_version {
		fmt.Println(PROJECT, VERSION)
		os.Exit(0)
	}

	signal_queue := make(chan os.Signal)
	signal.Notify(signal_queue, syscall.SIGTERM)
	signal.Notify(signal_queue, syscall.SIGINT)
	go func() {
		sig := <-signal_queue
		logger.Warnf("caught sig: %+v", sig)
		logger.Warn("Gracefully shutting down...")
		logger.Warn("Shutting down...")
		time.Sleep(500 * time.Millisecond)
		os.Exit(0)
	}()
}

func main() {
	logger.Debugf("%v:%v", PROJECT, VERSION)
	hostname, _ := os.Hostname()
	logger.Debug("Hostname: ", hostname)
	logger.Debug("GOOS: ", runtime.GOOS)
	logger.Debug("CPUS: ", runtime.NumCPU())
	logger.Debug("PID: ", os.Getpid())
	logger.Debug("Go Version: ", runtime.Version())
	logger.Debug("Go Arch: ", runtime.GOARCH)
	logger.Debug("Go Compiler: ", runtime.Compiler)
	logger.Debug("NumGoroutine: ", runtime.NumGoroutine())

	conf, err := config.Open(CONFIG_FILE)
	if nil != err {
		logger.Warn(err)
		logger.Info("Using default config settings")
		conf = config.New()
		conf.Save(CONFIG_FILE)
	}

	app, err := server.New(conf)
	if nil != err {
		panic(err)
	}

	err = app.ListenAndServe(":8080")
	if nil != err {
		panic(err)
	}
}
