//
// meteoclimatic_test.go
// Copyright (C) 2017 Paco Esteban <paco@onna.be>
//
// Distributed under terms of the MIT license.
//

package meteoclimatic

import (
	"fmt"
	"io/ioutil"
	"log"
	"reflect"
	"testing"
)

func mockXMLData() ([]byte, error) {
	data, err := ioutil.ReadFile("./assets/tests/meteo.xml")
	if err != nil {
		log.Fatal(err)
		return nil, fmt.Errorf("Could not read XML file: %q", err)
	}

	return data, nil
}

func TestGetStationInfo(t *testing.T) {
	t1 := []struct {
		item     string
		expected string
	}{
		{"Description", "Meteoclimatic - XML"},
		{"Docs", "http://meteoclimatic.net/index/wp/xml_es.html"},
		{"Link", "http://meteoclimatic.net/"},
		{"PubDate", "Mon, 18 Sep 2017 13:29:34 +0000"},
	}
	t2 := []struct {
		item     string
		expected string
	}{
		{"Id", "ESCAT0800000008301F"},
		{"Location", "Mataró - Port"},
	}
	t3 := []struct {
		item     string
		expected float64
	}{
		{"Min", 1015.4},
		{"Max", 1017.6},
		{"Now", 1015.6},
	}

	data, err := mockXMLData()
	if err != nil {
		t.Fatal(err)
	}
	m := MeteoClimatic{XMLData: data}

	err = m.GetStationInfo()
	if err != nil {
		t.Fatalf("GetStationInfo failed: %q", err)
	}
	r := reflect.ValueOf(&m)
	for _, tc := range t1 {
		f := reflect.Indirect(r).FieldByName("MD").FieldByName(tc.item)
		if f.String() != tc.expected {
			t.Fatalf("%s is not %s but %q", tc.item, tc.expected, f.String())
		}
	}
	for _, tc := range t2 {
		f := reflect.Indirect(r).FieldByName("MD").FieldByName("MeteoStations").FieldByName("Station").Index(0).FieldByName(tc.item)
		if f.String() != tc.expected {
			t.Fatalf("%s is not %s but %q", tc.item, tc.expected, f.String())
		}
	}
	for _, tc := range t3 {
		f := reflect.Indirect(r).FieldByName("MD").FieldByName("MeteoStations").FieldByName("Station").Index(0).FieldByName("StationData").FieldByName("Barometre").FieldByName(tc.item)
		if f.Float() != tc.expected {
			t.Fatalf("%s is not %f but %f", tc.item, tc.expected, f.Float())
		}
	}
}
