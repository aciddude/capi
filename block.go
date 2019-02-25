package capi

// Block provides the domain model for a blockchain block.
type Block struct {
	ID                string   `json:"id" storm:"id"`
	Hash              string   `json:"hash" storm:"index"`
	Confirmations     int      `json:"confirmations"`
	StrippedSize      int      `json:"strippedsize"`
	Size              int      `json:"size"`
	Weight            int      `json:"weight"`
	Height            int      `json:"height" storm:"index"`
	Version           int      `json:"version"`
	VersionHex        string   `json:"versionHex"`
	MerkleRoot        string   `json:"merkleroot"`
	BlockTransactions []string `json:"tx"`
	Time              int64    `json:"time"`
	Mediantime        int      `json:"mediantime"`
	Nonce             uint32   `json:"nonce"`
	Bits              string   `json:"bits"`
	Difficulty        float64  `json:"difficulty"`
	Chainwork         string   `json:"chainwork"`
	PreviousHash      string   `json:"previousBlockHash"`
	NextHash          string   `json:"nextBlockHash"`
}

// TODO: Block methods go here...
