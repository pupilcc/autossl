package util

import (
	"github.com/btcsuite/btcutil/base58"
	"github.com/google/uuid"
	"github.com/lithammer/shortuuid"
)

type base58Encoder struct{}

func (enc base58Encoder) Encode(u uuid.UUID) string {
	return base58.Encode(u[:])
}

func (enc base58Encoder) Decode(s string) (uuid.UUID, error) {
	return uuid.FromBytes(base58.Decode(s))
}

func GenerateID() string {
	enc := base58Encoder{}
	return shortuuid.NewWithEncoder(enc)
}
