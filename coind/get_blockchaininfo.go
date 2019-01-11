package coind

import (
	"encoding/json"
)

type GetBlockchainInfo struct {
	Chain                string  `json:"chain"`
	Blocks               int     `json:"blocks"`
	Headers              int     `json:"headers"`
	Bestblockhash        string  `json:"bestblockhash"`
	Difficulty           float64 `json:"difficulty"`
	Mediantime           int     `json:"mediantime"`
	Verificationprogress float64 `json:"verificationprogress"`
	Initialblockdownload bool    `json:"initialblockdownload"`
	Chainwork            string  `json:"chainwork"`
	SizeOnDisk           int64   `json:"size_on_disk"`
	Pruned               bool    `json:"pruned"`
	Softforks            []struct {
		ID      string `json:"id"`
		Version int    `json:"version"`
		Reject  struct {
			Status bool `json:"status"`
		} `json:"reject"`
	} `json:"softforks"`
	Bip9Softforks struct {
		Csv struct {
			Status    string `json:"status"`
			StartTime int    `json:"startTime"`
			Timeout   int    `json:"timeout"`
			Since     int    `json:"since"`
		} `json:"csv"`
		Segwit struct {
			Status    string `json:"status"`
			StartTime int    `json:"startTime"`
			Timeout   int    `json:"timeout"`
			Since     int    `json:"since"`
		} `json:"segwit"`
	} `json:"bip9_softforks"`
	Warnings string `json:"warnings"`
}

// RPC Command: getblockchaininfo - Returns a GetBlockChainInfo struct

func (d *Coind) GetBlockchainInfo() (i GetBlockchainInfo, err error) {
	r, err := d.client.call("getblockchaininfo", nil)
	if err = handleError(err, &r); err != nil {
		return
	}
	err = json.Unmarshal(r.Result, &i)

	return
}
