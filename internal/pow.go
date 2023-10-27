package internal

import (
	"crypto/sha256"
	"encoding/hex"
	"math"
	"math/big"
	"strconv"
)

func Run(data string, difficulty int) (string, int) {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-difficulty))

	var hash string
	var nonce int
	for nonce < math.MaxInt64 {
		hash := calculateHash(data, nonce)
		hashInt := new(big.Int)
		hashInt.SetString(hash, 16)

		if hashInt.Cmp(target) == -1 {
			break
		} else {
			nonce++
		}
	}

	return hash, nonce
}

func calculateHash(data string, nonce int) string {
	record := data + strconv.Itoa(nonce)

	hash := sha256.Sum256([]byte(record))
	return hex.EncodeToString(hash[:])
}
