package util

import gonanoid "github.com/matoous/go-nanoid/v2"

const ALPHABET = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenerateID() string {
	return gonanoid.MustGenerate(ALPHABET, 12)
}
