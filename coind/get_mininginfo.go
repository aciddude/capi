package coind

import "encoding/json"

//GetMiningInfo Struct

type GetMiningInfo struct {
	Blocks             int     `json:"blocks"`
	Currentblockweight int     `json:"currentblockweight"`
	Currentblocktx     int     `json:"currentblocktx"`
	Difficulty         float64 `json:"difficulty"`
	Networkhashps      float64 `json:"networkhashps"`
	Pooledtx           int     `json:"pooledtx"`
	Chain              string  `json:"chain"`
	Warnings           string  `json:"warnings"`
}

// RPC Command: getmininginfo  - Returns a GetMiningInfo struct

func (d *Coind) GetMiningInfo() (i GetMiningInfo, err error) {
	r, err := d.client.call("getmininginfo", nil)
	if err = handleError(err, &r); err != nil {
		return
	}
	err = json.Unmarshal(r.Result, &i)

	return i, err
}
