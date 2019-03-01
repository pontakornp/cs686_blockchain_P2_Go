package p2

import (
	"../p1"
)

type Header struct {
	height int32
	timestamp int64
	hash string
	parent_hash string
	size int32
}

type Block struct {
	header Header
	value  p1.MerklePatriciaTrie
}

func Initial(height int32, parent_hash string, value p1.MerklePatriciaTrie) (string, error) {
	return "", nil
}

func DecodeFromJSON(jsonString string) (Block, error) {
	return Block{}, nil
}

func EncodeToJSON(block Block) (string, error) {
	return "", nil
}


