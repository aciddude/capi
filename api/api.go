package api

import (
	"net/http"
	"log"
	"strings"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"encoding/json"
	"strconv"
	"os"
	"github.com/btcsuite/btcd/rpcclient"
	"time"
	"io/ioutil"
)


type Config struct {
	Coin				string	`json:"Coin"`
	Ticker				string	`json:"Ticker"`
	Daemon				string	`json:"Daemon"`
	RPCUser				string	`json:"RPCUser"`
	RPCPassword	 		string	`json:"RPCPassword"`
	HTTPPostMode		bool
	DisableTLS			bool
	EnableCoinCodexAPI 	bool	`json:"EnableCoinCodexAPI"`
}


// Block struct
type Block struct {
	Hash              string   `json:"hash"`
	Confirmations     int64    `json:"confirmations"`
	Size              int32    `json:"size"`
	StrippedSize  	  int32	   `json:"strippedSize"`
	Weight      	  int32	   `json:"weight"`
	Height            int64    `json:"height"`
	Version           int32	   `json:"version"`
	VersionHex  	  string   `json:"versionHex"`
	MerkleRoot        string   `json:"merkleRoot"`
	BlockTransactions []string `json:"tx"`
	Time              int64    `json:"time"`
	Nonce             uint32   `json:"nonce"`
	Bits              string   `json:"bits"`
	Difficulty        float64  `json:"difficulty"`
	PreviousHash      string   `json:"previousBlockHash"`
	NextHash          string   `json:"nextBlockHash"`

}

var Blocks []Block

/// Function to Load Config file from disk as type Config struct

func loadConfig(file string) (Config) {
	// get the local config from disk
	//filename is the path to the json config file

	var config Config
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		log.Fatal("ERROR: Could not find config file \n GoLang Error:  " , err)
	}

	decoder := json.NewDecoder(configFile)
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatal("ERROR: Could not decode json config  \n GoLang Error:  " , err)
	}
	return config
}

// config file from disk using loadConfig function
var configFile  = loadConfig("./config/config.json")

// coin client using coinClientConfig
var coinClient, _ = rpcclient.New(coinClientConfig, nil)

// coin client config for coinClient, loads values from configFile
var coinClientConfig = &rpcclient.ConnConfig {
	Host: 			configFile.Daemon,
	User: 			configFile.RPCUser,
	Pass: 			configFile.RPCPassword,
	HTTPPostMode:	configFile.HTTPPostMode,
	DisableTLS:		configFile.DisableTLS,
}


//GetBlockObject
func GetBlock(w http.ResponseWriter, r *http.Request)  {


	urlBlock := r.URL.Path
	if len(urlBlock) > 60{

		urlBlock = strings.TrimPrefix(urlBlock, "/block/")

		log.Println("Block Hash", urlBlock)

		hash, err := chainhash.NewHashFromStr(urlBlock)
		if err != nil {
			log.Print("Error with hash")
		}

		log.Println("Trying to get block hash: ", urlBlock)
		block, err := coinClient.GetBlockVerbose(hash)
		if err != nil {
			log.Print("Error with hash requested: ", urlBlock)
			http.Error(w, "ERROR: invalid block hash requested \n"+
				"Please use a block hash, eg: 4b6c3362e2f2a6b6317c85ecaa0f5415167e2bb333d2bf3d3699d73df613b91f", 500)
			return
		}
		blockArray := make([]Block, 1);
		index := 0;
		blockArray[index] = Block {
			Height:  block.Height,
			Hash: block.Hash,
			Bits: block.Bits,
			BlockTransactions: block.Tx,
			Confirmations: block.Confirmations,
			Difficulty: block.Difficulty,
			Nonce: block.Nonce,
			Time: block.Time,
			MerkleRoot: block.MerkleRoot,
			PreviousHash: block.PreviousHash,
			NextHash: block.NextHash,
			Size: block.Size,
			VersionHex: block.VersionHex,
			StrippedSize: block.StrippedSize,
			Weight: block.Weight,
			Version: block.Version,

		};

		log.Println(blockArray)

		jsonBlock, err := json.Marshal(&blockArray)
		data := json.RawMessage(jsonBlock)
		json.NewEncoder(w).Encode(data)


	} else {
		urlBlock = strings.TrimPrefix(urlBlock, "/block/")
		log.Println("Parsed Block Object from the URL", urlBlock)

		blockHeight, err := strconv.ParseInt(urlBlock, 10, 64)
		if err != nil {
			log.Println("ERROR: invalid block height specified" + " -- Go Error:" ,err)
			http.Error(w, "ERROR: invalid block height specified \n"+"Please chose a number like '0' for the genesis block or '444' for block 444", 404)
			return
		}

		log.Println("Block converted to int64", blockHeight)
		blockHash, err := coinClient.GetBlockHash(blockHeight)
		if err != nil {
			log.Println(err)
			http.Error(w, "ERROR Getting Block Hash from Height: \n"+err.Error(), 500)

		}

		block, err := coinClient.GetBlockVerbose(blockHash)
		if err != nil {
			log.Println(err)
			http.Error(w, "ERROR Getting Block from Block Hash:  "+err.Error(), 500)
		}

		blockArray := make([]Block, 1);
		index := 0;
		blockArray[index] = Block {
			Height:  block.Height,
			Hash: block.Hash,
			Bits: block.Bits,
			BlockTransactions: block.Tx,
			Confirmations: block.Confirmations,
			Difficulty: block.Difficulty,
			Nonce: block.Nonce,
			Time: block.Time,
			MerkleRoot: block.MerkleRoot,
			PreviousHash: block.PreviousHash,
			NextHash: block.NextHash,
			Size: block.Size,
			VersionHex: block.VersionHex,
			StrippedSize: block.StrippedSize,
			Weight: block.Weight,
			Version: block.Version,

		};

		log.Println(blockArray)

		jsonBlock, err := json.Marshal(&blockArray)
		data := json.RawMessage(jsonBlock)
		json.NewEncoder(w).Encode(data)
	}
}


//GetTX
func GetTX(w http.ResponseWriter, r *http.Request) {

	request := r.URL.Path
	request = strings.TrimPrefix(request, "/tx/")

	log.Println("Parsed txid from request ", request)

	requestHash, err := chainhash.NewHashFromStr(request)
	if err != nil {
		log.Println("ERROR:", err)
		return

	}

	txhash, err := coinClient.GetRawTransactionVerbose(requestHash)
	if err != nil {
		log.Println("ERROR:", err)
		http.Error(w, "ERROR: invalid transaction id specified \n"+err.Error(), 404)
		return

	}

	json.NewEncoder(w).Encode(txhash)

}



/// CoinCodex.com API for prices

type coincodexapi struct {
	Symbol 					string	`json:"symbol"`
	CoinName 				string  `json:"coin_name"`
	Price_TodayOpenUSD		float64	`json:"today_open"`
	Price_HighUSD			float64	`json:"price_high_24_usd"`
	Price_LowUSD			float64	`json:"price_low_24_usd"`
	Volume24USD				float64	`json:"volume_24_usd"`
	DataProvider			string	`json:"data_provider"`
}

func GetCoinCodexData(w http.ResponseWriter, r *http.Request) {
	if configFile.EnableCoinCodexAPI == false {
		return
	} else {

		url := "https://coincodex.com/api/coincodex/get_coin/"+configFile.Ticker

		client := http.Client{
			Timeout:time.Second * 5, // 5 second timeout
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
			DataProvider:"CoinCodex.com",
		}
		jsonError := json.Unmarshal(body, &jsonData)
		if jsonError != nil {
			log.Println("ERROR: ", jsonError)
		}

		json.NewEncoder(w).Encode(jsonData)
	}
}

