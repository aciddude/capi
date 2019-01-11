package main

import (
	"log"
	"net/http"

	"github.com/aciddude/capi/api"
	"github.com/gorilla/mux"
)

var configFile = api.ConfigFile

func main() {

	api.DBChecker()

	router := mux.NewRouter()
	router.HandleFunc("/block/{HeightOrHash}", api.GetBlock).Methods("GET")
	router.HandleFunc("/tx/{tx}", api.GetTX).Methods("GET")
	router.HandleFunc("/market", api.GetCoinCodexData).Methods("GET")
	router.HandleFunc("/blockchaininfo", api.GetBlockchainInfo).Methods("GET")
	router.HandleFunc("/blkdb/{Height}", api.GetBlockFromDBHeight).Methods("GET")
	router.HandleFunc("/txdb/{txID}", api.GetTXFromDB).Methods("GET")
	router.HandleFunc("/wallet/{WalletAddress}", api.GetWalletBlance).Methods("GET")
	router.HandleFunc("/", api.IndexRoute).Methods("GET")
	log.Println("capi v0.1 is running!")
	log.Fatal(http.ListenAndServe(configFile.CapiPort, router))
}
