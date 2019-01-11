package coind

import "strconv"

// GetBlockCount returns the current block height

func (d *Coind) GetBlockCount() (count uint64, err error) {
	r, err := d.client.call("getblockcount", nil)
	if err = handleError(err, &r); err != nil {
		return
	}
	count, err = strconv.ParseUint(string(r.Result), 10, 64)
	return
}
