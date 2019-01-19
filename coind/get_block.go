package coind

import (
	"encoding/json"
	"fmt"
	"strings"
)

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

type GetBlockResponse []struct {
	ID     int `json:"id"`
	Result struct {
		Hash          string   `json:"hash"`
		Confirmations int      `json:"confirmations"`
		Strippedsize  int      `json:"strippedsize"`
		Size          int      `json:"size"`
		Weight        int      `json:"weight"`
		Height        int      `json:"height"`
		Version       int      `json:"version"`
		VersionHex    string   `json:"versionHex"`
		Merkleroot    string   `json:"merkleroot"`
		Tx            []string `json:"tx"`
		Time          int      `json:"time"`
		Mediantime    int      `json:"mediantime"`
		Nonce         int      `json:"nonce"`
		Bits          string   `json:"bits"`
		Difficulty    float64  `json:"difficulty"`
		Chainwork     string   `json:"chainwork"`
		Nextblockhash string   `json:"nextblockhash"`
	} `json:"result"`
	Error interface{} `json:"error"`
}

// GetBlock returns json  Block from block hash
func (d *Coind) GetBlock(hash string) (block Block, err error) {
	r, err := d.client.call("getblock", []string{hash})
	if err = handleError(err, &r); err != nil {
		return
	}
	err = json.Unmarshal(r.Result, &block)
	return
}

func (d *Coind) MakeGetBlockListRequest(listsize int, hashlist []string) (response []rpcRequest, err error) {

	response = make([]rpcRequest, listsize)
	requestArrayIndex := 0
	ID := 0
	for requestArrayIndex < listsize {
		for range hashlist {
			response[requestArrayIndex] = rpcRequest{
				Id: int64(ID),
				//		Err:    []interface{}{err},
				Method: "getblock",
				Params: []interface{}{hashlist[requestArrayIndex], true},
			}
			if err != nil {
				fmt.Printf("ERROR! %s %s", err, hashlist)
			}

		}
		requestArrayIndex = requestArrayIndex + 1
		ID = ID + 1
	}

	//ID = ID + 1
	//fmt.Print(hashArray)

	return response, err

}

func (d *Coind) GetBlockList(params []rpcRequest) (response []rpcResponse, err error) {

	r, err := d.client.arraycall(params)
	if err = handleListError(err, &r); err != nil {
		return
	}

	return r, nil

}

// Parsehashlist takes the []rpcReponse, parses it and returns an array of []string as a list of block hashes

func ParseBlockHashList(hashlist []rpcResponse) (list []string, err error) {

	// for each hash in the hashlist above marshal the json and remove the quotes, append the string to the empty string array above.
	for _, hash := range hashlist {

		jsonhashlist, err := hash.Result.MarshalJSON()
		if err != nil {
			fmt.Errorf("ERROR", err)
		}
		hash := string(jsonhashlist)
		hash = strings.Replace(hash, `"`, "", -1)

		list = append(list, hash)
	}
	return list, err

}
