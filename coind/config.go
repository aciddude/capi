package coind

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	Coin               string `json:"Coin"`
	Ticker             string `json:"Ticker"`
	RPCHost            string `json:"RPCHost"`
	RPCUser            string `json:"RPCUser"`
	RPCPassword        string `json:"RPCPassword"`
	RPCPORT            string `json:"RPCPort"`
	RPCTimeout         int    `json:"RPCTimeout"`
	SSL                bool   `json:"SSL"`
	EnableCoinCodexAPI bool   `json:"EnableCoinCodexAPI"`
	CapiPort           string `json:"capi_port"`
}

func LoadConfig(file string) Config {
	// get the local config from disk
	//filename is the path to the json config file

	var config Config
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		log.Fatal("ERROR: Could not find config file \n GoLang Error:  ", err)
	}

	decoder := json.NewDecoder(configFile)
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatal("ERROR: Could not decode json config  \n GoLang Error:  ", err)
	}
	return config
}
