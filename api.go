package main

import (
	"encoding/csv"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

var bdongs []BDONG

type BDONG struct {
	Numero        int    `json:"numero"`
	Institucion   string `json:"institucion"`
	Departamento  string `json:"departamento"`
	Provincia     string `json:"provincia"`
	Distrito      string `json:"distrito"`
	Representante string `json:"representante"`
	Sector        string `json:"sector"`
}

func lineToStruc(lines [][]string) {
	// Recorre líneas y conviértete en objeto
	for _, line := range lines {
		Numero, _ := strconv.Atoi(strings.TrimSpace(line[0]))

		bdongs = append(bdongs, BDONG{
			Numero:        Numero,
			Institucion:   strings.TrimSpace(line[1]),
			Departamento:  strings.TrimSpace(line[2]),
			Provincia:     strings.TrimSpace(line[3]),
			Distrito:      strings.TrimSpace(line[4]),
			Representante: strings.TrimSpace(line[5]),
			Sector:        strings.TrimSpace(line[6]),
		})
	}
}

func readFileUrl(filePathUrl string) ([][]string, error) {
	// Abrir archivo CSV
	f, err := http.Get(filePathUrl)
	if err != nil {
		return [][]string{}, err
	}
	defer f.Body.Close()

	// Leer archivo en una variable
	lines, err := csv.NewReader(f.Body).ReadAll()
	if err != nil {
		return [][]string{}, err
	}
	return lines, nil
}

// Get all ONG
func getONGS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bdongs)
}
