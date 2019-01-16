package coind

import (
	"encoding/json"
	"fmt"
)

// RawTx represents a raw transaction
type RawTransaction struct {
	Hex           string `json:"hex"`
	Txid          string `json:"txid"`
	Version       uint32 `json:"version"`
	LockTime      uint32 `json:"locktime"`
	Vin           []Vin  `json:"vin"`
	Vout          []Vout `json:"vout"`
	BlockHash     string `json:"blockhash,omitempty"`
	Confirmations int64  `json:"confirmations,omitempty"`
	Time          int64  `json:"time,omitempty"`
	Blocktime     int64  `json:"blocktime,omitempty"`
}

// Vin represents an IN value
type Vin struct {
	Coinbase  string    `json:"coinbase"`
	Txid      string    `json:"txid"`
	Vout      int       `json:"vout"`
	ScriptSig ScriptSig `json:"scriptSig"`
	Sequence  uint32    `json:"sequence"`
}

// Vout represents an OUT value
type Vout struct {
	Value        float64      `json:"value"`
	N            int          `json:"n"`
	ScriptPubKey ScriptPubKey `json:"scriptPubKey"`
}

// ScriptPubKey represents a scriptpubkey
type ScriptPubKey struct {
	Asm       string   `json:"asm"`
	Hex       string   `json:"hex"`
	ReqSigs   int      `json:"reqSigs,omitempty"`
	Type      string   `json:"type"`
	Addresses []string `json:"addresses,omitempty"`
}

// A ScriptSig represents a scriptsig
type ScriptSig struct {
	Asm string `json:"asm"`
	Hex string `json:"hex"`
}

// GetRawTransaction returns raw transaction representation for given transaction id.
func (d *Coind) GetRawTransaction(txId string, verbose bool) (rawTx interface{}, err error) {
	intVerbose := 0
	if verbose {
		intVerbose = 1
	}
	r, err := d.client.call("getrawtransaction", []interface{}{txId, intVerbose})
	if err = handleError(err, &r); err != nil {
		return
	}
	if !verbose {
		err = json.Unmarshal(r.Result, &rawTx)
	} else {
		var t RawTransaction
		err = json.Unmarshal(r.Result, &t)
		rawTx = t
	}
	return
}

// MakeRawTxListRequest returns array of block hights
//
// We need to get the list of hashes and then get the rawtransaction for each one

func (d *Coind) MakeRawTxListRequest(txhaslist []string) (response []rpcRequest, err error) {

	response = make([]rpcRequest, len(txhaslist))
	requestArrayIndex := 0
	ID := 0
	for requestArrayIndex < len(txhaslist) {
		for _, txhash := range txhaslist {
			response[ID] = rpcRequest{
				Id:      int64(ID),
				Method:  "getrawtransaction",
				Params:  []interface{}{txhash, 2},
				JsonRpc: "2.0",
			}
			if err != nil {
				fmt.Printf("ERROR! %s", err, txhaslist)
			}
			requestArrayIndex = requestArrayIndex + 1
			ID = ID + 1
		}

	}
	return response, err
}

func (d *Coind) GetRawTransactionList(params []rpcRequest) (response []rpcResponse, err error) {

	r, err := d.client.arraycall(params)
	if err = handleListError(err, &r); err != nil {
		return
	}

	return r, err

}

// Parserawtxlist takes a []byte, parses it and returns an array of []string as a list of block hashes

func ParseBlockTX(hashlist []byte) (list []string, err error) {

	// for each rawtx in the hashlist above marshal the json and remove the quotes, append the string to the empty string array above.
	// if the blockheight is 0 then remove that TX from the array as the coin daemon cannot decode it

	var blocks GetBlockResponse

	var genesistx string

	json.Unmarshal([]byte(hashlist), &blocks)

	for i, block := range blocks {
		if block.Result.Height == 0 {
			genesistx = block.Result.Tx[i]
		}
		for _, tx := range block.Result.Tx {
			list = append(list, tx)
			removeGenesisTX(list, genesistx)
		}
	}

	return list, err
}

func removeGenesisTX(txlist []string, genesistx string) []string {
	for i, v := range txlist {
		if v == genesistx {
			return append(txlist[:i], txlist[i+1:]...)
		}
	}
	return txlist
}
