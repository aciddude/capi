package coind

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// A rpcClient represents a JSON RPC client (over HTTP(s)).
type rpcClient struct {
	serverAddr string
	user       string
	passwd     string
	httpClient *http.Client
	timeout    int
}

// rpcRequest represent a RCP request
type rpcRequest struct {
	Id      int64       `json:"id"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
	JsonRpc string      `json:"jsonrpc"`
}

// rpcResponse is the response given to the RPC request
type rpcResponse struct {
	Id     int64           `json:"id"`
	Result json.RawMessage `json:"result"`
	Err    interface{}     `json:"error"`
}

//// rpcResponse is the response given to the RPC request
//type rpcListResponse struct {
//	Id     int             `json:"id"`
//	Result json.RawMessage `json:"result"`
//	Err    interface{}     `json:"error"`
//}

// Create a new RPC client
// A New_Client accepts a config file
//
// Returns an *RPC_client and/or error

func newClient(config Config, timeout int) (client *rpcClient, err error) {

	if len(config.RPCHost) == 0 {
		err = errors.New("E0001: Missing Host in client config")
	}

	var serverAddress string
	var httpClient *http.Client

	if config.SSL {
		serverAddress = "https://"
		t := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		httpClient = &http.Client{Transport: t}
	} else {
		serverAddress = "http://"
		httpClient = &http.Client{}
	}
	client = &rpcClient{
		serverAddr: serverAddress + config.RPCHost + ":" + config.RPCPORT,
		user:       config.RPCUser,
		passwd:     config.RPCPassword,
		httpClient: httpClient,
		timeout:    config.RPCTimeout,
	}

	return
}

// doTimeoutRequest process a HTTP request with timeout
func (c *rpcClient) doTimeoutRequest(timer *time.Timer, req *http.Request) (*http.Response, error) {
	type result struct {
		resp *http.Response
		err  error
	}
	done := make(chan result, 1)
	go func() {
		resp, err := c.httpClient.Do(req)
		done <- result{resp, err}
	}()
	// Wait for the read or the timeout
	select {
	case r := <-done:
		return r.resp, r.err
	case <-timer.C:
		return nil, errors.New("Timeout reading data from server")
	}
}

//  Prepare the request and call the daemon

func (c *rpcClient) call(method string, params interface{}) (rr rpcResponse, err error) {
	connectTimer := time.NewTimer(time.Duration(c.timeout) * time.Second)
	rpcR := rpcRequest{time.Now().UnixNano(), method, params, "1.0"}
	payloadBuffer := &bytes.Buffer{}
	jsonEncoder := json.NewEncoder(payloadBuffer)
	err = jsonEncoder.Encode(rpcR)
	if err != nil {
		return
	}
	req, err := http.NewRequest("POST", c.serverAddr, payloadBuffer)
	if err != nil {
		return
	}
	req.Header.Add("Content-Type", "application/json;charset=utf-8")
	req.Header.Add("Accept", "application/json")

	// Auth ?
	if len(c.user) > 0 || len(c.passwd) > 0 {
		req.SetBasicAuth(c.user, c.passwd)
	}

	resp, err := c.doTimeoutRequest(connectTimer, req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	//fmt.Println(string(data))
	if err != nil {
		return
	}
	if resp.StatusCode != 200 {
		err = errors.New("HTTP error: " + resp.Status)
		return
	}
	err = json.Unmarshal(data, &rr)
	return
}

/// arraycall Request takes an array of []rpcRequest and gives back an array of []rpcResponse used for getting a range of block hashes or blocks

func (c *rpcClient) arraycall(params []rpcRequest) (rr []rpcResponse, err error) {
	connectTimer := time.NewTimer(time.Duration(c.timeout) * time.Second)
	payloadBuffer := &bytes.Buffer{}

	jsonEncoder := json.NewEncoder(payloadBuffer)
	err = jsonEncoder.Encode(params)
	if err != nil {
		return
	}
	req, err := http.NewRequest("POST", c.serverAddr, payloadBuffer)
	if err != nil {
		return
	}
	req.Header.Add("Content-Type", "application/json;charset=utf-8")
	req.Header.Add("Accept", "application/json")

	// Auth ?
	if len(c.user) > 0 || len(c.passwd) > 0 {
		req.SetBasicAuth(c.user, c.passwd)
	}

	resp, err := c.doTimeoutRequest(connectTimer, req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	//fmt.Printf("%s", data)
	if err != nil {
		return
	}
	if resp.StatusCode != 200 {
		fmt.Printf("Request that Errored %v", req)
		err = errors.New("HTTP error: " + resp.Status)
		return
	}
	err = json.Unmarshal(data, &rr)
	return
}

// handleListError handle error returned by client.blocklist &  client.hashlist
func handleListError(err error, r *[]rpcResponse) error {
	if err != nil {
		return err
	}
	for _, r := range *r {
		if r.Err != nil {
			rr := r.Err.(map[string]interface{})
			return errors.New(fmt.Sprintf("(%v) %s", rr["code"].(float64), rr["message"].(string)))
		}
	}
	return nil
}

// handleError handle error returned by client.call
func handleError(err error, r *rpcResponse) error {
	if err != nil {
		return err
	}
	if r.Err != nil {
		rr := r.Err.(map[string]interface{})
		return errors.New(fmt.Sprintf("(%v) %s", rr["code"].(float64), rr["message"].(string)))

	}
	return nil
}
