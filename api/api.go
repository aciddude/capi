package api

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/aciddude/capi/coind"
	"github.com/aciddude/capi/datastore"

	"github.com/asdine/storm"
	bolt "go.etcd.io/bbolt"
)

var Blocks []coind.Block

var ConfigFile = coind.LoadConfig("./config/config.json")

var coinDaemon, err = coind.New(ConfigFile)

//GetBlockObject
func GetBlock(w http.ResponseWriter, r *http.Request) {

	urlBlock := r.URL.Path
	if len(urlBlock) > 60 {

		urlBlock = strings.TrimPrefix(urlBlock, "/block/")

		log.Println("Block Hash", urlBlock)

		log.Println("Trying to get block hash: ", urlBlock)
		block, err := coinDaemon.GetBlock(urlBlock)
		if err != nil {
			log.Print("Error with hash requested: ", urlBlock)
			http.Error(w, "ERROR: invalid block hash requested \n"+
				"Please use a block hash, eg: 4b6c3362e2f2a6b6317c85ecaa0f5415167e2bb333d2bf3d3699d73df613b91f", 500)
			return
		}

		jsonBlock, err := json.Marshal(&block)
		data := json.RawMessage(jsonBlock)
		json.NewEncoder(w).Encode(data)

	} else {
		urlBlock = strings.TrimPrefix(urlBlock, "/block/")
		log.Println("Parsed Block Object from the URL", urlBlock)

		blockHeight, err := strconv.ParseInt(urlBlock, 10, 64)
		if err != nil {
			log.Println("ERROR: invalid block height specified"+" -- Go Error:", err)
			http.Error(w, "ERROR: invalid block height specified \n"+"Please chose a number like '0' for the genesis block or '444' for block 444", 404)
			return
		}

		log.Println("Block converted to int64", blockHeight)
		blockHash, err := coinDaemon.GetBlockHash(blockHeight)
		if err != nil {
			log.Println(err)
			http.Error(w, "ERROR Getting Block Hash from Height: \n"+err.Error(), 500)

		}

		block, err := coinDaemon.GetBlock(blockHash)
		if err != nil {
			log.Println(err)
			http.Error(w, "ERROR Getting Block from Block Hash:  "+err.Error(), 500)
		}

		jsonBlock, err := json.Marshal(&block)
		data := json.RawMessage(jsonBlock)
		json.NewEncoder(w).Encode(data)
	}
}

//GetTX
func GetTX(w http.ResponseWriter, r *http.Request) {

	request := r.URL.Path
	request = strings.TrimPrefix(request, "/tx/")

	log.Println("Parsed txid from request ", request)

	txhash, err := coinDaemon.GetRawTransaction(request, true)
	if err != nil {
		log.Println("ERROR:", err)
		http.Error(w, "ERROR: invalid transaction id specified \n"+err.Error(), 404)
		return

	}

	json.NewEncoder(w).Encode(txhash)

}

//GetBlockchainInfo
func GetBlockchainInfo(w http.ResponseWriter, r *http.Request) {

	getblockchaininfo, err := coinDaemon.GetBlockchainInfo()
	if err != nil {
		log.Println("ERROR: ", err)
		http.Error(w, "ERROR: \n"+err.Error(), 500)
		return
	}

	json.NewEncoder(w).Encode(getblockchaininfo)

}

/// CoinCodex.com API for prices

type coincodexapi struct {
	Symbol             string  `json:"symbol"`
	CoinName           string  `json:"coin_name"`
	LastPrice          float64 `json:"last_price_usd"`
	Price_TodayOpenUSD float64 `json:"today_open"`
	Price_HighUSD      float64 `json:"price_high_24_usd"`
	Price_LowUSD       float64 `json:"price_low_24_usd"`
	Volume24USD        float64 `json:"volume_24_usd"`
	DataProvider       string  `json:"data_provider"`
}

func GetCoinCodexData(w http.ResponseWriter, r *http.Request) {
	if ConfigFile.EnableCoinCodexAPI == false {
		return
	} else {

		url := "https://coincodex.com/api/coincodex/get_coin/" + ConfigFile.Ticker

		client := http.Client{
			Timeout: time.Second * 5, // 5 second timeout
		}

		request, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			log.Println("ERROR: ", err)
		}
		request.Header.Set("User-Agent", "capi v0.1")

		response, getError := client.Do(request)
		if getError != nil {
			log.Println("ERROR: ", getError)

		}
		body, readError := ioutil.ReadAll(response.Body)
		if readError != nil {
			log.Println("ERROR:", readError)
		}

		jsonData := coincodexapi{
			DataProvider: "CoinCodex.com",
		}
		jsonError := json.Unmarshal(body, &jsonData)
		if jsonError != nil {
			log.Println("ERROR: ", jsonError)
		}

		json.NewEncoder(w).Encode(jsonData)
	}
}

/// To show a simple index page with the coin price info
func IndexRoute(w http.ResponseWriter, r *http.Request) {

	tmpl, err := template.ParseFiles("templates/index.tmpl")
	if err != nil {
		log.Println("ERROR: Parsing template file index.tmpl", err)
	}

	url := "https://coincodex.com/api/coincodex/get_coin/" + ConfigFile.Ticker

	client := http.Client{
		Timeout: time.Second * 5, // 5 second timeout
	}

	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Println("ERROR: ", err)
	}
	request.Header.Set("User-Agent", "capi v0.1")

	response, getError := client.Do(request)
	if getError != nil {
		log.Println("ERROR: ", getError)

	}
	body, readError := ioutil.ReadAll(response.Body)
	if readError != nil {
		log.Println("ERROR:", readError)
	}

	jsonData := coincodexapi{
		DataProvider: "CoinCodex.com",
	}
	jsonError := json.Unmarshal(body, &jsonData)
	if jsonError != nil {
		log.Println("ERROR: ", jsonError)
	}

	tmpl.Execute(w, jsonData)
}

/// GetBlock from DB using Height. Couldn't use height has Primary index as bolt doesn't like '0' or block 0
/// Made my own dbBlock type with it's own "ID"
func GetBlockFromDBHeight(w http.ResponseWriter, r *http.Request) {

	db, err := storm.Open("blocks.db", storm.BoltOptions(600, &bolt.Options{Timeout: 5 * time.Second}))
	if err != nil {
		log.Println("ERROR: Cannot open DB", err)
	}

	request := r.URL.Path
	request = strings.TrimPrefix(request, "/blkdb/")

	var block datastore.BlockDB
	log.Println("Parse Request URL", request)
	requestInt, err := strconv.ParseUint(request, 10, 64)
	if err != nil {
		log.Println("ERROR: cannot parse URL", err)
		http.Error(w, "ERROR: Could not parse URL \n"+err.Error(), 500)
		db.Close()
		return
	}
	/// Silly fix, Height  0 = DB Index 1
	if requestInt == 0 {
		err := db.One("ID", 1, &block)
		if err != nil {
			log.Println("ERROR: Block not found in DB", err)
			http.Error(w, "ERROR: Block not found in DB \n"+err.Error(), 404)
			db.Close()
			return
		}
	} else {
		err := db.One("ID", requestInt+1, &block)
		if err != nil {
			log.Println("ERROR: Block not found in DB", err)
			http.Error(w, "ERROR: Block not found in DB \n"+err.Error(), 404)
			db.Close()
			return
		}
	}

	log.Println("DB Request: ", request, block)
	json.NewEncoder(w).Encode(block)

	db.Close()

}

// Get TX from TXID
func GetTXFromDB(w http.ResponseWriter, r *http.Request) {

	db, err := storm.Open("tx.db", storm.BoltOptions(600, &bolt.Options{Timeout: 5 * time.Second}))
	if err != nil {
		log.Println("ERROR: Cannot open TX DB", err)
	}

	request := r.URL.Path
	request = strings.TrimPrefix(request, "/txdb/")

	var tx datastore.TransactionDB
	log.Println("Parse Request URL", request)
	if err != nil {
		log.Println("ERROR: cannot parse URL", err)
		http.Error(w, "ERROR: Could not parse URL \n"+err.Error(), 500)
		db.Close()
		return
	}
	err = db.One("Txid", request, &tx)
	if err != nil {
		log.Println("ERROR: TX not found in DB", err)
		http.Error(w, "ERROR: TX not found in DB \n"+err.Error(), 404)
		db.Close()
		return
	}

	log.Println("DB Request: ", request, tx)
	json.NewEncoder(w).Encode(tx)
	db.Close()

}

func GetWalletTransactions(w http.ResponseWriter, r *http.Request) {

	db, err := storm.Open("addresses.db", storm.BoltOptions(600, &bolt.Options{Timeout: 5 * time.Second}))
	if err != nil {
		log.Println("ERROR: Cannot open TX DB", err)
	}

	request := r.URL.Path
	request = strings.TrimPrefix(request, "/wallet/")
	log.Println("Parse Request URL", request)

	var address []datastore.AddressDB

	query := db.Find("Address", []string{request}, &address)
	if err != nil {
		log.Println("ERROR: cannot parse URL", err)
		http.Error(w, "ERROR: Could not find Wallet \n"+err.Error(), 500)
		db.Close()
		return
	}
	log.Println("HELLO", query)

	json.NewEncoder(w).Encode(address)
	db.Close()
}
