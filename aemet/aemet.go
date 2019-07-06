//
// aemet.go
// Copyright (C) 2017 Paco Esteban <paco@onna.be>
//
// Distributed under terms of the MIT license.
//

package aemet

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

const (
	baseURL string = "https://opendata.aemet.es/opendata"
)

type RespData struct {
	Descripcion string `json:"descripcion"`
	Estado      int    `json:"estado"`
	Datos       string `json:"datos"`
	Metadatos   string `json:"metadatos"`
}

type HourlyForecast []struct {
	Elaborado string `json:"elaborado"`
	ID        string `json:"id"`
	Nombre    string `json:"nombre"`
	Origen    struct {
		Copyright string `json:"copyright"`
		Enlace    string `json:"enlace"`
		Language  string `json:"language"`
		NotaLegal string `json:"notaLegal"`
		Productor string `json:"productor"`
		Web       string `json:"web"`
	} `json:"origen"`
	Prediccion struct {
		Dia []struct {
			EstadoCielo []struct {
				Descripcion string `json:"descripcion"`
				Periodo     string `json:"periodo"`
				Value       string `json:"value"`
			} `json:"estadoCielo"`
			Fecha           string `json:"fecha"`
			HumedadRelativa []struct {
				Periodo string `json:"periodo"`
				Value   string `json:"value"`
			} `json:"humedadRelativa"`
			Nieve []struct {
				Periodo string `json:"periodo"`
				Value   string `json:"value"`
			} `json:"nieve"`
			Ocaso         string `json:"ocaso"`
			Orto          string `json:"orto"`
			Precipitacion []struct {
				Periodo string `json:"periodo"`
				Value   string `json:"value"`
			} `json:"precipitacion"`
			ProbNieve []struct {
				Periodo string `json:"periodo"`
				Value   string `json:"value"`
			} `json:"probNieve"`
			ProbPrecipitacion []struct {
				Periodo string `json:"periodo"`
				Value   string `json:"value"`
			} `json:"probPrecipitacion"`
			ProbTormenta []struct {
				Periodo string `json:"periodo"`
				Value   string `json:"value"`
			} `json:"probTormenta"`
			SensTermica []struct {
				Periodo string `json:"periodo"`
				Value   string `json:"value"`
			} `json:"sensTermica"`
			Temperatura []struct {
				Periodo string `json:"periodo"`
				Value   string `json:"value"`
			} `json:"temperatura"`
			VientoAndRachaMax []struct {
				Direccion []string `json:"direccion,omitempty"`
				Periodo   string   `json:"periodo"`
				Velocidad []string `json:"velocidad,omitempty"`
				Value     string   `json:"value,omitempty"`
			} `json:"vientoAndRachaMax"`
		} `json:"dia"`
	} `json:"prediccion"`
	Provincia string `json:"provincia"`
	Version   string `json:"version"`
}

// ForecastTownHourly returns a HourlyForecast struct given
// a Town id
func ForecastTownHourly(id string) (HourlyForecast, error) {
	var hf HourlyForecast

	url, err := getDataURL("/api/prediccion/especifica/municipio/horaria/" + id)
	if err != nil {
		return HourlyForecast{}, err
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	ua := &http.Client{Timeout: 10 * time.Second, Transport: tr}

	r, err := ua.Get(url)
	if err != nil {
		return hf, fmt.Errorf("Could not get forecast info: %q", err)
	}
	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(&hf); err != nil {
		return hf, fmt.Errorf("could not decode JSON: %v", err)
	}

	return hf, nil
}

// getDataURL performs the auth request and returns a url to
// get data from.
func getDataURL(subq string) (url string, err error) {
	var resp RespData
	token, ok := os.LookupEnv("AEMET_TOKEN")
	if !ok {
		return "", fmt.Errorf("AEMET_TOKEN variable is not set.")
	}

	reqURL := fmt.Sprintf("%s%s", baseURL, subq)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	ua := &http.Client{Timeout: 10 * time.Second, Transport: tr}
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return "", fmt.Errorf("Could not create request: %q", err)
	}
	req.Header.Add("api_key", token)
	req.Header.Add("Accept", "application/json")

	r, err := ua.Do(req)
	if err != nil {
		return "", fmt.Errorf("Could not get data URL: %q", err)
	}
	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(&resp); err != nil {
		return "", fmt.Errorf("could not decode Data URL JSON: %v", err)
	}

	if resp.Estado != 200 {
		return "", fmt.Errorf("Bad status: %v", r.Status)
	}

	return resp.Datos, nil
}
