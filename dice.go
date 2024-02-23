package main

import "math/rand"

var dice = [][]string{
	{"N", "D", "S", "E", "A", "O"},
	{"A", "O", "U", "E", "A", "I"},
	{"N", "I", "T", "A", "G", "U"},
	{"V", "O", "N", "J", "S", "L"},
	{"E", "S", "O", "Ñ", "A", "D"},
	{"E", "Q", "O", "S", "H", "D"},
	{"C", "E", "N", "O", "L", "S"},
	{"D", "T", "A", "R", "O", "I"},
	{"C", "N", "I", "R", "T", "F"},
	{"P", "S", "C", "E", "L", "O"},
	{"E", "H", "I", "X", "U", "R"},
	{"B", "O", "M", "L", "E", "Z"},
	{"A", "R", "E", "C", "M", "A"},
	{"S", "A", "C", "E", "N", "O"},
	{"P", "O", "D", "E", "T", "A"},
	{"B", "R", "A", "E", "L", "A"},
}

func ThrowDice(seed int64) []string {
	src := rand.NewSource(seed)
	rnd := rand.New(src)

	var result []string

	// Randomize dice faces
	for _, d := range dice {
		result = append(result, d[rnd.Intn(len(d)-1)])
	}

	// Randomize positions in board
	for i := 0; i < len(result); i++ {
		newPos := rnd.Intn(len(result))
		result[i], result[newPos] = result[newPos], result[i]
	}
	return result
}
