package main

import (
	"testing"
)

func TestTrieLoadAndSearch(t *testing.T) {
	words := NewTrie()
	LoadWordsListFromFile(words)
	found := words.Search("zapato")
	if !found {
		t.Fatalf("expected word not found")
	}
}
