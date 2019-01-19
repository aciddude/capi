package datastore

import (
	"capi/coind"
	"encoding/json"
	"fmt"
	"log"

	"github.com/asdine/storm"
)

func StoreBlocks() {

	daemonConfig := coind.LoadConfig("./config/config.json")
	coinDaemon, err := coind.New(daemonConfig, daemonConfig.RPCTimeout)
	if err != nil {
		fmt.Printf("DATA STORE ERROR: Cannot connect to coin daemon \n %v", err)
	}

	listsize := endHeight - startHeight

	getblockhashrequest, err := coind.MakeBlockHashListRequest(startHeight, endHeight)
	if err != nil {
		fmt.Printf("DATA STORE ERROR: Cannot create getblockhash list request \n %v", err)
	}

	getblockreponse, err := coinDaemon.GetBlockHashList(getblockhashrequest)
	if err != nil {
		fmt.Printf("DATA STORE ERROR: Cannot get block hash list \n %v", err)
	}

	hashlist, err := coind.ParseBlockHashList(getblockreponse)
	if err != nil {
		fmt.Printf("DATA STORE ERROR: Cannot parse the block hash list \n %v", err)
	}

	getblockrequest, err := coinDaemon.MakeGetBlockListRequest(listsize, hashlist)
	if err != nil {
		fmt.Printf("DATA STORE ERROR: Cannot create getblock list request \n %v", err)
	}

	blocklist, err := coinDaemon.GetBlockList(getblockrequest)
	if err != nil {
		fmt.Printf("DATA STORE ERROR: Cannot get block list \n %v", err)

	}
	jsonblocklist, _ := json.Marshal(blocklist)

	var blocks coind.GetBlockResponse

	json.Unmarshal([]byte(jsonblocklist), &blocks)

	db, err := storm.Open("blocks.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var Data BlockDB
	var block coind.Block

	for _, result := range blocklist {

		json.Unmarshal(result.Result, &block)

		Data = BlockDB{

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

		err = db.Save(&Data)
		if err != nil {
			fmt.Errorf("could not save config, %v", err)
		}

	}
	db.Close()
	log.Println("Finished storing blocks")
}
