//
// aemet_test.go
// Copyright (C) 2017 Paco Esteban <paco@onna.be>
//
// Distributed under terms of the MIT license.
//

package aemet

import (
	"fmt"
	"testing"
	"time"
)

func TestForecastTownHourly(t *testing.T) {
	hf, err := ForecastTownHourly("08121")
	if err != nil {
		t.Fatal(err)
	}

	if hf[0].ID != "08121" {
		t.Fatal("Id is not right")
	}
	if hf[0].Version != "1.0" {
		t.Fatal("Version is not right")
	}
	if hf[0].Origen.Enlace != "http://www.aemet.es/es/eltiempo/prediccion/municipios/horas/mataro-id08121" {
		t.Fatal("Origin is not right")
	}

	now := time.Now()

	if hf[0].Elaborado != fmt.Sprintf("%d-%02d-%02d", now.Year(), now.Month(), now.Day()) {
		t.Fatalf("Report not elaborated today: %s", hf[0].Elaborado)
	}

}
