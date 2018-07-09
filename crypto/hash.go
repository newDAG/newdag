package crypto

import (
	"crypto/sha256"
)

func SHA256(hashBytes []byte) []byte {
	hasher := sha256.New()
	hasher.Write(hashBytes)
	hash := hasher.Sum(nil)
	return hash
}

func SimpleHashFromTwoHashes(left []byte, right []byte) []byte {
	var hasher = sha256.New()
	hasher.Write(left)
	hasher.Write(right)
	return hasher.Sum(nil)
}
