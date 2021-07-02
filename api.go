//https://www.datosabiertos.gob.pe/dataset/base-de-datos-organizaciones-no-gubernamentales-de-desarrollo-ongd/resource/8a550153-0d26

package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var bdongs []BDONG

type knnNode struct {
	Distancia float64
	x         int
	y         int
	estado    string
}
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

func knn(usuario *BDONG) bool {
	var knnNodes = [100]knnNode{}
	chDistancia := make(chan float64)
	chY := make(chan int)
	chX := make(chan int)
	chEstado := make(chan string)
	for i := 0; i < 100; i++ {
		knnNodes[i].Distancia = <-chDistancia
		knnNodes[i].y = <-chY
		knnNodes[i].x = <-chX
		knnNodes[i].estado = <-chEstado
	}
	log.Println(knnNodes)
	for i := 1; i < 100; i++ {
		for j := 0; j < 100-i; j++ {
			if knnNodes[j].Distancia > knnNodes[j+1].Distancia {
				knnNodes[j], knnNodes[j+1] = knnNodes[j+1], knnNodes[j]
			}
		}
	}
	log.Println(knnNodes)
	count := 0
	for i := 0; i < 6; i++ {
		if knnNodes[i].estado == "ONG" {
			count++
		}
	}
	if count >= 3 {
		log.Println("-------------------------------------------")
		return true
	} else {
		log.Println("-------------------------------------------")
		return false
	}
}

// Add new ong
func createOng(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var bdong BDONG
	_ = json.NewDecoder(r.Body).Decode(&bdong)
	bdongs = append(bdongs, bdong)
	json.NewEncoder(w).Encode(bdong)
}

func realizarKnn(res http.ResponseWriter, req *http.Request) {
	allowedHeaders := "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization,X-CSRF-Token"
	log.Println("Llamada al endpoint /knn")
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.Header().Set("Access-Control-Allow-Origin", "*")
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.Header().Set("Access-Control-Allow-Origin", "*")
	res.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	res.Header().Set("Access-Control-Allow-Headers", allowedHeaders)
	res.Header().Set("Access-Control-Expose-Headers", "Authorization")

	conn, _ := net.Dial("tcp", "localhost:8001")
	defer conn.Close()

	ln, _ := net.Listen("tcp", "localhost:8000")
	defer ln.Close()
}

func main() {
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
	r.HandleFunc("/knn", realizarKnn).Methods("POST")

	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})
	origins := handlers.AllowedOrigins([]string{"*"})

	// Start server
	port := ":8000"
	fmt.Println("Escuchando en " + port)
	log.Fatal(http.ListenAndServe(port, handlers.CORS(headers, methods, origins)(r)))

}
