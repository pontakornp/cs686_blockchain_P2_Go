package p2

type Blockchain struct {
	chain map[int32][]Block
	length int32 // highest block height
}

// get list of blocks from the height
func (*Blockchain) Get(height int32) ([]Block, error) {
	return []Block{}, nil
}

func (*Blockchain) Insert(block Block) {
}

func (*Blockchain) EncodeToJSON() (string, error) {
	return "", nil
}

func (*Blockchain) DecodeFromJSON(json string) (Block, error) {
	return Block{}, nil
}

