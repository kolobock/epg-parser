package main

import "encoding/xml"

// Programme struct has <programme start="start_time" stop="stop_time" channel="channel_id"> attribute,
// and contains <title> and <desc> elements as children. The <desc> element is optional.
type Programme struct {
	XMLName     xml.Name `xml:"programme"`
	Start       string   `xml:"start,attr"`
	Stop        string   `xml:"stop,attr"`
	Channel     string   `xml:"channel,attr"`
	Title       string   `xml:"title"`
	Description string   `xml:"desc,omitempty"`
}
