package warpgen

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

type Warp struct {
	FileName    string   `json:"file_name"`
	FileSize    int      `json:"file_size"`
	FileHash    string   `json:"file_hash"`
	TotalChunks int      `json:"total_chunks"`
	ChunkSize   int      `json:"chunk_size"`
	Chunk       []string `json:"chunk"`
}

const chunkSize = 1 << 20

// function to create warp file
func CreateWarpFile(filePath string) {

	// get absolute file path
	absFilePath, err := filepath.Abs(filePath)
	if err != nil {
		log.Fatalf("error getting absolute path: %v", err)
	}

	// getting file meta data
	fileInfo, err := os.Stat(absFilePath)
	if err != nil {
		log.Fatalf("error reading file stats: %v", err)
	}

	// open file pointer
	file, err := os.Open(absFilePath)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer file.Close()

	// create warp
	warp := &Warp{
		FileName:  fileInfo.Name(),
		FileSize:  int(fileInfo.Size()),
		ChunkSize: chunkSize,
	}

	// stream reading from file
	buffer := make([]byte, chunkSize)
	for {
		n, err := file.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error reading file: %v", err)
		}

		// hash the chunk
		hashStr := hash(buffer[:n])
		warp.Chunk = append(warp.Chunk, hashStr)

		// keep on computing file hash
		warp.FileHash = hash([]byte(warp.FileHash + hashStr))
	}

	// set total chunks
	warp.TotalChunks = len(warp.Chunk)

	// create warp file
	outputDir := filepath.Dir(absFilePath)
	warpFilePath := filepath.Join(outputDir, warp.FileName+".json")
	warpFile, err := os.Create(warpFilePath)
	if err != nil {
		log.Fatalf("error creating warp file: %v", err)
	}
	defer warpFile.Close()

	json.NewEncoder(warpFile).Encode(warp)

	fmt.Printf("warp file created: %s\n", warpFilePath)
}

// function to read warp file
func ReadWarpFile(warpFilePath string) *Warp {
	// open warpfile in read mode
	data, err := os.ReadFile(warpFilePath)
	if err != nil {
		log.Fatalf("error reading warp file%v", err)
	}

	var warp *Warp

	json.Unmarshal(data, &warp)

	return warp
}

// function to hash in sha256
func hash(data []byte) string {
	hash := sha256.Sum256(data)
	return fmt.Sprintf("%x", hash)
}
