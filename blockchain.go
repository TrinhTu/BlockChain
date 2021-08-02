package main

import (
	"github.com/boltdb/bolt"
	"log"
	"sync"
)

type Blockchain struct {
	tip []byte
	db *bolt.DB
}

type BlockchainIterator struct {
	currentHash []byte
	db			*bolt.DB
}

var instantiated *Blockchain
var once sync.Once
const dbFile = "blockchain.db"
const blocksBucket = "blocks"

func NewBlockchain() *Blockchain {

	// final hash value store in DB
	var tip []byte

	// open boltDB file, if file doesn't exist => not return err
	db, err := bolt.Open(dbFile, 0600, nil)

	// open transaction read/write to store block to DB
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))

		// doesn't have any block in DB => create init block (genesis)
		if b == nil {
			genesis :=  NewGenesisBlock("Beginning Block")
			b, err = tx.CreateBucketIfNotExists([]byte(blocksBucket))
			err = b.Put(genesis.Hash, genesis.Serialize())
			err = b.Put([]byte("l"), genesis.Hash)
			tip = genesis.Hash
		} else {
			tip = b.Get([]byte("l"))
		}
		return nil
	})

	// create blockchain's instance
	bc := Blockchain{tip, db}
	return &bc
}

// AddBlock store block to DB
func (bc *Blockchain) AddBlock(data string) {
	var lastHash []byte

	// open transaction readonly to get the last hash
	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash = b.Get([]byte("l"))
		return nil
	})

	newBlock := NewBlock(data, lastHash)
	err = bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		err = b.Put(newBlock.Hash, newBlock.Serialize())
		err = b.Put([]byte("l"), newBlock.Hash)
		bc.tip = newBlock.Hash
		return nil
	})
}

func (bc *Blockchain) Iterator() *BlockchainIterator {
	bci := &BlockchainIterator{bc.tip, bc.db}
	return bci
}

func (i *BlockchainIterator) Next() *Block {
	var block *Block

	err := i.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		encodedBlock := b.Get(i.currentHash)
		block = DeserializeBlock(encodedBlock)

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	i.currentHash = block.PrevBlockHash

	return block
}
