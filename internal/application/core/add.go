package core

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func compressFile(filename, inputFilePath, outputFilePath string) error {

	err := os.MkdirAll(outputFilePath, os.ModePerm)
	if err != nil {
		return err
	}

	// Read the content of the input file
	inputData, err := os.ReadFile(inputFilePath)
	if err != nil {
		return err
	}

	// Create a buffer to store the compressed data
	var compressedBuffer bytes.Buffer

	// Create a zlib writer with default compression level
	writer := zlib.NewWriter(&compressedBuffer)

	// Write the uncompressed data to the zlib writer
	_, err = writer.Write(inputData)
	if err != nil {
		return err
	}

	// Close the writer to flush any remaining data
	writer.Close()
	// Write the compressed data to the output file
	err = os.WriteFile(outputFilePath+"/"+filename, compressedBuffer.Bytes(), 0644)
	if err != nil {
		return err
	}

	return nil
}

func calculateSHA1(input string) string {
	inputBytes := []byte(input)
	hasher := sha1.New()
	hasher.Write(inputBytes)
	hashSum := hasher.Sum(nil)
	hashString := hex.EncodeToString(hashSum)

	return hashString
}

func createNewObject(status string, newmp *map[string]FileInfo, files []FileInfo, objectFilePath string) {
	for _, v := range files {
		// Get File Content
		content, err := os.ReadFile(v.FilePath)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		// Construct the header with object type+filesize+\0
		header := "blob" + " " + strconv.Itoa(int(v.FileSize)) + "\\0"

		// Find sha1 hash of the header+content
		sha1Hash := calculateSHA1(header + string(content))

		outputFilePath := objectFilePath + "/" + sha1Hash[:2]

		err = compressFile(sha1Hash[2:], v.FilePath, outputFilePath)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		v.SHA1 = sha1Hash
		v.FileStatus = status
		(*newmp)[v.FilePath] = v
	}
}

func CoreAdd(path string, untracked, modified, deleted []FileInfo) {

	objectFilePath := filepath.Join(path, "./.lit/objects")
	indexFilePath := filepath.Join(path, "./.lit/index")

	mp := GetIndexFileContent(indexFilePath)

	// change the status as modified in the index file
	createNewObject("M", &mp, modified, objectFilePath)

	createNewObject("N", &mp, untracked, objectFilePath)

	// change the status as deleted in index file
	for _, v := range deleted {
		currentStruct := mp[v.FilePath]
		currentStruct.FileStatus = "D"
		mp[v.FilePath] = currentStruct
	}

	// Loop over the untracked data and store it in compressed format

	writeToIndex(mp, indexFilePath)

}

func writeToIndex(mp map[string]FileInfo, path string) {

	// fmt.Println(Data)

	var rawData []string
	for _, v := range mp {
		fmt.Println(v.FilePath)
		line := fmt.Sprintf("%v %v %v %v %v %v %v %v", v.FileName, v.FileModifiedAt, v.FileSize, v.FilePerm, v.SHA1, strings.Replace(v.FilePath, " ", "||", -1), v.FileStatus, v.CommitStatus)

		// fmt.Printf("Name: %v\nPath: %v\nTime: %v\nPerm: %v\nSize: %v\nSHA1: %v\n\n", v.FileName, strings.Replace(v.FilePath," ","||",-1), v.FileModifiedAt, v.FilePerm, v.FileSize,v.SHA1)
		rawData = append(rawData, line)
	}

	err := os.WriteFile(path, []byte(strings.Join(rawData, "\n")), 0644)
	if err != nil {
		fmt.Println(err.Error())
	}
}