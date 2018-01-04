package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/machinae/portfwd"
	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	verbose bool

	quiet bool

	configPath string

	timeout time.Duration
)

func init() {
	flag.BoolVarP(&verbose, "verbose", "v", false, "Verbose output")
	flag.BoolVarP(&quiet, "quiet", "q", false, "Suppress output")
	flag.DurationVarP(&timeout, "timeout", "t", 30*time.Second, "Dial timeout per connection")
	flag.StringVarP(&configPath, "config", "c", "portfwd.cfg", "Path to config file")
}

func main() {
	if err := parseFlags(); err != nil {
		log.Fatal(err)
	}
	if err := parseConfig(); err != nil {
		log.Fatal(err)
	}

	portfwd.Start()

	chSig := make(chan os.Signal, 1)
	signal.Notify(chSig, syscall.SIGINT, syscall.SIGTERM)

	// Wait for stop signal
	<-chSig
	portfwd.Stop()
}

func parseFlags() error {
	flag.Parse()
	if verbose {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	// Quiet overrides verbose
	if quiet {
		log.SetLevel(log.FatalLevel)
	}

	portfwd.Timeout = timeout

	viper.SetConfigType("toml")
	viper.SetConfigFile(configPath)

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	return nil
}
