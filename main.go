// This program filters an XML TV file based on a list of channels to keep.
// It reads the input XML file, processes the channels and programmes, and writes
// the filtered output to a new XML file. It also logs the channels that were
// not found in the input XML file if configured to do so.
package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync/atomic"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Starting XML TV filter...")

	config := initConfig("config.json")
	log.Printf("\tConfiguration loaded: %+v\n", config)

	// Parse the XML file
	if err := parseXML(config); err != nil {
		log.Fatalf("Error parsing XML file: %v", err)
	}

	log.Printf("Successfully parsed XML file and created output file at %s\n", config.OutputFile)
}

// Parse the XML file and keep channel and programme for ChannelsFile only list
func parseXML(config Config) error {
	channelsToKeep, err := readChannelsFromFile(config.ChannelsFile)
	if err != nil {
		return err
	}

	if len(channelsToKeep) < 1 {
		return fmt.Errorf("No channels to keep defined at %s. exiting", config.ChannelsFile)
	}

	log.Printf("\tRead channels to keep, count: %d\n", len(channelsToKeep))

	outputFile, err := os.Create(config.OutputFile)
	if err != nil {
		return fmt.Errorf("error creating output XML file: %v", err)
	}
	defer outputFile.Close()

	inputFile, err := os.Open(config.InputFile)
	if err != nil {
		return fmt.Errorf("error opening input XML file: %v", err)
	}
	defer inputFile.Close()

	// Create a map to keep track of channels that were found in the XML file
	foundChannels := make(map[string]bool)

	// Create a new XML encoder for the output file
	encoder := xml.NewEncoder(outputFile)
	encoder.Indent("", "  ")

	// Create a new XML decoder for the input file
	decoder := xml.NewDecoder(inputFile)

	log.Printf("\tStart parsing XML file %s\n", config.InputFile)

	var tokensCounter atomic.Int64

	for {
		tok, err := decoder.Token()
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("error reading XML token: %v", err)
		}

		tokensCounter.Add(1)
		if tokensCounter.Load()%10000 == 0 {
			os.Stdout.Write([]byte(fmt.Sprintf("\tProcessed %d XML tokens\r", tokensCounter.Load())))
		}

		switch se := tok.(type) {
		case xml.StartElement:
			switch se.Name.Local {
			case "channel":
				var channel Channel
				if err := decoder.DecodeElement(&channel, &se); err != nil {
					log.Printf("\t\terror decoding channel element: %v\n", err)

					continue
				}

				if err := processChannel(&channel, channelsToKeep, foundChannels, config.KeepIcon, encoder); err != nil {
					return err
				}
			case "programme":
				var programme Programme
				if err := decoder.DecodeElement(&programme, &se); err != nil {
					log.Printf("\t\terror decoding programme element: %v\n", err)

					continue
				}

				if foundChannels[programme.Channel] {
					if err := encoder.Encode(programme); err != nil {
						return fmt.Errorf("error encoding programme element: %v", err)
					}
				}
			default:
				// For other elements, just copy them to the output file
				if err := encoder.EncodeToken(se); err != nil {
					return fmt.Errorf("error encoding start element: %v", err)
				}
			}
		case xml.EndElement:
			if err := encoder.EncodeToken(se); err != nil {
				return fmt.Errorf("error encoding end element: %v", err)
			}

			if se.Name.Local == "tv" {
				break
			}
		case xml.CharData:
			// For character data, just copy it to the output file except for whitespace-only character data.
			if strings.TrimSpace(string(se)) == "" {
				continue
			}

			if err := encoder.EncodeToken(se); err != nil {
				return fmt.Errorf("error encoding char data: %v", err)
			}
		default:
			if err := encoder.EncodeToken(tok); err != nil {
				return fmt.Errorf("error encoding default token: %v", err)
			}
		}
	}

	if config.LogNotFoundChannels {
		if err := logNotFoundChannels(channelsToKeep, config); err != nil {
			return err
		}
	}

	// Flush the encoder to write any remaining tokens to the output file
	if err := encoder.Flush(); err != nil {
		return fmt.Errorf("error flushing encoder: %v", err)
	}

	return nil
}
