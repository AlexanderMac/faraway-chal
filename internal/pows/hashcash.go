package pows

import (
	"crypto/sha256"
	"encoding/hex"
	"math"
	"math/big"
	"strconv"
)

type Hashcash struct{}

func (hashcash *Hashcash) Solve(challenge string, difficulty int) int {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-difficulty))

	var hash string
	var nonce int
	for nonce < math.MaxInt64 {
		hash = hashcash.calculateHash(challenge, nonce)
		hashInt := new(big.Int)
		hashInt.SetString(hash, 16)

		if hashInt.Cmp(target) == -1 {
			break
		} else {
			nonce++
		}
	}

	return nonce
}

func (hashcash *Hashcash) Validate(challenge string, solutionNonce int, difficulty int) (bool, error) {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-difficulty))

	hashStr := hashcash.calculateHash(challenge, solutionNonce)
	hash, err := hex.DecodeString(hashStr)
	if err != nil {
		return false, err
	}
	var hashInt big.Int
	hashInt.SetBytes(hash[:])

	return hashInt.Cmp(target) == -1, nil
}

func (hashcash *Hashcash) calculateHash(data string, nonce int) string {
	record := data + strconv.Itoa(nonce)

	hash := sha256.Sum256([]byte(record))
	return hex.EncodeToString(hash[:])
}
