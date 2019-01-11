package coind

import "encoding/json"

type GetNetworkInfo struct {
	Version         int    `json:"version"`
	Subversion      string `json:"subversion"`
	Protocolversion int    `json:"protocolversion"`
	Localservices   string `json:"localservices"`
	Localrelay      bool   `json:"localrelay"`
	Timeoffset      int    `json:"timeoffset"`
	Networkactive   bool   `json:"networkactive"`
	Connections     int    `json:"connections"`
	Networks        []struct {
		Name                      string `json:"name"`
		Limited                   bool   `json:"limited"`
		Reachable                 bool   `json:"reachable"`
		Proxy                     string `json:"proxy"`
		ProxyRandomizeCredentials bool   `json:"proxy_randomize_credentials"`
	} `json:"networks"`
	Relayfee       float64 `json:"relayfee"`
	Incrementalfee float64 `json:"incrementalfee"`
	Localaddresses []struct {
		Address string `json:"address"`
		Port    int    `json:"port"`
		Score   int    `json:"score"`
	} `json:"localaddresses"`
	Warnings string `json:"warnings"`
}

type GetNetworkInfoLocalAddresses []struct {
	Address string `json:"address"`
	Port    int    `json:"port"`
	Score   int    `json:"score"`
}

// RPC Command: getnetworkinfo  - Returns a GetNetworkInfo struct

func (d *Coind) GetNetworkInfo() (i GetNetworkInfo, err error) {
	r, err := d.client.call("getnetworkinfo", nil)
	if err = handleError(err, &r); err != nil {
		return
	}
	err = json.Unmarshal(r.Result, &i)

	return
}
