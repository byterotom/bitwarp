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
	outputDir := "storage/warp/"
	os.MkdirAll(outputDir, os.ModePerm)
	warpFilePath := outputDir + warp.FileName + ".json"
	
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

	// get absolute file path
	absFilePath, err := filepath.Abs(warpFilePath)
	if err != nil {
		log.Fatalf("error getting absolute path: %v", err)
	}

	// open warpfile in read mode
	data, err := os.ReadFile(absFilePath)
	if err != nil {
		log.Fatalf("error reading warp file%v", err)
	}

	var warp *Warp

	json.Unmarshal(data, &warp)

	return warp
}

func (w *Warp) MergeChunks() {

	chunkDir := "storage/temp/" + w.FileHash + "/"
	fileDir := "storage/downloads/"

	filePath := fileDir + w.FileName

	err := os.MkdirAll(fileDir, os.ModePerm)
	if err != nil {
		log.Fatalf("error creating file directory: %v", err)
	}

	outFile, err := os.Create(filePath)
	if err != nil {
		log.Fatalf("error opening output file %v", err)
	}
	defer outFile.Close()

	for i := range w.TotalChunks {

		chunkPath := chunkDir + fmt.Sprint(i)

		inFile, err := os.Open(chunkPath)
		if err != nil {
			os.Remove(filePath)
			log.Fatalf("error merging: %v", err)
		}

		_, err = io.Copy(outFile, inFile)
		inFile.Close()
		if err != nil {
			os.Remove(filePath)
			log.Fatalf("error merging: %v", err)
		}
	}

	err = os.RemoveAll(chunkDir)
	if err != nil {
		log.Fatalf("error removing chunks: %v", err)
	}
}

func (w *Warp) ReadChunk(filePath string, chunkNo int) {

}

func CreateChunk(fileHash string, chunkNo int, data []byte) {

	chunkDir := "storage/temp/" + fileHash + "/"
	chunkPath := chunkDir + fmt.Sprint(chunkNo)

	err := os.MkdirAll(chunkDir, os.ModePerm)
	if err != nil {
		log.Printf("error creating directory %s: %v", chunkDir, err)
		return
	}

	err = os.WriteFile(chunkPath, data, 0644)
	if err != nil {
		log.Printf("error writing chunk no. %d: %v", chunkNo, err)
	}

}

// function to hash in sha256
func hash(data []byte) string {
	hash := sha256.Sum256(data)
	return fmt.Sprintf("%x", hash)
}
