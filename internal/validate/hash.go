package validate

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"math/big"
)

const sha256Bits = sha256.Size * 8

var (
	errZeroCount   = errors.New("less zeroes")
	errHashNotReal = errors.New("hash not from this data")
)

type Parser interface {
	GetByteTime(d []byte) []byte
	GetRndPart(d []byte) []byte
	GetByteNonce(d []byte) []byte
	GetHash(d []byte) []byte
}

type HashCash struct {
	target *big.Int
	parser Parser
}

func NewHashCash(zeroCount int, p Parser) *HashCash {
	target := big.NewInt(1)
	target.Lsh(target, uint(sha256Bits-zeroCount))
	return &HashCash{
		target: target,
		parser: p,
	}
}
func (h *HashCash) Validate(ipAddress string, data []byte) (err error) {
	var cmp big.Int
	hash := h.parser.GetHash(data)
	cmp.SetBytes(hash)
	if cmp.Cmp(h.target) != -1 {
		return errZeroCount
	}

	digest := sha256.Sum256(h.prepareData(
		h.parser.GetByteTime(data),
		h.parser.GetRndPart(data),
		h.parser.GetByteNonce(data),
		ipAddress))

	if !bytes.Equal(hash, digest[:]) {
		return errHashNotReal
	}

	return nil
}

func (h *HashCash) prepareData(generateTime, rndBytes, nonce []byte, address string) []byte {
	return bytes.Join(
		[][]byte{
			[]byte(address),
			generateTime,
			rndBytes,
			nonce,
		},
		[]byte{},
	)
}
