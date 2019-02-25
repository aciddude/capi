// dev provides the development build, which may have unimplemented features,
// or dragons.. Who knows...
package main

import (
	// Standard Library Imports
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	// External Imports
	log "github.com/sirupsen/logrus"

	// Internal Imports
	"github.com/aciddude/capi"
	"github.com/aciddude/capi/coind"
	"github.com/aciddude/capi/datastore"
)

const (
	// chunkSize si how many blocks to request at once.
	chunkSize = 1000
)

var (
	configPath = "../../config/config.yaml"
)

func init() {
	log.SetLevel(log.DebugLevel)
}

func updateCoin(store *datastore.Datastore, coin capi.Coin) (err error) {
	logger := log.WithFields(log.Fields{
		"app":  "capi-dev",
		"func": "updateCoin",
		"coin": coin.Code,
	})

	cfg := coind.Config{
		Coin:        coin.Name,
		Ticker:      coin.Code,
		RPCHost:     coin.Host,
		RPCUser:     coin.Username,
		RPCPassword: coin.Password,
		RPCPORT:     strconv.Itoa(int(coin.Port)),
		RPCTimeout:  coin.Timeout,
		SSL:         coin.SSL,
	}
	coinDaemon, err := coind.New(cfg, coin.Timeout)
	if err != nil {
		logger.WithError(err).Fatal()
	}

	// TODO: create an interface to deal with returning a datastore sesison for stores that require it.
	ctx := context.TODO()
	var startHeight int
	lastStored, err := store.Blocks.Last(ctx, coin.Code)
	if err != nil {
		switch err {
		case datastore.ErrNotFound:
			logger.WithError(err).Debug("last block not found, starting sync from 0")
			// last block not found in the datastore, therefore, start at zero.
			startHeight = 0

		default:
			logger.WithError(err).Error("error getting last block")
			return err
		}
	}

	if lastStored != nil {
		startHeight = lastStored.Height + 1
	}

	blockCount, err := coinDaemon.GetBlockCount()
	if err != nil {
		logger.WithError(err).Error("error getting current block height")
		return err
	}
	currentBlockHeight := int(blockCount)

	for startHeight < currentBlockHeight {
		endHeight := chunkSize + startHeight
		if endHeight > currentBlockHeight {
			endHeight = currentBlockHeight
		}

		blocks, err := getBlocks(coinDaemon, coin.Code, startHeight, endHeight)
		if err != nil {
			logger.WithError(err).Error("error getting blocks")
			return err
		}

		blocks, err = store.Blocks.CreateBulk(ctx, coin.Code, blocks)
		if err != nil {
			logger.WithError(err).Error("error storing blocks")
			return err
		}

		startHeight = endHeight

		for _, block := range blocks {
			fmt.Printf("Block %d hash: %s\n", block.Height, block.Hash)
		}
	}

	return err
}

// getBlocks uses a coin daemon to get a list of blocks, given a start and
// end height.
func getBlocks(client *coind.Coind, coin string, startHeight, endHeight int) (blocks []*capi.Block, err error) {
	logger := log.WithFields(log.Fields{
		"app":  "capi-dev",
		"func": "getBlocks",
		"coin": coin,
	})

	getblockhashrequest, err := coind.MakeBlockHashListRequest(startHeight, endHeight)
	if err != nil {
		logger.WithError(err).Debug("MakeBlockHashListRequest error")
		return nil, err
	}

	getblockreponse, err := client.GetBlockHashList(getblockhashrequest)
	if err != nil {
		logger.WithError(err).Debug("GetBlockHashList error")
		return nil, err
	}

	hashlist, err := coind.ParseBlockHashList(getblockreponse)
	if err != nil {
		logger.WithError(err).Debug("ParseBlockHashList error")
		return nil, err
	}

	listsize := endHeight - startHeight
	getblockrequest, err := client.MakeGetBlockListRequest(listsize, hashlist)
	if err != nil {
		logger.WithError(err).Debug("MakeGetBlockListRequest error")
		return nil, err
	}

	blocklist, err := client.GetBlockList(getblockrequest)
	if err != nil {
		logger.WithError(err).Debug("GetBlockList error")
		return nil, err
	}

	for _, blockRes := range blocklist {
		newBlock := &capi.Block{}
		err := json.Unmarshal(blockRes.Result, newBlock)
		if err != nil {
			logger.WithError(err).Error("error parsing block")
			return nil, err
		}

		blocks = append(blocks, newBlock)
	}

	return blocks, nil
}

func main() {
	logger := log.WithFields(log.Fields{
		"app": "capi-dev",
	})

	config, err := capi.NewConfig(configPath)
	if err != nil {
		logger.WithError(err).Error("error processing config")
		os.Exit(1)
	}

	store, err := datastore.NewDatastore(config)
	if err != nil {
		logger.WithError(err).Error("error starting datastore")
		os.Exit(1)
	}

	for _, coin := range config.Coins {
		err = updateCoin(store, coin)
		if err != nil {
			logger.WithError(err).Error()
		}
	}

	// TODO: Bind store into service.
	// TODO: configure API endpoints based on config.

	err = store.Close()
	if err != nil {
		logger.WithError(err).Error("error closing datastore connections")
		os.Exit(1)
	}

	// Close happily
	os.Exit(0)
}
