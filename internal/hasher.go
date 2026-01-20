package internal

import (
	"crypto/sha256"
	"fmt"
	"hash"
	"io"
)

func NewHasher(input io.Reader) Hasher {
	hasher := sha256.New()
	output := io.TeeReader(input, hasher)

	return Hasher{
		output: output,
		hasher: hasher,
	}
}

type Hasher struct {
	output io.Reader
	hasher hash.Hash
}

func (h *Hasher) Reader() io.Reader {
	return h.output
}

func (h *Hasher) After(output *Asset) {
	sum := h.Sum()
	output.Hash = fmt.Sprintf("%x", sum)
}

func (h *Hasher) Sum() []byte {
	return h.hasher.Sum(nil)
}
