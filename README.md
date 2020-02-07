# Blockchain Implementation using Merkle Patricia Trie (MPT)

This project is continued development using Go on top of the MPT implementation https://github.com/pontakornp/cs686_blockchain_P1_Go.

The main classes include Block and Blockchain.

Block class consist of the following functions:
- Block Initialization
- Encode block to JSON
- Decode JSON to block
- Hash Block

Blockchain class consist of the following functions:
- Blockchain Initialization
- Insert - insert block into blockchain
- Get - get list of blocks from the height
- Contains Block - check if block list contains a given block
- Encode blockchain to JSON
- Decode JSON to blockchain
