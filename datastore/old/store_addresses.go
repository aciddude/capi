package old

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/asdine/storm"

	"github.com/aciddude/capi/coind"
)

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
				TxInBlock:     tx.Blockhash,
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
