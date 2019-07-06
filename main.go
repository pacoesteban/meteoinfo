//
// main.go
// Copyright (C) 2017 Paco Esteban <paco@onna.be>
//
// Distributed under terms of the MIT license.
//

package main

import (
	"fmt"
	"log"

	"github.com/pacoesteban/meteoinfo/aemet"
	mtc "github.com/pacoesteban/meteoinfo/meteoclimatic"
)

func main() {
	m, err := mtc.New("ESCAT0800000008301F")
	if err != nil {
		log.Fatalf("Meteoclimatic error, %v", err)
	}

	forecast, err := aemet.ForecastTownHourly("08121")
	if err != nil {
		log.Fatalf("AEMET Forecast ERROR: %q", err)
	}

	err = m.GetStationInfo()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(m.MD.Ttl)
	fmt.Println(forecast[0].Elaborado)
}
