package api

import (
	"log"
	"os"
	"strconv"

	"github.com/asdine/storm"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
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
	ID            int64          `json:"-",storm:"increment,index"` /// The Primary Key
	Hex           string         `json:"hex"`
	Txid          string         `storm:"index"`
	Hash          string         `json:"hash,omitempty"`
	Size          int32          `json:"size,omitempty"`
	Vsize         int32          `json:"vsize,omitempty"`
	Version       int32          `json:"version"`
	LockTime      uint32         `json:"locktime"`
	Vin           []btcjson.Vin  `json:"vin"`
	Vout          []btcjson.Vout `json:"vout"`
	BlockHash     string         `json:"blockhash,omitempty"`
	Confirmations uint64         `json:"confirmations,omitempty"`
	Time          int64          `json:"time,omitempty"`
	Blocktime     int64          `json:"blocktime,omitempty"`
}

// Block struct
type Block struct {
	Hash              string   `json:"hash"`
	Confirmations     int64    `json:"confirmations"`
	Size              int32    `json:"size"`
	StrippedSize      int32    `json:"strippedSize"`
	Weight            int32    `json:"weight"`
	Height            int64    `json:"height"`
	Version           int32    `json:"version"`
	VersionHex        string   `json:"versionHex"`
	MerkleRoot        string   `json:"merkleRoot"`
	BlockTransactions []string `json:"tx"`
	Time              int64    `json:"time"`
	Nonce             uint32   `json:"nonce"`
	Bits              string   `json:"bits"`
	Difficulty        float64  `json:"difficulty"`
	PreviousHash      string   `json:"previousBlockHash"`
	NextHash          string   `json:"nextBlockHash"`
}

type TxRaw struct {
	Hex           string         `json:"hex"`
	Txid          string         `json:"txid"`
	Hash          string         `json:"hash,omitempty"`
	Size          int32          `json:"size,omitempty"`
	Vsize         int32          `json:"vsize,omitempty"`
	Version       int32          `json:"version"`
	LockTime      uint32         `json:"locktime"`
	Vin           []btcjson.Vin  `json:"vin"`
	Vout          []btcjson.Vout `json:"vout"`
	BlockHash     string         `json:"blockhash,omitempty"`
	Confirmations uint64         `json:"confirmations,omitempty"`
	Time          int64          `json:"time,omitempty"`
	Blocktime     int64          `json:"blocktime,omitempty"`
}

type ScriptPubKeyResult struct {
	Asm       string   `json:"asm"`
	Hex       string   `json:"hex,omitempty"`
	ReqSigs   int32    `json:"reqSigs,omitempty"`
	Type      string   `json:"type"`
	Addresses []string `json:"addresses,omitempty"`
}

type ScriptSig struct {
	Asm string `json:"asm"`
	Hex string `json:"hex"`
}

func BlockRanger(client *rpcclient.Client, startIndex int64, endIndex int64, blockCount int64) {

	//For each item (block) in the array index, print block details
	log.Println("starting block ranger at", strconv.FormatInt(startIndex, 10))

	logIndex := startIndex
	blockArray := make([]Block, 500)
	blockArrayIndex := 0
	for ; startIndex <= endIndex; startIndex++ {
		blockHash, err := client.GetBlockHash(startIndex)
		if err != nil {
			log.Println("Error getting block hash from height ", err)
		}

		block, err := client.GetBlockVerbose(blockHash)
		if err != nil {
			log.Println("Error getting block hash ", err)
		}

		blockArray[blockArrayIndex] = Block{
			Hash:              block.Hash,
			Confirmations:     block.Confirmations,
			Size:              block.Size,
			StrippedSize:      block.StrippedSize,
			Weight:            block.Weight,
			Height:            block.Height,
			Version:           block.Version,
			VersionHex:        block.VersionHex,
			MerkleRoot:        block.MerkleRoot,
			BlockTransactions: block.Tx,
			Time:              block.Time,
			Nonce:             block.Nonce,
			Bits:              block.Bits,
			Difficulty:        block.Difficulty,
			PreviousHash:      block.PreviousHash,
			NextHash:          block.NextHash,
		}
		blockArrayIndex = blockArrayIndex + 1

	}
	log.Println("Batch completed from:", logIndex, " to ", endIndex)
	log.Println("Sending to DB")

	for _, block := range blockArray {

		dbIndex := block.Height
		dbIndex++
		//log.Println(block, dbIndex)

		WriteBlock(block, dbIndex)
		for _, tx := range block.BlockTransactions {

			if tx == "97ddfbbae6be97fd6cdf3e7ca13232a3afff2353e29badfab7f73011edd4ced9" {

			} else {

				txHash, err := chainhash.NewHashFromStr(tx)
				if err != nil {
					log.Println(err)
				}
				txDaemon, err := client.GetRawTransactionVerbose(txHash)
				if err != nil {
					log.Println(err)
				}

				txIndexArray := block.Height
				txIndexArray++
				log.Println(tx)
				WriteTX(txDaemon, txIndexArray)
			}

		}
	}

	if endIndex < blockCount {
		newEndIndex := int64(0)
		if (endIndex + 500) > blockCount {
			newEndIndex = blockCount
		} else {
			newEndIndex = endIndex + 500
		}
		BlockRanger(client, endIndex+1, newEndIndex, blockCount)
	}

}

func WriteBlock(block Block, blockArrayIndex int64) {

	//cmd := exec.Command("chmod", "666", "blocks.db")
	db, err := storm.Open("blocks.db")

	if err != nil {
		log.Println("ERROR: Cannot open DB", err)
	}

	blockDB := dbBlock{
		ID:                blockArrayIndex,
		Hash:              block.Hash,
		Confirmations:     block.Confirmations,
		Size:              block.Size,
		StrippedSize:      block.StrippedSize,
		Weight:            block.Weight,
		Height:            block.Height,
		Version:           block.Version,
		VersionHex:        block.VersionHex,
		MerkleRoot:        block.MerkleRoot,
		BlockTransactions: block.BlockTransactions,
		Time:              block.Time,
		Nonce:             block.Nonce,
		Bits:              block.Bits,
		Difficulty:        block.Difficulty,
		PreviousHash:      block.PreviousHash,
		NextHash:          block.NextHash,
	}

	//initialise DB
	db.Init(&dbBlock{})
	// The block ranger array is every 500 blocks, Once we get empty blocks stop writing to DB
	if block.Hash == "" {
		db.Close()
		return
	}
	db.Save(&blockDB)
	db.Close()

}

func WriteTX(tx *btcjson.TxRawResult, txIndexArray int64) {

	//cmd := exec.Command("chmod", "666", "tx.db")
	db, err := storm.Open("tx.db")
	if err != nil {
		log.Println("ERROR: Cannot open TX DB", err)
	}

	txDB := dbTX{
		ID:            txIndexArray,
		Hex:           tx.Hex,
		Txid:          tx.Txid,
		Hash:          tx.Hash,
		Size:          tx.Size,
		Vsize:         tx.Vsize,
		Version:       tx.Version,
		LockTime:      tx.LockTime,
		Vin:           tx.Vin[:],
		Vout:          tx.Vout[:],
		BlockHash:     tx.BlockHash,
		Confirmations: tx.Confirmations,
		Time:          tx.Time,
		Blocktime:     tx.Blocktime,
	}
	log.Println("HELLO", txDB)

	//initialise DB
	db.Init(&dbTX{})
	// The block ranger array is every 500 blocks, Once we get empty blocks stop writing to DB
	if tx.Hash == "" {
		db.Close()
		return
	}
	db.Save(&txDB)
	db.Close()

}

func DBChecker() {

	if _, err := os.Stat("blocks.db"); os.IsNotExist(err) {
		log.Println("Running Block Ranger")
		GoBlockRanger()
	} else {
		log.Println("DB Exists, checking DB.......")
	}

}
