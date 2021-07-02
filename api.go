package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type knnNode struct {
	Distancia float64
	x         int
	y         int
	estado    string
}
type Respuesta struct {
	Mensaje string
}

type Ong struct {
	Numero        int    `json:"numero"`
	Institucion   string `json:"institucion"`
	Departamento  string `json:"departamento"`
	Provincia     string `json:"provincia"`
	Distrito      string `json:"distrito"`
	Representante string `json:"representante"`
	Sector        string `json:"sector"`
}

//variables globales
var Dataset = [1000]Ong{}
var eschucha_funcion bool
var remotehost string
var chCont chan int
var n, min, valorUsuario int

func enviar(num int) { //enviar el numero mayor al host remoto
	conn, _ := net.Dial("tcp", remotehost)
	defer conn.Close()
	//envio el número
	fmt.Fprintf(conn, "%d\n", num)
}

func enviar_Principal(num int) { //enviar el numero mayor al host remoto
	conn, _ := net.Dial("tcp", "localhost:8000")
	defer conn.Close()
	//envio el número
	fmt.Fprintf(conn, "%d\n", num)
}

func manejador_respueta(conn net.Conn) bool {
	defer conn.Close()
	eschucha_funcion = false
	bufferIn := bufio.NewReader(conn)
	numStr, _ := bufferIn.ReadString('\n')
	numStr = strings.TrimSpace(numStr)
	numero, _ := strconv.Atoi(numStr)
	strNumero := strconv.Itoa(numero)
	if strNumero[1] == 49 {
		return true
	} else {
		return false
	}
}

func knn(usuario *Ong) bool {
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
		log.Println("-----------------------------------")
		return true
	} else {
		log.Println("-----------------------------------")
		return false
	}
}

func LeerDataSetFromGit() {
	response, err := http.Get("https://raw.githubusercontent.com/sigiandre/TF-Programacion-Concurrente-y-Distribuida-Backend/master/dataset/Base-de-Datos-de-las-ONGD-I-Trimestre-2018_0.csv") //use package "net/http"
	if err != nil {
		log.Println(err)
		return
	}
	defer response.Body.Close()
	reader := csv.NewReader(response.Body)
	reader.Comma = ','
	if err != nil {
		log.Println(nil)
	}
	log.Println(Dataset)
}

func mostrarDataset(res http.ResponseWriter, req *http.Request) {
	allowedHeaders := "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization,X-CSRF-Token"
	log.Println("Llamada al endpoint /dataset")
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.Header().Set("Access-Control-Allow-Origin", "*")
	res.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	res.Header().Set("Access-Control-Allow-Headers", allowedHeaders)
	res.Header().Set("Access-Control-Expose-Headers", "Authorization")
	jsonBytes, _ := json.MarshalIndent(Dataset, "", "\t")
	io.WriteString(res, string(jsonBytes))
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
	temp := "1"
	conn, _ := net.Dial("tcp", "localhost:8001")
	defer conn.Close()

	i, _ := strconv.Atoi(temp)
	log.Println(i)

	fmt.Fprintf(conn, "%d\n", i)

	ln, _ := net.Listen("tcp", "localhost:8000")
	defer ln.Close()
	eschucha_funcion = true
}

func handleRequest() {

	http.HandleFunc("/dataset", mostrarDataset)
	http.HandleFunc("/knn", realizarKnn)
	log.Fatal(http.ListenAndServe(":9200", nil))
}

func main() {
	bufferIn := bufio.NewReader(os.Stdin)
	LeerDataSetFromGit()
	//tipo de nodo
	log.Print("Ingrese el tipo de nodo (i:inicio -n:intermedio - f:final): ")
	tipo, _ := bufferIn.ReadString('\n')
	tipo = strings.TrimSpace(tipo)

	if tipo == "i" {
		handleRequest()
	}
	if tipo == "n" {
		//establecer el identificador del host local (IP:puerto)
		log.Print("Ingrese el puerto local: ")
		puerto, _ := bufferIn.ReadString('\n')
		puerto = strings.TrimSpace(puerto)
		localhost := ("localhost:" + puerto)

		//establecer el identificador del host remoto (IP:puerto)
		log.Print("Ingrese el puerto remoto:")
		puerto, _ = bufferIn.ReadString('\n')
		puerto = strings.TrimSpace(puerto)
		remotehost = ("localhost:" + puerto)

		//Cantidad de numero a recibir x nodo

		//canal para el contador
		chCont = make(chan int, 1) //canal asincrono
		chCont <- 0

		//establecer el modo escucha del nodo
		ln, _ := net.Listen("tcp", localhost)
		defer ln.Close()
		for {
			//manejador de conexiones
			conn, _ := ln.Accept()
			go manejador_respueta(conn)
		}
	}
	if tipo == "f" {
		//establecer el identificador del host local (IP:puerto)
		log.Print("Ingrese el puerto local: ")
		puerto, _ := bufferIn.ReadString('\n')
		puerto = strings.TrimSpace(puerto)
		localhost := ("localhost:" + puerto)

		//establecer el identificador del host remoto (IP:puerto)

		//Cantidad de numero a recibir x nodo

		//canal para el contador
		chCont = make(chan int, 1) //canal asincrono
		chCont <- 0

		//establecer el modo escucha del nodo
		ln, _ := net.Listen("tcp", localhost)
		defer ln.Close()
		for {
			//manejador de conexiones
			conn, _ := ln.Accept()
			go manejador_respueta(conn)
		}
	}

}
