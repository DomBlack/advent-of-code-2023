package trie

import (
	"github.com/cockroachdb/errors"
)

// Trie is a tree data structure that stores strings against values.
//
// The trie (prefix tree) allows for character by character lookup of a string.
type Trie[V any] struct {
	root *Node[V]
}

// Node is a node in the Trie.
type Node[V any] struct {
	hasValue bool
	value    V
	children map[rune]*Node[V]
}

type RuneFinder[V any] interface {
	// Next returns the next node in the tree for the given rune
	// or nil if no such node exists.
	Next(r rune) RuneFinder[V]

	// Value returns the value stored at this node.
	Value() (value V, found bool)
}

// New creates a new Trie.
func New[V any]() *Trie[V] {
	return &Trie[V]{
		root: &Node[V]{
			children: make(map[rune]*Node[V]),
		},
	}
}

// Insert inserts a new key/value pair into the Trie.
//
// It returns an error if the key already exists.
func (t *Trie[V]) Insert(key string, value V) error {
	n := t.root
	for _, r := range key {
		if n.children[r] == nil {
			n.children[r] = &Node[V]{
				children: make(map[rune]*Node[V]),
			}
		}
		n = n.children[r]
	}

	if n.hasValue {
		return errors.Newf("key already exists: %q", key)
	}

	n.hasValue = true
	n.value = value

	return nil
}

// MustInsert inserts a new key/value pair into the Trie, however
// unlike [Insert], it panics if the key already exists.
func (t *Trie[V]) MustInsert(key string, value V) *Trie[V] {
	if err := t.Insert(key, value); err != nil {
		panic(err)
	}

	return t
}

// Find returns the value for the given key, or false if the key does not exist.
func (t *Trie[V]) Find(key string) (value V, found bool) {
	n := t.root
	for _, r := range key {
		if n.children[r] == nil {
			return value, false
		}
		n = n.children[r]
	}

	return n.value, n.hasValue
}

// SubstrMatches returns all the matches within the given text in the
// order that they where found.
//
// If there are multiple overlapping matches, all will be returned.
func (t *Trie[V]) SubstrMatches(text string) (values []V) {
	possibleMatches := make([]RuneFinder[V], 0)
	nextSet := make([]RuneFinder[V], 0)

	for _, r := range text {
		// Find all the next possible matches
		nextSet = nextSet[:0] // reset
		for _, pos := range possibleMatches {
			if next := pos.Next(r); next != nil {
				nextSet = append(nextSet, next)
			}
		}
		if next := t.root.Next(r); next != nil {
			nextSet = append(nextSet, next)
		}

		// Swap the slices
		possibleMatches, nextSet = nextSet, possibleMatches

		// Check if any of the possible matches are valid
		for _, pos := range possibleMatches {
			if value, found := pos.Value(); found {
				values = append(values, value)
			}
		}
	}
	return values
}

func (t *Trie[V]) Next(r rune) RuneFinder[V] {
	return t.root.Next(r)
}

func (t *Trie[V]) Value() (V, bool) {
	return t.root.Value()
}

func (n *Node[V]) Next(r rune) RuneFinder[V] {
	if n.children == nil {
		return nil
	}

	child, found := n.children[r]
	if !found {
		return nil
	}
	return child
}

func (n *Node[V]) Value() (V, bool) {
	return n.value, n.hasValue
}
