package datastore

import "github.com/aciddude/capi/coind"

type AddressDB struct {
	ID            int      `storm:"id,increment" json:"-"` // Primary Key
	TxID          string   `storm:"index" json:"txid"`
	Address       []string `storm:"index" json:"address"`
	Received      float64  `json:"received"`
	Confirmations int64    `json:"confirmations"`
	TxInBlock     string   `json:"block_hash"`
	TxTime        int64    `json:"tx_time"`
}

// DB Block Struct

type BlockDB struct {
	ID                int64  // The Primary Key
	Hash              string `storm:"index"`
	Confirmations     int64
	Size              int32
	StrippedSize      int32
	Weight            int32
	Height            int64 `storm:"index"`
	Version           int32
	VersionHex        string
	MerkleRoot        string
	BlockTransactions []string
	Time              int64
	Nonce             uint32
	Bits              string
	Difficulty        float64
	PreviousHash      string
	NextHash          string
}

type TransactionDB struct {
	ID            int64        `storm:"id,increment" json:"-"` // The Primary Key
	Hex           string       `json:"hex"`
	Txid          string       `storm:"index"`
	Hash          string       `json:"hash,omitempty"`
	Size          int64        `json:"size,omitempty"`
	Vsize         int64        `json:"vsize,omitempty"`
	Version       int64        `json:"version"`
	LockTime      int64        `json:"locktime"`
	Vin           []coind.Vin  `storm:"inline"`
	Vout          []coind.Vout `storm:"inline"`
	BlockHash     string       `json:"blockhash,omitempty"`
	Confirmations int64        `json:"confirmations,omitempty"`
	Time          int64        `json:"time,omitempty"`
	Blocktime     int64        `json:"blocktime,omitempty"`
}
