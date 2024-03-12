package core

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/UtkarshM-hub/Lit/internal/application/core/util"
	color "github.com/gookit/color"
)

func CreateBranch(pathToLit, branchname string) error {
	// create file in logs/refs/heads/ folder
	branch_logfile_path := filepath.Join(pathToLit, "/.lit/logs/refs/heads/"+branchname)
	file, err := os.Create(branch_logfile_path)
	if err != nil {
		return err
	}
	defer file.Close()
	// create file in refs/heads/ folder
	branch_reffile_path := filepath.Join(pathToLit, "/.lit/refs/heads/"+branchname)
	file, err = os.Create(branch_reffile_path)
	if err != nil {
		return err
	}
	defer file.Close()
	return nil
}

func CurrentActiveBranch(pathToLit string) (string, error) {
	// get current active branch
	HEAD_file_path := filepath.Join(pathToLit, "/.lit/HEAD")
	data, err := os.ReadFile(HEAD_file_path)
	if err != nil {
		return "", err
	}

	data_arr := strings.Split(string(data), "/")
	current_active_branch := data_arr[len(data_arr)-1]
	return current_active_branch, nil
}

func ChangeActiveBranch(pathToLit, branchname string) error {

	HEAD_file_path := filepath.Join(pathToLit, "/.lit/HEAD")
	logs_HEAD_file_path := filepath.Join(pathToLit, "/.lit/logs/HEAD")

	prev_active_branch, err := CurrentActiveBranch(pathToLit)

	prev_active_branch_path := filepath.Join(pathToLit, "/.lit/refs/heads/"+prev_active_branch)
	current_active_branch_path := filepath.Join(pathToLit, "/.lit/refs/heads/"+branchname)

	if err != nil {
		return err
	}

	// change the file content
	err = os.WriteFile(HEAD_file_path, []byte(fmt.Sprintf("ref: /refs/heads/%v", branchname)), 0644)
	if err != nil {
		return err
	}

	// add logs into logs file
	prev_commit_hash, err := util.ReadFile(prev_active_branch_path)
	if err != nil {
		return err
	}
	current_commit_hash, err := util.ReadFile(current_active_branch_path)
	if err != nil {
		return err
	}

	if prev_commit_hash == "" {
		prev_commit_hash = "0000000000000000000000000000000000000000"
	}
	if current_commit_hash == "" && prev_commit_hash == "" {
		current_commit_hash = "0000000000000000000000000000000000000000"
	} else if current_commit_hash == "" {
		current_commit_hash = prev_commit_hash
	}

	commiter := strings.Replace("Utkarsh Mandape", " ", "||", -1)
	msg := strings.Replace(fmt.Sprintf("checkout: moving from %v to %v", prev_active_branch, branchname), " ", "||", -1)
	NewData := fmt.Sprintf("%v %v %v %v %v %v", prev_commit_hash, current_commit_hash, commiter, "utmandape4@gmail.com", time.Now().String(), msg)

	err = LogsBranchChange(logs_HEAD_file_path, NewData)

	return err
}

func LogsBranchChange(logFilePath, msg string) error {
	err := util.DoesExists(logFilePath)
	if err != nil {
		os.Create(logFilePath)
	}
	data, err := os.ReadFile(logFilePath)
	if err != nil {
		return err
	}

	// modify the commit message and username in-order to replace space with ||
	data_string := strings.Split(string(data), "\n")
	data_string = append(data_string, msg)
	err = os.WriteFile(logFilePath, []byte(strings.Join(data_string, "\n")), 0644)
	if err != nil {
		return err
	}
	return nil
}

func ListBranches(pathToLit string) error {
	var branches []string

	refs_file_path := filepath.Join(pathToLit, "/.lit/refs/heads")

	// get all branches
	filepath.WalkDir(refs_file_path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Println(err)
			return err
		}

		if !d.IsDir() {
			branches = append(branches, d.Name())
		}
		return nil
	})

	current_active_branch, err := CurrentActiveBranch(pathToLit)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// print to terminal
	for _, v := range branches {
		if v == current_active_branch {
			color.Green.Printf("* %v\n", v)
			continue
		}
		fmt.Printf("  %v\n", v)
	}
	return nil
}

func Checkout(PathToLit, BranchName string){
	// branch_pointer_file_path:=filepath.Join(PathToLit,"/.lit/refs/heads/"+BranchName)
	
}

func DecompressFile(inputFilePath, outputFilePath string) error {

	// Read the compressed data from the input file
	compressedData, err := os.ReadFile(inputFilePath)
	if err != nil {
		return err
	}

	// Create a zlib reader with default compression level
	reader, err := zlib.NewReader(bytes.NewReader(compressedData))
	if err != nil {
		return err
	}

	// Create a buffer to store the decompressed data
	var decompressedBuffer bytes.Buffer

	// Copy the data from the zlib reader to the decompressed buffer
	_, err = io.Copy(&decompressedBuffer, reader)
	if err != nil {
		return err
	}

	// Close the zlib reader
	reader.Close()

	// Write the decompressed data to the output file
	err = os.WriteFile(outputFilePath, decompressedBuffer.Bytes(), 0644)
	if err != nil {
		return err
	}

	return nil
}
