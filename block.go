package capi

// Block provides the domain model for a blockchain block.
type Block struct {
	ID            string `json:"id" storm:"id"`
	Hash          string `json:"hash"`
	Confirmations int64  `json:"confirmations"`
	Size          int32  `json:"size"`
}

// TODO: Block methods go here...
