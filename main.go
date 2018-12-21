package main

import (

	"log"
	"github.com/gorilla/mux"
	"net/http"
	"capi/api"
)


var GetBlock = api.GetBlock

var GetTX = api.GetTX

var GetBLK = api.GetBlockObject



func main() {


	router := mux.NewRouter()
	router.HandleFunc("/block/{HeightOrHash}", GetBlock).Methods("GET")
	router.HandleFunc("/tx/{tx}", GetTX).Methods("GET")
	router.HandleFunc("/blk/{blk}", GetBLK).Methods("GET")
	log.Println("capi v0.1 is running!")
	log.Fatal(http.ListenAndServe(":8000", router))


}
