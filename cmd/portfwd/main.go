package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/machinae/portfwd"
	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	verbose bool

	configPath string
)

func init() {
	flag.BoolVarP(&verbose, "verbose", "v", false, "Verbose output")
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

	viper.SetConfigType("toml")
	viper.SetConfigFile(configPath)

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	return nil
}
