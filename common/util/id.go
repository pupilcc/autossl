package util

import gonanoid "github.com/matoous/go-nanoid/v2"

const ALPHABET = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenerateID() string {
	if id, err := gonanoid.Generate(ALPHABET, 12); err == nil {
		return id
	}
	return ""
}
