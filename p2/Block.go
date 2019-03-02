package p2

import (
	"cs686_blockchain_P2_Go/p1"
	"fmt"
	"time"
)

type Header struct {
	height int32
	timestamp int64
	hash string
	parentHash string
	size int32
}

type Block struct {
	header Header
	value  *p1.MerklePatriciaTrie
}

func (b *Block) Initial(height int32, parentHash string, value *p1.MerklePatriciaTrie) {
	header := Header{
		height: height,
		timestamp: int64(time.Now().Unix()),
		hash: "",
		parentHash: parentHash,
		size: int32(len([]byte(fmt.Sprintf("%v", value)))),
	}
	b.header = header
	b.value = value
	b.header.hash = b.hashBlock()
}

func DecodeFromJSON(jsonString string) (Block, error) {
	return Block{}, nil
}

func EncodeToJSON(block Block) (string, error) {
	return "", nil
}

func (b *Block) hashBlock() string {
	hash_str := string(b.header.height) +
		string(b.header.timestamp) +
		b.header.parentHash +
		b.value.GetRoot()+
		string(b.header.size)
	return hash_str
}

func createGenesisBlock() *Block {
	mpt := new(p1.MerklePatriciaTrie)
	mpt.Initial()
	b := new(Block)
	b.Initial(0, "", mpt)
	return b
}

func TestBlock() {
	b := createGenesisBlock()
	fmt.Println(b.value.String())
}


