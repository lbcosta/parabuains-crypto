package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

var bc *Blockchain

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
func (bc *Blockchain) AddBlock(newBlock *Block) {
	newBlockChain := append(bc.blocks, newBlock)
	bc.ReplaceChain(newBlockChain)
}

func (bc *Blockchain) GetLatestBlock() *Block {
	return bc.blocks[len(bc.blocks)-1]
}

func NewGenesisBlock() *Block {
	return NewBlock("Genesis Block", &Block{})
}

func NewBlockchain() *Blockchain {
	return &Blockchain{[]*Block{NewGenesisBlock()}}
}

func (bc *Blockchain) ReplaceChain(candidateChain []*Block) {
	if len(candidateChain) > len(bc.blocks) {
		bc.blocks = candidateChain
	}
}

func IsBlockValid(block, prevBlock Block) bool {
	if block.Index != prevBlock.Index+1 {
		return false
	}

	if bytes.Compare(block.PrevHash, prevBlock.Hash) != 0 {
		return false
	}

	if bytes.Compare(block.Hash, GenerateHash(&block)) != 0 {
		return false
	}

	return true
}

func HandleGetBlockchain(c *gin.Context) {
	c.JSON(200, bc.blocks)
}

func HandleAddBlock(c *gin.Context) {
	newTransaction := c.Query("transaction")
	if newTransaction == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "a transaction is needed to create a block",
		})
		return
	}

	newBlock := NewBlock(newTransaction, bc.blocks[len(bc.blocks)-1])

	if !IsBlockValid(*newBlock, *bc.GetLatestBlock()) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "block is not valid",
		})
		return
	}

	bc.AddBlock(newBlock)
	c.JSON(http.StatusCreated, bc.GetLatestBlock())
}

func runHttpServer() error {
	r := gin.Default()
	if err := r.SetTrustedProxies(nil); err != nil {
		log.Fatal(err)
	}

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Parabuains!",
		})
	})
	r.GET("/blockchain", HandleGetBlockchain)
	r.POST("/block", HandleAddBlock)

	addr := "localhost:" + os.Getenv("PORT")
	if err := r.Run(addr); err != nil {
		return err
	}

	return nil
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println()
	fmt.Println("$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$")
	fmt.Println("$$$$$$$$$$ Parabuains (PRBS)! $$$$$$$$$$")
	fmt.Println("$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$")
	fmt.Println()

	bc = NewBlockchain()

	log.Fatal(runHttpServer())
}
