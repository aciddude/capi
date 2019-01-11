package coind

import "encoding/json"

//GetInfo Struct

type GetInfo struct {
	Balance         float64 `json:"balance"`
	Blocks          int     `json:"blocks"`
	Connections     int     `json:"connections"`
	Difficulty      float64 `json:"difficulty"`
	Errors          string  `json:"errors"`
	Keypoololdest   int     `json:"keypoololdest"`
	Keypoolsize     int     `json:"keypoolsize"`
	Paytxfee        float64 `json:"paytxfee"`
	Protocolversion int     `json:"protocolversion"`
	Proxy           string  `json:"proxy"`
	Relayfee        float64 `json:"relayfee"`
	Testnet         bool    `json:"testnet"`
	Timeoffset      int     `json:"timeoffset"`
	Version         int     `json:"version"`
	Walletversion   int     `json:"walletversion"`
}

// RPC Command: getinfo  - Returns a GetInfo struct

func (d *Coind) GetInfo() (i GetInfo, err error) {
	r, err := d.client.call("getinfo", nil)
	if err = handleError(err, &r); err != nil {
		return
	}
	err = json.Unmarshal(r.Result, &i)

	return
}
