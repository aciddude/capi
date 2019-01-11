package coind

// a Coind is a client for the coin daemon you're using
type Coind struct {
	client *rpcClient
}

// New return a new bitcoind
func New(config Config, timeoutParam ...int) (*Coind, error) {
	var timeout = config.RPCTimeout
	// If the timeout is specified in timeoutParam, allow it.
	if len(timeoutParam) != 0 {
		timeout = timeoutParam[0]
	}

	rpcClient, err := newClient(config, timeout)
	if err != nil {
		return nil, err
	}
	return &Coind{client: rpcClient}, nil
}
