package old

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/asdine/storm"

	"github.com/aciddude/capi/coind"
)

// delcare start and end height of the scan....
// to do.  Do a getblockcount to get current height and assign that as end height

var startHeight int = 0
var endHeight int = 500

// Check the DB height vs Daemon height and return both
func DBvsDaemon() (dbHeight, daemonHeight int) {

	/// First we get the current block height of the daemon.
	daemonConfig := coind.LoadConfig("./config/config.json")
	coinDaemon, err := coind.New(daemonConfig, daemonConfig.RPCTimeout)
	if err != nil {
		fmt.Printf("DBvsDaemon could not connect to daemon %v:", err)
	}

	coindheight, err := coinDaemon.GetBlockCount()
	if err != nil {
		fmt.Printf("DBvsDaemon could not get block height %v:", err)
	}

	// now we have the current height blockchain from the daemon as currentHeight

	// Second we get the height of the DB

	db, err := storm.Open("blocks.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var block []BlockDB

	// search the db for all blocks, limit to 1 block returned and reverse to get the last block in the db
	err = db.All(&block, storm.Limit(1), storm.Reverse())
	fmt.Println(err)

	for range block {
		if err == nil && block[0].Height == 0 {
			dbHeight = 0
			fmt.Println(err)
		} else {
			dbHeight = int(block[0].Height)
		}
	}

	daemonHeight = int(coindheight)

	// print daemon and db height
	fmt.Printf("Daemon Height: %v \n", coindheight)
	fmt.Printf("DB Height: %v \n", dbHeight)

	return dbHeight, daemonHeight
}

/// check to see if the blocks DB exists
func DatabaseExists() bool {

	if _, err := os.Stat("blocks.db"); os.IsNotExist(err) {
		log.Println("Running Datastore")
		return false

	} else {
		log.Println("DB Exists, checking DB.......")
	}

	return true

}

func BlockRanger() {

	daemonConfig := coind.LoadConfig("./config/config.json")
	coinDaemon, err := coind.New(daemonConfig, daemonConfig.RPCTimeout)
	if err != nil {
		fmt.Printf("DBvsDaemon could not connect to daemon %v:", err)
	}

	localHeight, DaemonHeight := DBvsDaemon()

	blocks := make([]int, DaemonHeight)
	batch := 500

	for index := localHeight; index < len(blocks); index += batch {

		endindex := index + batch
		if endindex >= len(blocks) {
			endindex = len(blocks) - 1
		}

		listsize := endindex - index

		fmt.Printf("Storing txns from Height [%v to %v] \n", index, endindex)
		fmt.Printf("Storing Blocks from Height [%v to %v] \n", index, endindex)
		fmt.Printf("Storing Addresses from Height [%v to %v] \n", index, endindex)

		blockhashes, err := coind.MakeBlockHashListRequest(index, endindex)
		if err != nil {
			fmt.Printf("DBvsDaemon could not create block hash list request %v:", err)
		}
		getblockreponse, err := coinDaemon.GetBlockHashList(blockhashes)
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

		txdb, err := storm.Open("transactions.db", storm.Batch())
		if err != nil {
			log.Fatal(err)
		}
		defer txdb.Close()

		var TXData TransactionDB
		var tx coind.RawTransaction

		for _, result := range rawtxlist {

			json.Unmarshal(result.Result, &tx)

			TXData = TransactionDB{
				Hex:           tx.Hex,
				Txid:          tx.Txid,
				Hash:          tx.Hash,
				Size:          tx.Size,
				Vsize:         tx.Vsize,
				Version:       tx.Version,
				LockTime:      tx.Locktime,
				Vin:           tx.Vin,
				Vout:          tx.Vout,
				BlockHash:     tx.Blockhash,
				Confirmations: tx.Confirmations,
				Time:          tx.Time,
				Blocktime:     tx.Blocktime,
			}
			err = txdb.Save(&TXData)
			if err != nil {
				fmt.Errorf("could not save config, %v", err)
			}

		}

		var blocks coind.GetBlockResponse

		json.Unmarshal([]byte(jsonblocklist), &blocks)

		blockdb, err := storm.Open("blocks.db", storm.Batch())
		if err != nil {
			log.Fatal(err)
		}
		defer blockdb.Close()

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
			err = blockdb.Init(&Data)
			if err != nil {
				fmt.Errorf("could not save blocks to db, %v", err)
			}

			err = blockdb.Save(&Data)
			if err != nil {
				fmt.Errorf("could not save blocks to db, %v", err)
			}
		}

		addrdb, err := storm.Open("addresses.db", storm.Batch())
		if err != nil {
			log.Fatal(err)
		}
		defer addrdb.Close()

		var AddrData AddressDB

		for _, result := range rawtxlist {

			json.Unmarshal(result.Result, &tx)

			for i, vout := range tx.Vout {
				//fmt.Printf("Parsing Addresses and Transactions to DB...... %v \n", vout.ScriptPubKey.Addresses)
				AddrData = AddressDB{
					TxID:          tx.Txid,
					Address:       vout.ScriptPubKey.Addresses,
					Received:      tx.Vout[i].Value,
					Confirmations: tx.Confirmations,
					TxInBlock:     tx.Blockhash,
					TxTime:        tx.Time,
				}
				err = addrdb.Save(&AddrData)
				if err != nil {
					fmt.Errorf("could not save config, %v", err)
				}

			}

		}
		blockdb.Close()
		txdb.Close()
		addrdb.Close()

	}

}
