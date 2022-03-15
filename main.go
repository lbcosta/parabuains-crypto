package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"strconv"
	"time"
)

type Block struct {
	Index       int64
	Timestamp   int64
	PrevHash    []byte
	Transaction []byte
	Hash        []byte
}

func GenerateHash(block *Block) []byte {
	index := []byte(strconv.FormatInt(block.Index, 10))
	timestamp := []byte(strconv.FormatInt(block.Timestamp, 10))
	headers := bytes.Join([][]byte{index, timestamp, block.Transaction, block.PrevHash}, []byte{})
	hash := sha256.Sum256(headers)

	return hash[:]
}

func NewBlock(transaction string, prevBlock *Block) *Block {
	block := &Block{prevBlock.Index + 1,
		time.Now().Unix(),
		prevBlock.Hash,
		[]byte(transaction),
		[]byte{}}
	block.Hash = GenerateHash(block)
	return block
}

type Blockchain struct {
	blocks []*Block // TODO: Use map of arrays to make a hash->block search
}

// AddBlock is a Proof-Of-Work free method
func (bc *Blockchain) AddBlock(transaction string) {
	latestBlock := bc.blocks[len(bc.blocks)-1]
	newBlock := NewBlock(transaction, latestBlock)
	bc.blocks = append(bc.blocks, newBlock)
}

func NewGenesisBlock() *Block {
	return NewBlock("Genesis Block", &Block{})
}

func NewBlockchain() *Blockchain {
	return &Blockchain{[]*Block{NewGenesisBlock()}}
}

func (bc *Blockchain) replaceChain(candidateChain []*Block) {
	if len(candidateChain) > len(bc.blocks) {
		bc.blocks = candidateChain
	}
}

//func IsBlockValid(block, prevBlock Block) bool {
//	if block.Index != prevBlock.Index+1 {
//		return false
//	}
//
//	if bytes.Compare(block.PrevHash, prevBlock.Hash) != 0 {
//		return false
//	}
//
//	if bytes.Compare(block.Hash, GenerateHash(&block)) != 0 {
//		return false
//	}
//
//	return true
//}

func main() {
	fmt.Println("$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$")
	fmt.Println("$$$$$$$$$$ Parabuains (PRBS)! $$$$$$$$$$")
	fmt.Println("$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$")
	fmt.Println()

	bc := NewBlockchain()

	bc.AddBlock("Send 1 PRBS to Allan")
	bc.AddBlock("Send 2 more PRBS to Allan")

	for _, block := range bc.blocks {
		fmt.Printf("Prev. hash: %x\n", block.PrevHash)
		fmt.Printf("Transaction: %s\n", block.Transaction)
		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Println()
	}
}
