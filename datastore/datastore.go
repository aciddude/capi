package datastore

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/aciddude/capi/coind"

	"github.com/asdine/storm"
)

// delcare start and end hight of the scan....
// to do.  Do a getblockcount to get current hight and assign that as end hight

var startHeight int = 0
var endHeight int = 8000

func DatabaseExists() bool {

	if _, err := os.Stat("blocks.db"); os.IsNotExist(err) {
		log.Println("Running Datastore")
		return false

	} else {
		log.Println("DB Exists, checking DB.......")
	}
	return true
}

func StoreAddresses() {

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
	txlist, err := coind.ParseBlockTX(jsonblocklist)
	if err != nil {
		fmt.Printf("DATA STORE ERROR: Cannot parse transactions from block list \n %v", err)

	}

	getrawtxrequest, err := coinDaemon.MakeRawTxListRequest(txlist)
	if err != nil {
		fmt.Printf("DATA STORE ERROR: Cannot create getrawtransaction list request \n %v", err)
	}

	rawtxlist, err := coinDaemon.GetRawTransactionList(getrawtxrequest)
	if err != nil {
		fmt.Printf("ERROR:\nRaw Transacaction List Request %v ", err)
	}

	db, err := storm.Open("addresses.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var Data AddressDB
	var tx coind.RawTransaction

	for _, result := range rawtxlist {

		json.Unmarshal(result.Result, &tx)

		for i, vout := range tx.Vout {
			//fmt.Printf("Parsing Addresses and Transactions to DB...... %v \n", vout.ScriptPubKey.Addresses)
			Data = AddressDB{
				TxID:          tx.Txid,
				Address:       vout.ScriptPubKey.Addresses,
				Received:      tx.Vout[i].Value,
				Confirmations: tx.Confirmations,
				TxInBlock:     tx.BlockHash,
				TxTime:        tx.Time,
			}
			err = db.Save(&Data)
			if err != nil {
				fmt.Errorf("could not save config, %v", err)
			}

		}

	}
	db.Close()
	log.Println("Finished storing addresses")
}

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
