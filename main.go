package main

import (

	"log"
	"github.com/gorilla/mux"
	"net/http"
	"capi/api"
)



func main() {

	router := mux.NewRouter()
	router.HandleFunc("/block/{HeightOrHash}", api.GetBlock).Methods("GET")
	router.HandleFunc("/tx/{tx}", api.GetTX).Methods("GET")
	router.HandleFunc("/market", api.GetCoinCodexData).Methods("GET")
	router.HandleFunc("/blockchaininfo", api.GetBlockchainInfo).Methods("GET")
	log.Println("capi v0.1 is running!")
	log.Fatal(http.ListenAndServe(":8000", router))
}
