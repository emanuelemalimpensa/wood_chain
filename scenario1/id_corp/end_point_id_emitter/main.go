package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

type IDs struct {
	Num int64  `json:"Num"`
	ID  string `json:"ID"`
}

func IDsRequest(w http.ResponseWriter, r *http.Request) {
	fmt.Println("AUDIT REQUEST")
	type Quantita struct {
		Quantita   string `json:Quantita`
		ReqCompany string `json:ReqCompany`
		Tipo       string `json:Tipo`
	}
	var prova Quantita
	err := json.NewDecoder(r.Body).Decode(&prova)
	//GENERAZIONE DEGLI ID RICHIESTI CON LA AUDIT REQUEST
	myurl := "http://localhost:5000/api/product/audit"
	bodyVal := fmt.Sprintf("{\"Quantita\":\"%s\",\"Owner\":\"%s\"}", prova.Quantita, prova.ReqCompany)
	body := strings.NewReader(bodyVal)
	_, err2 := http.Post(myurl, "application/json", body)
	if err2 != nil {
		log.Fatal(err2)
		os.Exit(1)
	}

	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	url := os.Getenv("URL")

	//ACCESSO AGLI ID DISPONIBILI (APPENA GENERATI)
	response, err := http.Get("http://" + host + ":" + port + url)
	var dataJ2Struct []IDs

	if err != nil {
		fmt.Printf("the HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		err = json.Unmarshal(data, &dataJ2Struct)
		if err != nil {
			log.Fatal(err)
		}

		// err = json.NewDecoder(r.Body).Decode(&prova)
		data2, _ := ioutil.ReadAll(r.Body)
		fmt.Printf("%v\n", dataJ2Struct)
		//RICHIESTE DI INSERIMENTO AD API CHE INTERAGISCE CON CHAINCODE
		//UNA PER OGNI ID DISPONIBILE FORNITO
		q, _ := strconv.Atoi(prova.Quantita)
		fmt.Print(data2)
		fmt.Print(q)
		for i := 0; i < q; i++ {
			fmt.Printf("%v\n", dataJ2Struct[i])
			const myurl = "http://localhost:3000/mock?channelid=mychannel&chaincodeid=basic2"
			// BODY {
			// 	ID:     ID,
			// 	Tipo:   Tipo,
			// 	Owner:  Owner,
			// 	Issuer: Issuer,
			// 	Free:   "true",
			// }
			bodyVal := fmt.Sprintf("{\"ID\":\"%s\",\"Issuer\":\"EmitterCompany\",\"Owner\":\"%s\",\"Tipo\":\"%s\"}", dataJ2Struct[i].ID, prova.ReqCompany, prova.Tipo)

			body := strings.NewReader(bodyVal)
			//INSERIMENTO IN BLOCKCHAIN
			_, err := http.Post(myurl, "application/json", body)
			if err != nil {
				log.Fatal(err)
				os.Exit(1)
			}
			//ETICHETTATURA ID COME GIA' IMPIEGATO E QUINDI NON PIU' UTILIZZABILE
			myurl2 := fmt.Sprintf("http://localhost:5000/api/product/inserted/%s", dataJ2Struct[i].ID)
			_, err2 := http.Post(myurl2, "application/json", nil)
			if err2 != nil {
				log.Fatal(err2)
				os.Exit(1)
			}

		}
	}

}

func main() {
	serverPort := os.Getenv("SERVERPORT")
	fmt.Printf("server listening on port %s ", serverPort)
	router := mux.NewRouter()
	router.HandleFunc("/api/product/audit", IDsRequest).Methods("POST")
	err := http.ListenAndServe(":"+serverPort, router)

	if err != nil {
		fmt.Println(err)
	}
}
