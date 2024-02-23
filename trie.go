package main

import "unicode"

type Trie struct {
	root *TrieNode
}

type TrieNode struct {
	children map[rune]*TrieNode
	isEnd    bool
}

func NewTrie() *Trie {
	return &Trie{
		root: &TrieNode{
			children: make(map[rune]*TrieNode),
		},
	}
}

func (t *Trie) Insert(word string) {
	node := t.root
	for _, v := range word {
		c := unicode.ToLower(v)
		if _, ok := node.children[c]; !ok {
			node.children[c] = &TrieNode{children: make(map[rune]*TrieNode)}
		}
		node = node.children[c]
	}
	node.isEnd = true
}

func (t *Trie) Search(word string) bool {
	node := t.root
	for _, v := range word {
		c := unicode.ToLower(v)
		if _, ok := node.children[c]; !ok {
			return false
		}
		node = node.children[c]
	}
	return node.isEnd
}
