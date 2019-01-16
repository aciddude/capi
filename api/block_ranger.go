package api

import (
	"log"
	"os"

	"capi/coind"
)

// DB Block Struct

type dbBlock struct {
	ID                int64  `json:"-"` /// The Primary Key
	Hash              string `storm:"index"`
	Confirmations     int64
	Size              int32
	StrippedSize      int32
	Weight            int32
	Height            int64 `storm:"index"`
	Version           int32
	VersionHex        string
	MerkleRoot        string
	BlockTransactions []string `storm:"index"`
	Time              int64
	Nonce             uint32
	Bits              string
	Difficulty        float64
	PreviousHash      string
	NextHash          string
}

type dbTX struct {
	ID            int64        `json:"-",storm:"increment,index"` /// The Primary Key
	Hex           string       `json:"hex"`
	Txid          string       `storm:"index"`
	Hash          string       `json:"hash,omitempty"`
	Size          int32        `json:"size,omitempty"`
	Vsize         int32        `json:"vsize,omitempty"`
	Version       int32        `json:"version"`
	LockTime      uint32       `json:"locktime"`
	Vin           []coind.Vin  `storm:"inline"`
	Vout          []coind.Vout `storm:"inline"`
	BlockHash     string       `json:"blockhash,omitempty"`
	Confirmations uint64       `json:"confirmations,omitempty"`
	Time          int64        `json:"time,omitempty"`
	Blocktime     int64        `json:"blocktime,omitempty"`
}

func DBChecker() {

	if _, err := os.Stat("blocks.db"); os.IsNotExist(err) {
		log.Println("Running Block Ranger")

	} else {
		log.Println("DB Exists, checking DB.......")
	}

}
