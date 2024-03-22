package core

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/UtkarshM-hub/Lit/internal/application/core/util"
)

type TreeInfo struct {
	Type        string
	Perm        int64
	SHA1        string
	FileName    string
	Modified_at string
	FilePath    string
	FileSize    int
}

func Commit(commitMessage string) error {
	// Get index file path
	dir, err := util.FindDirectory(".lit")
	if err != nil {
		return nil
	}

	current_active_branch, err := CurrentActiveBranch(dir)
	if err != nil {
		fmt.Println(err)
		return err
	}

	indexFilePath := filepath.Join(dir, "./.lit/index")

	// take content and append
	logsHEAD_Append := filepath.Join(dir, "./.lit/logs/HEAD")
	logsFilePath := filepath.Join(dir, "./.lit/logs/refs/heads")

	// Replace the file content
	refsFilePath := filepath.Join(dir, "./.lit/refs/heads")

	objectFilePath := filepath.Join(dir, "./.lit/objects")
	actualPath := filepath.Join(dir, "/")

	mp := GetIndexFileContent(indexFilePath)

	MainTree, err := GetTree(&mp, dir, actualPath)
	if err != nil {
		return err
	}

	var content string
	content += fmt.Sprintf("tree %v\nauthor %v <%v> %v\ncommitter %v <%v> %v\n\n%v", MainTree.SHA1, "Utkarsh Mandape", "utmandape4@gmail.com", time.Now(), "Utkarsh Mandape", "utmandape4@gmail.com", time.Now(), commitMessage)

	// find SHA1
	header := "commit" + " " + "\\0"

	// Find sha1 hash of the header+content
	sha1Hash := calculateSHA1(header + string(content))

	outputFilePath := objectFilePath + "/" + sha1Hash[:2]

	err = compressCommitContent(sha1Hash[2:], []byte(content), outputFilePath)
	if err != nil {
		return err
	}

	// // write commit object to logs/refs/heads/branchname and refs/heads/branchname
	commit_time := time.Now().String()
	err = AppendToFiles(logsFilePath+"/"+current_active_branch, "Utkarsh Mandape", "utmandape4@gmail.com", commitMessage, sha1Hash, commit_time)

	if err != nil {
		return err
	}
	err = AppendToFiles(logsHEAD_Append, "Utkarsh Mandape", "utmandape4@gmail.com", commitMessage, sha1Hash, commit_time)

	if err != nil {
		return err
	}

	err = os.MkdirAll(refsFilePath, os.ModePerm)
	if err != nil {
		return err
	}

	err = os.WriteFile(fmt.Sprintf("%v/%v", refsFilePath+"/", current_active_branch), []byte(sha1Hash), 0644)
	if err != nil {
		return err
	}

	// store to index file

	for i, v := range mp {
		currentFile := mp[v.FilePath]
		currentFile.CommitStatus = "C"
		mp[i] = currentFile
	}

	// write to index
	writeToIndex(mp, indexFilePath)

	return nil
}

func AppendToFiles(filePath, commiter, email, msg, SHA1, time string) error {
	err := util.DoesExists(filePath)
	if err != nil {
		os.Create(filePath)
	}
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	var parentHash string

	// modify the commit message and username in-order to replace space with ||
	commiter = strings.Replace(commiter, " ", "||", -1)
	msg = strings.Replace(msg, " ", "||", -1)
	data_string := strings.Split(string(data), "\n")

	if len(data) == 0 || data_string[0] == " " {
		parentHash = "0000000000000000000000000000000000000000"
	} else {
		parentHash = strings.Split(data_string[len(data_string)-1], " ")[1]
	}
	NewData := fmt.Sprintf("%v %v %v %v %v %v", parentHash, SHA1, commiter, email, time, msg)
	data_string = append(data_string, NewData)
	err = os.WriteFile(filePath, []byte(strings.Join(data_string, "\n")), 0644)
	if err != nil {
		return err
	}
	return nil
}

func compressCommitContent(filename string, content []byte, outputFilePath string) error {
	err := os.MkdirAll(outputFilePath, os.ModePerm)
	if err != nil {
		return err
	}

	// Create a buffer to store the compressed data
	var compressedBuffer bytes.Buffer

	// Create a zlib writer with default compression level
	writer := zlib.NewWriter(&compressedBuffer)

	// Write the uncompressed data to the zlib writer
	_, err = writer.Write(content)
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

func createTreeObj(dirContent []TreeInfo, path, dirName string) (TreeInfo, error) {
	splitted := strings.Split(dirName, "/")
	directoryName := splitted[len(splitted)-1]
	objectFilePath := filepath.Join(path, "./.lit/objects")

	var fileContent []string

	for _, v := range dirContent {
		line := fmt.Sprintf("%v %v %v %v %v %v %v", v.Type, v.Modified_at, v.FileName, v.SHA1, v.Perm, v.FileSize, strings.Replace(v.FilePath, " ", "||", -1))
		fileContent = append(fileContent, line)
	}
	content := strings.Join(fileContent, "\n")
	fmt.Println(dirName)
	fmt.Println(content)
	fmt.Println("--------------------------------------------")

	// find SHA1
	header := "tree" + " " + directoryName + "\\0"

	// Find sha1 hash of the header+content
	sha1Hash := calculateSHA1(header + string(content))

	outputFilePath := objectFilePath + "/" + sha1Hash[:2]

	err := compressCommitContent(sha1Hash[2:], []byte(content), outputFilePath)
	if err != nil {
		fmt.Println(err.Error())
		return TreeInfo{}, err
	}

	return TreeInfo{
		Type:     "Tree",
		SHA1:     sha1Hash,
		FileName: directoryName,
		FilePath: dirName,
		Perm:     040000,
	}, nil
}

func GetTree(indexFile *map[string]FileInfo, mainDir, dir string) (TreeInfo, error) {
	var dirContent []TreeInfo
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if info.IsDir() {
			if info.Name() == ".git" || info.Name() == ".lit" {
				return filepath.SkipDir
			}

			// skip if is it calling dir path
			if path == dir {
				return nil
			}

			SubtreeInfo, err := GetTree(indexFile, mainDir, path)

			if err != nil {
				return err
			}

			dirContent = append(dirContent, SubtreeInfo)

			return filepath.SkipDir
		}

		// check if it is present in the index file
		file, exists := (*indexFile)[path]
		if !exists {
			return nil
		}

		NewTreeInfo := TreeInfo{Perm: 100644, FileName: file.FileName, SHA1: file.SHA1, Type: "blob", Modified_at: file.FileModifiedAt.String(), FileSize: int(file.FileSize), FilePath: strings.Replace(file.FilePath, " ", "||", -1)}

		dirContent = append(dirContent, NewTreeInfo)

		return nil
	})
	if err != nil {
		return TreeInfo{}, err
	}

	MainTreeObject, err := createTreeObj(dirContent, mainDir, dir)
	if err != nil {
		return TreeInfo{}, err
	}

	return MainTreeObject, nil
}
