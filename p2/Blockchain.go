package p2

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
)

type Blockchain struct {
	Chain map[int32][]Block
	Length int32 // highest block height
}

// initialize blockchain
func (bc *Blockchain) Initial() {
	bc.Chain = make(map[int32][]Block, 0)
	bc.Length = 0
}

// check if blockchain is empty
func (bc *Blockchain) isEmpty() bool {
	return reflect.DeepEqual(bc, nil)
}

// get list of blocks from the height
func (bc *Blockchain) Get(height int32) ([]Block, error) {
	if value, ok := bc.Chain[height]; ok {
		return value, nil
	}
	return []Block{}, errors.New("no block found")
}

// check is block list contains a given block
func containsBlock(blockList []Block, block Block) bool {
	for _, b := range blockList {
		if b == block {
			return true
		}
	}
	return false
}

// insert block into blockchain
func (bc *Blockchain) Insert(block Block) error {
	if block.isEmpty() {
		return errors.New("block is empty")
	}
	height := block.Header.Height
	blockList := bc.Chain[height]
	// if block hash already already
	if containsBlock(blockList, block) {
		return errors.New("contains block")
	}
	// append block to the block chain list at the height
	bc.Chain[height] = append(bc.Chain[height], block)
	// update highest block height
	if height > bc.Length {
		bc.Length = height
	}
	return nil
}

// encode blockchain to json string
func (bc *Blockchain) EncodeToJson() (string, error) {
	if bc.isEmpty() {
		return "", errors.New("empty blockchain")
	}
	// loop the map to get the block list
	// loop the block list to get each block
	blockJsonList := make([]BlockJson, 0)
	for _, blockList := range bc.Chain {
		for _, block := range blockList {
			blockJson := block.ConvertBlocktoBlockJson()
			blockJsonList = append(blockJsonList, blockJson)
		}
	}
	jsonStr, err := json.Marshal(blockJsonList)
	if err != nil {
		fmt.Println(err)
	}
	return string(jsonStr), nil
}

// decode json string to blockchain
func DecodeJsonToBlockChain(jsonStr string) (Blockchain, error) {
	// decodes the JSON string back to a list of block JSON strings
	// decodes each block JSON string back to a block instance
	//bList := []BlockJson{}
	bList := make([]BlockJson, 0)
	json.Unmarshal([]byte(jsonStr), &bList)
	// check if block list is empty or not
	if len(bList) == 0 {
		return Blockchain{}, errors.New("fail to decode")
	}
	// inserts every block into the blockchain
	blockchain := Blockchain{}
	blockchain.Initial()
	for _, blockJson := range bList {
		b, _ := DecodeJsonHelper(blockJson)
		blockchain.Insert(b)
	}
	return blockchain, nil
}

// print string
func (bc *Blockchain) String() string {
	content := ""
	for _, blockList := range bc.Chain {
		for _, b := range blockList {
			content += b.String()
		}
	}
	return content
}