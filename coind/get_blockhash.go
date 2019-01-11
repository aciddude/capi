package coind

import (
	"encoding/json"
	"fmt"
)

// GetBlockHash returns hash of block in best-block-chain at block <height>
func (d *Coind) GetBlockHash(height int64) (hash string, err error) {
	r, err := d.client.call("getblockhash", []int64{height})
	if err = handleError(err, &r); err != nil {
		return
	}
	err = json.Unmarshal(r.Result, &hash)
	return
}

// MakeBlockHashListRequest returns array of requests for block hashes
func MakeBlockHashListRequest(heightStart, heightEnd int) (hashArray []rpcRequest, err error) {

	var arraySize int
	arraySize = heightEnd - heightStart
	hashArray = make([]rpcRequest, arraySize)

	requestArrayIndex := 0
	ID := int64(0)
	for ; requestArrayIndex < arraySize; heightStart++ {
		hashArray[requestArrayIndex] = rpcRequest{
			Id:      ID,
			Method:  "getblockhash",
			Params:  []interface{}{heightStart},
			JsonRpc: "2.0",
		}
		if err != nil {
			fmt.Printf("ERROR! %s", err)
		}
		ID = ID + 1
		requestArrayIndex = requestArrayIndex + 1
	}

	//fmt.Printf("HASH ARRAY \n %d", heightEnd-heightStart)

	return hashArray, err
}

// GetBlockHashList sends the array of requests to the daemon and get a response back
func (d *Coind) GetBlockHashList(params []rpcRequest) (response []rpcResponse, err error) {

	r, err := d.client.arraycall(params)
	if err = handleListError(err, &r); err != nil {
		return
	}

	return r, nil

}
