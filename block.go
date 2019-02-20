package capi

// Block provides the domain model for a blockchain block.
type Block struct {
	ID            string `json:"id"     storm:"id"`
	Height        int    `json:"height" storm:"index"`
	Hash          string `json:"hash"   storm:"index"`
	Confirmations int64  `json:"confirmations"`
	Size          int32  `json:"size"`
}

// TODO: Block methods go here...
