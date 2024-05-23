package application

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"os"

	"github.com/vitor-thomazini/video-to-picture/app/domain"
)

type SaveStorageBatch struct {
}

func NewSaveStorageBatch() SaveStorageBatch {
	return SaveStorageBatch{}
}

func (b SaveStorageBatch) Execute(latestText string, latestIndex int, data map[string][]domain.Resource) {
	var network bytes.Buffer // Stand-in for a network connection
	enc := gob.NewEncoder(&network)

	d := domain.StorageBatch{
		LatestText:  latestText,
		LatestIndex: latestIndex,
		Data:        data,
	}

	err := enc.Encode(d)
	if err != nil {
		log.Fatal("encode error:", err)
	}

	f, err := os.Create("data/xpto")
	check(err)

	defer f.Close()

	n2, err := f.Write(network.Bytes())
	check(err)
	fmt.Printf("wrote %d bytes\n", n2)

	f.Sync()
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
