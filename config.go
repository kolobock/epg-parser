package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
)

// Config struct to hold the configuration for the program
type Config struct {
	InputFile           string `json:"input_file"`
	OutputFile          string `json:"output_file"`
	ChannelsFile        string `json:"channels_file"`
	KeepIcon            bool   `json:"keep_icon"`
	LogNotFoundChannels bool   `json:"log_not_found_channels"`
}

// Read the config file
func readConfig(filename string) (Config, error) {
	var config Config
	file, err := os.Open(filename)
	if err != nil {
		return config, fmt.Errorf("error opening config file: %v", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return config, fmt.Errorf("error decoding config file: %v", err)
	}

	return config, nil
}

// initConfig reads the configuration from the specified file and returns a Config struct.
func initConfig(filename string) Config {
	config, err := readConfig(filename)
	if err != nil {
		log.Printf("error reading config file: %v", err)
	}

	flag.StringVar(&config.InputFile, "input", config.InputFile, "Path to the input XML file")
	flag.StringVar(&config.OutputFile, "output", config.OutputFile, "Path to the output XML file")
	flag.StringVar(&config.ChannelsFile, "channels", config.ChannelsFile, "Path to the channels file to keep")
	flag.BoolVar(&config.KeepIcon, "keep-icon", config.KeepIcon, "Whether to keep icon elements in the output XML file")
	flag.BoolVar(&config.LogNotFoundChannels, "log-not-found-channels", config.LogNotFoundChannels, "Whether to log channels that were not found in the XML file to not-found-channels.txt")
	flag.Parse()

	return config
}
