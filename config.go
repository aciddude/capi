package capi

type datastore string

const (
	BoltDB datastore = "boltdb"
)

// Config provides the domain structure to enable configuring CAPI.
type Config struct {
	// Port configures the API port CAPI will serve from.
	Port string

	// Coin Configuration.
	Coins []Coin

	// Datastore configuration.
	Datastore datastore
	// BoltDB specific datastore configuration.
	BoltDB ConfigBoltDB
}

// Coin provides the configuration required to connect to a coin daemon API.
type Coin struct {
	// Name is the human readable name of the coin. For example, "Feathercoin".
	Name string
	// Code is the 3 letter coin code. For example, "FTC".
	Code string
	// Host is the Coin's API daemon hostname on which to connect to the API.
	Host string
	// Port is the Coin's API daemon port on which to connect to the API.
	Port string
	// Username is the username to use in order to authenticate to the Coin's
	// API daemon.
	Username string
	// Password is the password to use in order to authenticate to the Coin's
	// API daemon.
	Password string
	// Timeout is how long to wait before timing out API requests.
	Timeout int
	// SSL is whether to connect over SSL.
	// If not specified, the default is false.
	SSL bool
	// EnableCoinCodexAPI is whether to enable the coin's codex API.
	// If not specified, the default is false.
	EnableCoinCodexAPI bool
}

// ConfigBoltDB provides specific configuration customisation for BoltDB.
type ConfigBoltDB struct {
	// How long to wait before timing out a query.
	Timeout int
}
