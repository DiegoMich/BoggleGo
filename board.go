package main

import (
	"strings"
)

type Cell struct {
	neighbors []*Cell
	value     string
	visited   bool
}

type Board []*Cell

func NewBoard(dice []string) *Board {
	var c1, c2, c3, c4, c5, c6, c7, c8, c9,
		c10, c11, c12, c13, c14, c15, c16 Cell
	c1.value = dice[0]
	c2.value = dice[1]
	c3.value = dice[2]
	c4.value = dice[3]
	c5.value = dice[4]
	c6.value = dice[5]
	c7.value = dice[6]
	c8.value = dice[7]
	c9.value = dice[8]
	c10.value = dice[9]
	c11.value = dice[10]
	c12.value = dice[11]
	c13.value = dice[12]
	c14.value = dice[13]
	c15.value = dice[14]
	c16.value = dice[15]

	c1.neighbors = []*Cell{&c2, &c5, &c6}
	c2.neighbors = []*Cell{&c1, &c3, &c5, &c6, &c7}
	c3.neighbors = []*Cell{&c2, &c4, &c6, &c7, &c8}
	c4.neighbors = []*Cell{&c3, &c7, &c8}
	c5.neighbors = []*Cell{&c1, &c2, &c6, &c9, &c10}
	c6.neighbors = []*Cell{&c1, &c2, &c3, &c5, &c7, &c9, &c10, &c11}
	c7.neighbors = []*Cell{&c2, &c3, &c4, &c6, &c8, &c10, &c11, &c12}
	c8.neighbors = []*Cell{&c3, &c4, &c7, &c11, &c12}
	c9.neighbors = []*Cell{&c5, &c6, &c10, &c13, &c14}
	c10.neighbors = []*Cell{&c5, &c6, &c7, &c9, &c11, &c13, &c14, &c15}
	c11.neighbors = []*Cell{&c6, &c7, &c8, &c10, &c12, &c14, &c15, &c16}
	c12.neighbors = []*Cell{&c7, &c8, &c11, &c15, &c16}
	c13.neighbors = []*Cell{&c9, &c10, &c14}
	c14.neighbors = []*Cell{&c9, &c10, &c11, &c13, &c15}
	c15.neighbors = []*Cell{&c10, &c11, &c12, &c14, &c16}
	c16.neighbors = []*Cell{&c11, &c12, &c15}

	b := Board{&c1, &c2, &c3, &c4, &c5, &c6, &c7, &c8,
		&c9, &c10, &c11, &c12, &c13, &c14, &c15, &c16}
	return &b
}

func WordExists(runes []rune, cells []*Cell) bool {
	//Exit condition: found
	if len(runes) == 0 {
		return true
	}

	// Get first letter from word (or sub-word)
	l := string(runes[0])
	// Iterate board
	for _, c := range cells {
		if c.visited {
			continue
		}
		if strings.ToLower(c.value) == l {
			c.visited = true
			if WordExists(runes[1:], c.neighbors) {
				return true
			}
			c.visited = false
		}
	}
	return false
}

func (board Board) Reset() {
	for _, c := range board {
		c.visited = false
	}
}
