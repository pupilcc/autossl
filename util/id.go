package util

import gonanoid "github.com/matoous/go-nanoid/v2"

func GenerateID() string {
	if id, err := gonanoid.New(12); err == nil {
		return id
	}
	return ""
}
