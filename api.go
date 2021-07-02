package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
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

// Get single ong
func getONG(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	for _, item := range bdongs {
		numero, _ := strconv.Atoi(params["id"])
		if item.Numero == numero {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&BDONG{})
}

// Add new ong
func createOng(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var bdong BDONG
	_ = json.NewDecoder(r.Body).Decode(&bdong)
	bdongs = append(bdongs, bdong)
	json.NewEncoder(w).Encode(bdong)
}

// Get single ong
func getCategory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	w.Header().Set("Content-Type", "application/json")

	var bdong BDONG
	_ = json.NewDecoder(r.Body).Decode(&bdong)

	k := 20 + rand.Intn(20)
	fmt.Println(k)

	bdongs = append(bdongs, bdong)

	//payload, _ := json.MarshalJson()
	//w.Write(payload)
}

func main() {
	//filePathUrl := "dataset/Base-de-Datos-de-las-ONGD-I-Trimestre-2018_0.csv"
	filePathUrl := "https://raw.githubusercontent.com/sigiandre/TF-Programacion-Concurrente-y-Distribuida-Backend/master/dataset/Base-de-Datos-de-las-ONGD-I-Trimestre-2018_0.csv"
	lines, err := readFileUrl(filePathUrl)
	if err != nil {
		panic(err)
	}
	fmt.Println("Leyo archivos")
	lineToStruc(lines)
	fmt.Println("Parseo Archivos")

	r := mux.NewRouter()

	r.HandleFunc("/ongs", getONGS).Methods("GET")
	r.HandleFunc("/ongs/{id}", getONG).Methods("GET")
	r.HandleFunc("/ongs", createOng).Methods("POST")
	r.HandleFunc("/knn", getCategory).Methods("POST")

	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})
	origins := handlers.AllowedOrigins([]string{"*"})

	// Start server
	port := ":8000"
	fmt.Println("Escuchando en " + port)
	//main3()
	log.Fatal(http.ListenAndServe(port, handlers.CORS(headers, methods, origins)(r)))

}
