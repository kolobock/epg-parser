package main

import (
	"encoding/xml"
	"fmt"
	"os"
	"slices"
	"strings"
)

// Channel struct has <channel id="channel_id"> and <display-name> elements.
// <icon> element is optional
type Channel struct {
	XMLName     xml.Name `xml:"channel"`
	ID          string   `xml:"id,attr"`
	DisplayName []string `xml:"display-name,omitempty"`
	Icon        *Icon    `xml:"icon,omitempty"`
}

// IsInTheChannelsList checks if the channel ID is in the list of channels to keep, or any of channel DisplayName
func (c *Channel) IsInTheChannelsList(channelsToKeep []string) int {
	for i, keep := range channelsToKeep {
		if keep == "" {
			continue
		} else if keep == c.ID {
			return i
		}
		for _, displayName := range c.DisplayName {
			if keep == displayName {
				return i
			}
		}
	}

	return -1
}

func processChannel(channel *Channel, channelsToKeep []string, foundChannels map[string]bool, keepIcon bool, encoder *xml.Encoder) error {
	if idx := channel.IsInTheChannelsList(channelsToKeep); idx != -1 {
		if !keepIcon {
			channel.Icon = nil
		}
		if err := encoder.Encode(channel); err != nil {
			return fmt.Errorf("error encoding channel element: %v", err)
		}

		foundChannels[channel.ID] = true
		_ = slices.Delete(channelsToKeep, idx, idx+1)
	}

	return nil
}

func readChannelsFromFile(filePath string) ([]string, error) {
	channelsToKeepStr, err := os.ReadFile(filePath)
	if err != nil {
		return []string{}, fmt.Errorf("error reading channels file: %v", err)
	}

	channelsToKeep := strings.Split(string(channelsToKeepStr), "\n")
	for i, channel := range channelsToKeep {
		channelsToKeep[i] = strings.TrimSpace(channel)
	}
	channelsToKeep = slices.DeleteFunc(channelsToKeep, func(s string) bool {
		return s == ""
	})

	return channelsToKeep, nil
}

// Write not-found channels to the file
func logNotFoundChannels(channelsToKeep []string, config Config) error {
	notFoundFile, err := os.Create(fmt.Sprintf("%s-not-found.txt", config.ChannelsFile))
	if err != nil {
		return fmt.Errorf("error creating not-found %s file: %v", config.ChannelsFile, err)
	}
	defer notFoundFile.Close()

	if err := notFoundFile.Truncate(0); err != nil {
		return fmt.Errorf("error truncating not-found %s file: %v", config.ChannelsFile, err)
	}

	_, err = notFoundFile.WriteString(fmt.Sprintf("Channels not found in the XML file %s:\n", config.InputFile))
	if err != nil {
		return fmt.Errorf("error writing to not-found %s file: %v", config.ChannelsFile, err)
	}

	if len(channelsToKeep) > 0 {
		for _, channel := range channelsToKeep {
			if channel != "" {
				_, err = notFoundFile.WriteString(channel + "\n")
				if err != nil {
					return fmt.Errorf("error writing to not-found %s file: %v", config.ChannelsFile, err)
				}
			}
		}
	}

	return nil
}
