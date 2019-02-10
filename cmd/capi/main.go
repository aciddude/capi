package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/aciddude/capi/api"
	datastore "github.com/aciddude/capi/datastore/old"
)

var configFile = api.ConfigFile

func main() {

	DBExists := datastore.DatabaseExists()

	if DBExists == false {
		go datastore.StoreBlocks()
		go datastore.StoreAddresses()
		go datastore.StoreTransactions()
	}

	router := mux.NewRouter()
	router.HandleFunc("/block/{HeightOrHash}", api.GetBlock).Methods("GET")
	router.HandleFunc("/tx/{tx}", api.GetTX).Methods("GET")
	router.HandleFunc("/market", api.GetCoinCodexData).Methods("GET")
	router.HandleFunc("/blockchaininfo", api.GetBlockchainInfo).Methods("GET")
	router.HandleFunc("/blkdb/{Height}", api.GetBlockFromDBHeight).Methods("GET")
	router.HandleFunc("/txdb/{txID}", api.GetTXFromDB).Methods("GET")
	router.HandleFunc("/address/{WalletAddress}", api.GetWalletTransactions).Methods("GET")
	router.HandleFunc("/", api.IndexRoute).Methods("GET")
	log.Println("capi v0.1 is running!")
	log.Fatal(http.ListenAndServe(configFile.CapiPort, router))
}
