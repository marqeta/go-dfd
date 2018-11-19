package dfd

import (
	"crypto/rand"
	"math"
	"math/big"
	"strconv"
)

func genID() string {
	xid, _ := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
	return strconv.FormatInt(xid.Int64(), 10)
}

func idToID64(id string) int64 {
	xid64, _ := strconv.ParseInt(id, 10, 64)
	return xid64
}
