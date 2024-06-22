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

// compress and write the file to /objects directory
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

// Create object file for the files and change the status
// It takes status to be applied, index file (present in map), files on which the changes are going to be applied
// object file path (where the object is going to be stored)
func createNewObject(status string, newmp *map[string]FileInfo, files []FileInfo, objectFilePath string) {
	// iterate over the files which we want to add in the staging area
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

		// compress and write it to the /objects directory with first 2 characters as directory name and remaining as file name
		err = compressFile(sha1Hash[2:], v.FilePath, outputFilePath)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		// change the information accordingly
		v.SHA1 = sha1Hash
		v.FileStatus = status

		// modify the changes in the index file map
		(*newmp)[v.FilePath] = v
	}
}

// takes current file information from status command and writes it to the index file 
// resulting into moving those files into staging area
// takes root directory filepath and fileinfo
func CoreAdd(path string, untracked, modified, deleted []FileInfo) {

	objectFilePath := filepath.Join(path, "./.bit/objects")
	indexFilePath := filepath.Join(path, "./.bit/index")

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

	// write the changed data to index which will be new state of index file
	writeToIndex(mp, indexFilePath)

}

// write the changed index file content to index file
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