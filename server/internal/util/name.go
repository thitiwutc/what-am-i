package util

import (
	"fmt"
	"math/rand/v2"
)

var adjectives = []string{
	"Swift", "Silent", "Iron", "Lucky", "Crimson",
	"Mighty", "Rapid", "Golden", "Brave", "Shadow",
	"Silver", "Scarlet", "Azure", "Emerald", "Frost",
	"Storm", "Thunder", "Blazing", "Frozen", "Burning",
	"Wild", "Fearless", "Nimble", "Clever", "Wise",
	"Bold", "Fierce", "Gentle", "Calm", "Steady",
	"Flying", "Soaring", "Glowing", "Radiant", "Shining",
	"Hidden", "Secret", "Mystic", "Ancient", "Royal",
	"Noble", "Epic", "Cosmic", "Lunar", "Solar",
	"Electric", "Crystal", "Stone", "Steel", "Velvet",
	"Arctic", "Desert", "Ocean", "Forest", "Mountain",
	"Bright", "Dark", "Quick", "Sly", "Daring",
	"Heroic", "Valiant", "Dynamic", "Supreme", "Majestic",
	"Fearless", "Agile", "Vivid", "Glacial", "Whispering",
	"Howling", "Dancing", "Roaring", "Smiling", "Cheerful",
	"Playful", "Curious", "Cunning", "Gallant", "Stellar",
	"Galactic", "Phantom", "Thunderous", "Sparkling", "Verdant",
}

var nouns = []string{
	"Tiger", "Falcon", "Wolf", "Panda", "Fox",
	"Eagle", "Shark", "Owl", "Rhino", "Bear",
	"Lion", "Leopard", "Jaguar", "Panther", "Cobra",
	"Viper", "Dragon", "Phoenix", "Griffin", "Kraken",
	"Bison", "Buffalo", "Moose", "Stag", "Boar",
	"Hawk", "Raven", "Crow", "Sparrow", "Robin",
	"Dolphin", "Whale", "Octopus", "Seal", "Otter",
	"Lynx", "Cheetah", "Badger", "Wolverine", "Cougar",
	"Ram", "Goat", "Yak", "Camel", "Horse",
	"Penguin", "Koala", "Kangaroo", "Platypus", "Beaver",
	"Squirrel", "Rabbit", "Hedgehog", "Raccoon", "Skunk",
	"Tortoise", "Lizard", "Gecko", "Frog", "Salamander",
	"Butterfly", "Dragonfly", "Firefly", "Scorpion", "Beetle",
	"Comet", "Meteor", "Nova", "Nebula", "Galaxy",
	"Mountain", "River", "Volcano", "Hurricane", "Tornado",
	"Blizzard", "Thunder", "Lightning", "Tempest", "Cyclone",
}

func GeneratePlayerName() string {
	adjective := adjectives[rand.IntN(len(adjectives))]
	noun := nouns[rand.IntN(len(nouns))]
	number := rand.IntN(1000)

	return fmt.Sprintf("%s%s%03d", adjective, noun, number)
}
