package main

import "encoding/xml"

// Icon struct has <icon src="icon_url"> element, we need to keep the src attribute if keepIcon config is set to true
type Icon struct {
	XMLName xml.Name `xml:"icon"`
	Src     string   `xml:"src,attr"`
}
