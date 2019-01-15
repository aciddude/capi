package datastore

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/aciddude/capi/coind"
	"github.com/asdine/storm"
)

type DataStoreAddressSchema struct {
	ID          int      `storm:"id,increment" json:"-"` // primary key
	Address     []string `storm:"index"`
	Transaction string   `storm:"unique"`
}

func StoreAddresses() {

	daemonConfig := coind.LoadConfig("./config/config.json")
	coinDaemon, err := coind.New(daemonConfig, daemonConfig.RPCTimeout)
	if err != nil {
		fmt.Printf("DATA STORE ERROR: Cannot connect to coin daemon \n %v", err)
	}

	startHeight := 0
	endHeight := 100
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

	db, err := storm.Open("test.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var transaction coind.RawTransaction

	for _, result := range rawtxlist {

		json.Unmarshal(result.Result, &transaction)

		for _, add := range transaction.Vout {

			DBAddr := DataStoreAddressSchema{
				Address:     add.ScriptPubKey.Addresses,
				Transaction: transaction.Txid,
			}
			fmt.Printf("WRITING TO DB, %v \n", DBAddr)
			err := db.Save(&DBAddr)
			if err != nil {
				fmt.Errorf("could not save config, %v", err)
			}
		}

	}
	db.Close()
}
