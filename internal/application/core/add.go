package core

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/UtkarshM-hub/Lit/internal/application/core/util"
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
	fmt.Println(outputFilePath + "/" + filename)
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

func CoreAdd(data []FileInfo) {

	// Get the directory path
	dir, err := util.FindDirectory(".lit")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	objectFilePath := filepath.Join(dir, "./.lit/objects")

	// Loop over the data and store it in compressed format
	for i, v := range data {
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

		data[i].SHA1 = sha1Hash
		// Store the data in index file in following form
		// Name Last-Modified SHA-1-Hash Permissions
	}
	writeToIndex(data, dir)
}

func writeToIndex(Data []FileInfo, path string) {
	// get index file path
	indexFilePath := filepath.Join(path, "./.lit/index")

	// sort the data first
	sort.Slice(Data,func(i,j int)bool{
		return Data[i].FileName>Data[j].FileName
	})

	// fmt.Println(Data)

	var rawData []string
	for _, v := range Data {
		line := fmt.Sprintf("%v %v %v %v %v %v %v", v.FileName, v.FileModifiedAt, v.FileSize, v.FilePerm, v.SHA1, strings.Replace(v.FilePath," ","||",-1),v.FileStatus)

		// fmt.Printf("Name: %v\nPath: %v\nTime: %v\nPerm: %v\nSize: %v\nSHA1: %v\n\n", v.FileName, strings.Replace(v.FilePath," ","||",-1), v.FileModifiedAt, v.FilePerm, v.FileSize,v.SHA1)
		rawData = append(rawData, line)
	}

	err := os.WriteFile(indexFilePath, []byte(strings.Join(rawData, "\n")), 0644)
	if err != nil {
		fmt.Println(err.Error())
	}
}
