package util

import (
	"fmt"
	"math/rand/v2"
)

var adjectives = []string{
	"Swift", "Silent", "Iron", "Lucky", "Crimson",
	"Mighty", "Rapid", "Golden", "Brave", "Shadow",
}

var nouns = []string{
	"Tiger", "Falcon", "Wolf", "Panda", "Fox",
	"Eagle", "Shark", "Owl", "Rhino", "Bear",
}

func GeneratePlayerName() string {
	adjective := adjectives[rand.IntN(len(adjectives))]
	noun := nouns[rand.IntN(len(nouns))]
	number := rand.IntN(1000)

	return fmt.Sprintf("%s%s%03d", adjective, noun, number)
}
