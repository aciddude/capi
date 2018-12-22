package main

import (

	"log"
	"github.com/gorilla/mux"
	"net/http"
	"capi/api"
)


var GetBlock = api.GetBlock

var GetTX = api.GetTX

var GetCoinCodexData = api.GetCoinCodexData




func main() {

	router := mux.NewRouter()
	router.HandleFunc("/block/{HeightOrHash}", GetBlock).Methods("GET")
	router.HandleFunc("/tx/{tx}", GetTX).Methods("GET")
	router.HandleFunc("/market", GetCoinCodexData).Methods("GET")
	log.Println("capi v0.1 is running!")
	log.Fatal(http.ListenAndServe(":8000", router))


}
