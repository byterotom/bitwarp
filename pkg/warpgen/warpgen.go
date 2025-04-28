package warpgen

import (
	"bufio"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)

type Warp struct {
	FileName    string   `json:"file_name"`
	FileSize    int      `json:"file_size"`
	FileHash    string   `json:"file_hash"`
	TotalChunks int      `json:"total_chunks"`
	ChunkSize   int      `json:"chunk_size"`
	Chunk       []string `json:"chunk"`
}

const CHUNK_SIZE = 1 << 20

// function to create warp file
func CreateWarpFile(filePath string) {

	// getting file meta data
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		log.Fatalf("error reading file stats: %v", err)
	}

	// open file pointer
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer file.Close()

	// create warp
	warp := &Warp{
		FileName:  fileInfo.Name(),
		FileSize:  int(fileInfo.Size()),
		ChunkSize: CHUNK_SIZE,
	}

	// stream reading from file
	buffer := make([]byte, CHUNK_SIZE)
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

	// load file in memory
	data, err := os.ReadFile(warpFilePath)
	if err != nil {
		log.Fatalf("error reading warp file%v", err)
	}

	// unmarshal json to warp
	var warp Warp
	if err := json.Unmarshal(data, &warp); err != nil {
		log.Fatalf("error unmarshalling warp file: %v", err)
	}
	return &warp
}

// function to merge downloaded chunks
func (w *Warp) MergeChunks() {
	log.Printf("Merging....")

	chunkDir := "storage/temp/" + w.FileHash + "/"
	fileDir := "storage/downloads/"

	filePath := fileDir + w.FileName

	err := os.MkdirAll(fileDir, os.ModePerm)
	if err != nil {
		log.Fatalf("error creating file directory: %v", err)
	}

	outFile, err := os.Create(filePath)
	if err != nil {
		log.Fatalf("error opening output file: %v", err)
	}
	defer outFile.Close()

	writer := bufio.NewWriter(outFile)
	defer writer.Flush()

	// copy all chunks to a output file
	for i := range w.TotalChunks {

		// log.Printf("merging %d", i)

		chunkPath := chunkDir + fmt.Sprint(i)

		inFile, err := os.Open(chunkPath)
		if err != nil {
			os.Remove(filePath)
			log.Fatalf("error opening chunk file: %v", err)
		}

		reader := bufio.NewReader(inFile)

		_, err = io.Copy(writer, reader)

		if err != nil {
			os.Remove(filePath)
			log.Fatalf("error merging: %v", err)
		}

	}

	// delete all chunks after successfull merge
	err = os.RemoveAll(chunkDir)
	if err != nil {
		log.Fatalf("error removing chunks: %v", err)
	}

	log.Printf("file downloaded successfully at %s", filePath)
}

// function to read chunk to send
func (w *Warp) ReadChunk(chunkNo int, isSeeder bool) ([]byte, error) {

	// read entire downloaded chunk from /temp if not a seeder
	if !isSeeder {
		filePath := fmt.Sprintf("storage/temp/%s/%d", w.FileHash, chunkNo)
		data, err := os.ReadFile(filePath)
		if err != nil {
			log.Printf("error reading chunk %v", err)
			return []byte{}, err
		}
		return data, nil
	}

	// read chunk directly from file if seeder
	buf := make([]byte, w.ChunkSize)
	filePath := "storage/downloads/" + w.FileName
	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("error opening file: %v", err)
		return []byte{}, err
	}
	defer file.Close()

	offset := int64(chunkNo * w.ChunkSize)
	section := io.NewSectionReader(file, offset, int64(w.ChunkSize))

	n, err := section.Read(buf)
	if err != nil && err != io.EOF {
		log.Printf("error reading chunk: %v", err)
		return []byte{}, err
	}

	return buf[:n], nil
}

// function to create(write) chunk in temp
func CreateChunk(fileHash string, chunkNo int64, data []byte) {

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

// function to verify chunk hash
func Verify(computedHash string, data []byte) bool {
	return hash(data) == computedHash
}
