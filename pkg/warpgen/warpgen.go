package warpgen

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

type Warp struct {
	FileName    string   `json:"file_name"`
	FileType    string   `json:"file_type"`
	FileHash    string   `json:"file_hash"`
	FileSize    int      `json:"file_size"`
	TotalChunks int      `json:"total_chunks"`
	ChunkSize   int      `json:"chunk_size"`
	Chunk       []string `json:"chunk"`
}

func CreateWarpFile(filePath string) {

	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("error %v", err)
	}

	fileSize := len(data)
	chunkSize := 1 << 20

	warp := &Warp{
		FileName:  "sp",
		FileType:  "pdf",
		FileSize:  fileSize,
		ChunkSize: chunkSize,
	}

	st := 0
	for st < fileSize {

		en := min(fileSize, chunkSize+st)
		currChunk := data[st:en]

		hashStr := hash(currChunk)

		warp.Chunk = append(warp.Chunk, hashStr)
		st = en
	}

	warp.TotalChunks = len(warp.Chunk)

	warpHashStr := warp.FileName + warp.FileType + strings.Join(warp.Chunk, "")

	warp.FileHash = hash([]byte(warpHashStr))

	file, _ := os.Create("storage/" + warp.FileName + "." + warp.FileType + ".json")
	defer file.Close()

	json.NewEncoder(file).Encode(warp)
}

func ReadWarpFile(warpFilePath string) *Warp {
	data, err := os.ReadFile(warpFilePath)
	if err != nil {
		log.Fatalf("error %v", err)
	}

	var warp *Warp

	json.Unmarshal(data, &warp)

	return warp
}

func hash(data []byte) string {
	hash := sha256.Sum256(data)
	return fmt.Sprintf("%x", hash)
}
