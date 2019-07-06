//
// meteoclimatic.go
// Copyright (C) 2017 Paco Esteban <paco@onna.be>
//
// Distributed under terms of the MIT license.
//
// This is overcomplicated, I know. It's just a test
// to have a package that is simple to use and testable.
// I took some ideas from here and there ...

package meteoclimatic

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"golang.org/x/net/html/charset"
)

const (
	baseURL string = "https://www.meteoclimatic.net/feed/xml/"
)

type MeteoClimatic struct {
	XMLData []byte
	MD      Meteodata
}

type Meteodata struct {
	XMLName       xml.Name  `xml:"meteodata"`
	Description   string    `xml:"description"`
	Link          string    `xml:"link"`
	Copyright     string    `xml:"copyright"`
	Ttl           int       `xml:"ttl"`
	PubDate       string    `xml:"pubDate"`
	Docs          string    `xml:"docs"`
	MeteoStations MStations `xml:"stations"`
}

type MStations struct {
	PubDate string     `xml:"pubDate"`
	Station []MStation `xml:"station"`
}
type MStation struct {
	Id          string `xml:"id"`
	Location    string `xml:"location"`
	Homepage    string `xml:"homepage"`
	Datasheet   string `xml:"datasheet"`
	Author      string `xml:"author"`
	PubDate     string `xml:"pubDate"`
	QOS         string `xml:"qos"`
	StationData struct {
		Temperature struct {
			Unit string  `xml:"unit"`
			Now  float64 `xml:"now"`
			Max  float64 `xml:"max"`
			Min  float64 `xml:"min"`
		} `xml:"temperature"`
		Humidity struct {
			Unit string  `xml:"unit"`
			Now  float64 `xml:"now"`
			Max  float64 `xml:"max"`
			Min  float64 `xml:"min"`
		} `xml:"humidity"`
		Barometre struct {
			Unit string  `xml:"unit"`
			Now  float64 `xml:"now"`
			Max  float64 `xml:"max"`
			Min  float64 `xml:"min"`
		} `xml:"barometre"`
		Wind struct {
			Unit    string  `xml:"unit"`
			Now     float64 `xml:"now"`
			Max     float64 `xml:"max"`
			Azimuth int     `xml:"azimuth"`
		} `xml:"wind"`
		Rain struct {
			Unit  string  `xml:"unit"`
			Total float64 `xml:"total"`
		} `xml:"rain"`
	} `xml:"stationdata"`
}

func New(id string) (MeteoClimatic, error) {
	data, err := getXMLInfo(id)
	if err != nil {
		return MeteoClimatic{}, err
	}

	return MeteoClimatic{XMLData: data}, nil
}

// method to get the XML from Meteoclimatic
// it is separated so we can mock the tests
func getXMLInfo(id string) ([]byte, error) {
	ua := &http.Client{Timeout: 10 * time.Second}
	stationURL := fmt.Sprintf(baseURL+"%s", id)

	r, err := ua.Get(stationURL)
	if err != nil {
		return nil, fmt.Errorf("Could not get station info: %q", err)
	}
	defer r.Body.Close()

	resBody, err := ioutil.ReadAll(r.Body)
	return resBody, nil
}

// Parse XML and return a struct with all the data
func (mc *MeteoClimatic) GetStationInfo() error {
	decoder := xml.NewDecoder(bytes.NewReader(mc.XMLData))
	decoder.CharsetReader = charset.NewReaderLabel
	err := decoder.Decode(&mc.MD)

	if err != nil {
		return fmt.Errorf("could not decode xml: %q", err)
	}

	return nil
}
