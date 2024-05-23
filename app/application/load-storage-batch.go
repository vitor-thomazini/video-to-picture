package application

import (
	"bytes"
	"encoding/gob"
	"log"
	"os"

	"github.com/vitor-thomazini/video-to-picture/app/domain"
)

type LoadStorageBatch struct {
}

func NewLoadStorageBatch() LoadStorageBatch {
	return LoadStorageBatch{}
}

func (b LoadStorageBatch) Execute() domain.StorageBatch {
	by, err := os.ReadFile("data/xpto")
	if err != nil {
		return domain.NewStorageBatch()
	}

	dec := gob.NewDecoder(bytes.NewBuffer(by))
	var q domain.StorageBatch
	err = dec.Decode(&q)
	if err != nil {
		log.Fatal("decode error:", err)
	}
	return q
}
