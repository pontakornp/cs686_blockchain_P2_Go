package p2

import (
	"crypto/sha256"
	"cs686_blockchain_P2_Go/p1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

type BlockHeader struct {
	Height int32
	Timestamp int64
	Hash string
	ParentHash string
	Size int32
}

type Block struct {
	Header BlockHeader
	Value  *p1.MerklePatriciaTrie
}

// struct for json string of blockchain
type BlockJson struct {
	BlockHeader
	MPT map[string]string
}

// check if block is empty
func (block *Block) isEmpty() bool {
	return reflect.DeepEqual(block, nil)
}

// initialize block
func (b *Block) Initial(Height int32, ParentHash string, Value *p1.MerklePatriciaTrie) {
	header := BlockHeader{
		Height: Height,
		Timestamp: int64(time.Now().Unix()),
		Hash: "",
		ParentHash: ParentHash,
		Size: int32(len([]byte(fmt.Sprintf("%v", Value)))),
	}
	b.Header = header
	b.Value = Value
	b.Header.Hash = b.hashBlock()
}

// helper function to decode json to block
func DecodeJsonHelper(blockJson BlockJson) (Block, error) {
	blockHeader := BlockHeader {
		blockJson.Height,
		blockJson.Timestamp,
		blockJson.Hash,
		blockJson.ParentHash,
		blockJson.Size,
	}
	m := blockJson.MPT // map of mpt key value
	mpt := new(p1.MerklePatriciaTrie)
	mpt.Initial()
	for k, v := range m {
		mpt.Insert(k, v)
	}
	b := Block {
		blockHeader,
		mpt,
	}
	return b, nil
}

// decode json to block
func DecodeJsonToBlock(jsonStr string) (Block, error) {
	blockJson := BlockJson{}
	json.Unmarshal([]byte(jsonStr), &blockJson)
	b, err := DecodeJsonHelper(blockJson)
	return b, err
}

// helper function to convert block to block json string
func (block *Block) ConvertBlocktoBlockJson() BlockJson {
	mpt := block.Value
	pairMap := mpt.GetPairMap()
	blockJson := BlockJson {
		block.Header,
		pairMap,
	}
	return blockJson
}

// encode block to json string
func (block *Block) EncodeToJson() (string, error) {
	blockJson := block.ConvertBlocktoBlockJson()
	jsonStr, err := json.Marshal(blockJson)
	for err != nil {
		fmt.Println("error:", err)
	}
	return string(jsonStr), nil
}

// hash the block to sha256
func (b *Block) hashBlock() string {
	str := strconv.Itoa(int(b.Header.Height)) +
		strconv.Itoa(int(b.Header.Timestamp)) +
		b.Header.ParentHash +
		b.Value.GetRoot() +
		strconv.Itoa(int(b.Header.Size))
	h := sha256.New()
	h.Write([]byte(str))
	hash_str := base64.URLEncoding.EncodeToString(h.Sum(nil))
	return hash_str
}

// create genesis block or the initial block
func createGenesisBlock() *Block {
	mpt := new(p1.MerklePatriciaTrie)
	mpt.Initial()
	b := new(Block)
	b.Initial(0, "", mpt)
	return b
}

// print block to readable string
func (b *Block) String() string {
	content := fmt.Sprintf("HEIGHT=%d\n", b.Header.Height)
	content += fmt.Sprintf("TIMESTAMP=%d\n", b.Header.Timestamp)
	content += fmt.Sprintf("HASH=%s\n", b.Header.Hash)
	content += fmt.Sprintf("PARENTHASH=%s\n", (b.Header).ParentHash)
	content += fmt.Sprintf("SIZE=%d\n", b.Header.Size)
	content += b.Value.String()
	return content
}


