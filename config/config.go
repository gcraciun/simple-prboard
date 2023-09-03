package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Collector struct {
	Owner    string `yaml:"owner"`
	Interval int    `yaml:"interval"`
	Refresh  bool   `yaml:"refresh"`
}

type Config struct {
	Collector Collector
	Repos     map[string][]string `yaml:"repos"`
	Token     string
}

func LoadConfig() (*Config, error) {
	viper.SetConfigFile("./conf/config.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("fatal error reading config file: %s", err)
	}

	viper.AutomaticEnv()
	viper.BindEnv("token", "GIT_TOKEN")

	// retrieve the owner from config
	var collectorOwner string
	if !viper.IsSet("collector.owner") {
		return nil, fmt.Errorf("could not identify Github owner from the config file")
	} else {
		collectorOwnerRaw := viper.Get("collector.owner")
		ownerValue, ok := collectorOwnerRaw.(string)
		if !ok {
			return nil, fmt.Errorf("invalied value for collector.owner. Expected a string, got %v", collectorOwnerRaw)
		}
		collectorOwner = ownerValue
	}

	// set our default for collector.interval if the value isn't set in the config file
	// else if the value we get from the config isn't an int, then error out
	// otherwise we got our collector.interval ok from the config file
	defaultCollectorInterval := 60
	var collectorInterval int
	if !viper.IsSet("collector.interval") {
		collectorInterval = defaultCollectorInterval
	} else {
		collectorIntervalRaw := viper.Get("collector.interval")
		intervalValue, ok := collectorIntervalRaw.(int)
		if !ok {
			return nil, fmt.Errorf("invalid value for collector.interval. Expected an integer, got %v", collectorIntervalRaw)
		}
		collectorInterval = intervalValue
	}

	// set our default for collector.refresh if the value isn't set in the config file
	// else if the value we get from the config isn't a book, then error out
	// otherwise we got our collector.refresh ok from the config file
	defaultCollectorRefresh := true // or false, depending on your preferred default
	var collectorRefresh bool

	if !viper.IsSet("collector.refresh") {
		collectorRefresh = defaultCollectorRefresh
	} else {
		collectorRefreshRaw := viper.Get("collector.refresh")
		refreshValue, ok := collectorRefreshRaw.(bool)
		if !ok {
			return nil, fmt.Errorf("invalid value for collector.refresh. Expected a boolean (true/false), got %v", collectorRefreshRaw)
		}
		collectorRefresh = refreshValue
	}

	// get the repos, not error checking this for now
	repos := viper.GetStringMapStringSlice("repos")
	// get the token, not error checking this for now
	token := viper.GetString("token")

	cfg := &Config{
		Collector: Collector{
			Owner:    collectorOwner,
			Interval: collectorInterval,
			Refresh:  collectorRefresh,
		},
		Repos: repos,
		Token: token,
	}
	return cfg, nil
}
