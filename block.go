package main

import (
	"bytes"
	"encoding/gob"
	"log"
	"time"
)

type Block struct {
	Hash          []byte
	PrevBlockHash []byte
	Data          []byte
	Timestamp     int64
	Nonce		  int
}

//func (b *Block) SetHash() {
//	bTimeStamp := []byte(strconv.FormatInt(b.Timestamp, 10))
//	blockAsBytes := bytes.Join([][]byte{b.PrevBlockHash, b.Data, bTimeStamp}, []byte{})
//	hash := sha256.Sum256(blockAsBytes)
//	b.Hash = hash[:]
//}

func NewBlock(data string, prevBlockHash []byte) *Block {
	block := &Block{[]byte{}, prevBlockHash, []byte(data), time.Now().Unix(), 0}
	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()
	block.Hash = hash[:]
	block.Nonce = nonce
	return block
}

func NewGenesisBlock(starting string) *Block {
	return NewBlock(starting, []byte{})
}

func (b *Block) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(b)
	if err != nil {
		log.Panic(err)
	}
	return result.Bytes()
}

func DeserializeBlock(d []byte) *Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}
	return &block
}


