package api

import (
	"github.com/btcsuite/btcd/rpcclient"
	"log"
	"strconv"
)

// Block Structs

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

type TxRaw struct {
	Hex           string `json:"hex"`
	Txid          string `json:"txid"`
	Hash          string `json:"hash,omitempty"`
	Size          int32  `json:"size,omitempty"`
	Vsize         int32  `json:"vsize,omitempty"`
	Version       int32  `json:"version"`
	LockTime      uint32 `json:"locktime"`
	Vin           []Vin  `json:"vin"`
	Vout          []Vout `json:"vout"`
	BlockHash     string `json:"blockhash,omitempty"`
	Confirmations uint64 `json:"confirmations,omitempty"`
	Time          int64  `json:"time,omitempty"`
	Blocktime     int64  `json:"blocktime,omitempty"`
}

type Vin struct {
	Coinbase  string     `json:"coinbase"`
	Txid      string     `json:"txid"`
	Vout      uint32     `json:"vout"`
	ScriptSig *ScriptSig `json:"scriptSig"`
	Sequence  uint32     `json:"sequence"`
	Witness   []string   `json:"txinwitness"`
}

type Vout struct {
	Value        float64            `json:"value"`
	N            uint32             `json:"n"`
	ScriptPubKey ScriptPubKeyResult `json:"scriptPubKey"`
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


func BlockRanger (client *rpcclient.Client, startIndex int64, endIndex int64, blockCount int64) {

	//For each item (block) in the array index, print block details
	log.Println("starting block ranger at", strconv.FormatInt(startIndex,10) )

	logIndex := startIndex
	blockArray := make([]Block, 500)
	blockArrayIndex := 0
	for ; startIndex <= endIndex; startIndex++{
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
		blockArrayIndex = blockArrayIndex +1
	}
	log.Println("Batch completed from:", logIndex, " to ", endIndex)
	log.Println("Sending to DB")
	if (endIndex < blockCount) {
		newEndIndex := int64(0)
		if ((endIndex + 500) > blockCount) {
			newEndIndex = blockCount
		} else {
			newEndIndex = endIndex + 500
		}
		BlockRanger(client, endIndex + 1, newEndIndex, blockCount)
	}

}

