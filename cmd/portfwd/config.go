package main

import (
	"errors"
	"fmt"

	"github.com/machinae/portfwd"
	"github.com/spf13/viper"
)

// parse config created by viper
func parseConfig() error {
	var allKeys []string
	for k := range viper.AllSettings() {
		allKeys = append(allKeys, k)
	}
	if len(allKeys) == 0 {
		return errors.New("At least one forward section is required")
	}
	for _, k := range allKeys {
		err := parseForwardConfig(k)
		if err != nil {
			return err
		}
	}
	return nil
}

// Parses config and sets it globally on portfwd
func parseForwardConfig(key string) error {
	defaultStrategy := "random"
	fk := key + ".strategy"
	viper.SetDefault(fk, defaultStrategy)

	strategy := viper.GetString(fk)
	hostList, err := portfwd.HostListForStrategy(strategy)
	if err != nil {
		return err
	}

	fk = key + ".from"
	lHost := viper.GetString(fk)
	if lHost == "" {
		return fmt.Errorf("%s: 'from' is required", key)
	}

	fk = key + ".to"
	rHosts := viper.GetStringSlice(fk)
	if len(rHosts) == 0 {
		return fmt.Errorf("%s: 'to' is required", key)
	}

	for _, rHost := range rHosts {
		err := hostList.AddHost(rHost)
		if err != nil {
			return err
		}
	}

	portfwd.AddForwarder(lHost, hostList)

	return nil

}
