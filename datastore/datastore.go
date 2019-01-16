package datastore

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/aciddude/capi/coind"

	"github.com/asdine/storm"
)

type AddressDB struct {
	ID            int      `storm:"id,increment" json:"-"` // primary key
	TxID          string   `storm:"index" json:"txid"`
	Address       []string `storm:"index" json:"address"`
	Received      float64  `json:"received"`
	Confirmations int64    `json:"confirmations"`
	TxInBlock     string   `json:"block_hash"`
	TxTime        int64    `json:"tx_time"`
}

func StoreAddresses() {

	daemonConfig := coind.LoadConfig("./config/config.json")
	coinDaemon, err := coind.New(daemonConfig, daemonConfig.RPCTimeout)
	if err != nil {
		fmt.Printf("DATA STORE ERROR: Cannot connect to coin daemon \n %v", err)
	}

	startHeight := 0
	endHeight := 1000 // parse 1K blocks + transactions + addresses + ammount received, then store in DB
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
			fmt.Printf("Parsing Addresses and Transactions to DB...... %v \n", vout.ScriptPubKey.Addresses)
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
}
