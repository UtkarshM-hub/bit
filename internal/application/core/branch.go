package core

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	Queue "github.com/UtkarshM-hub/bit/internal/application/core/Structs/Queue"
	"github.com/UtkarshM-hub/bit/internal/application/core/util"
	color "github.com/gookit/color"
)

// Function to create new branch
func CreateBranch(pathToLit, branchname string) error {
	// create file in logs/refs/heads/ folder
	branch_logfile_path := filepath.Join(pathToLit, "/.bit/logs/refs/heads/"+branchname)
	file, err := os.Create(branch_logfile_path)
	if err != nil {
		return err
	}
	defer file.Close()

	// create file in refs/heads/ folder
	branch_reffile_path := filepath.Join(pathToLit, "/.bit/refs/heads/"+branchname)
	file, err = os.Create(branch_reffile_path)
	if err != nil {
		return err
	}
	defer file.Close()
	return nil
}

// Function which returns the current active branch
func CurrentActiveBranch(pathToLit string) (string, error) {
	// get current active branch
	HEAD_file_path := filepath.Join(pathToLit, "/.bit/HEAD")
	data, err := os.ReadFile(HEAD_file_path)
	if err != nil {
		return "", err
	}

	data_arr := strings.Split(string(data), "/")
	current_active_branch := data_arr[len(data_arr)-1]
	return current_active_branch, nil
}

// List all the branches and show the current active branch
func ListBranches(pathToLit string) error {
	var branches []string

	refs_file_path := filepath.Join(pathToLit, "/.bit/refs/heads")

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

// Perform checkout operation to switch to other branch
func Checkout(PathToLit, BranchName string) {
	branch_pointer_file_path := filepath.Join(PathToLit, "/.bit/refs/heads/"+BranchName)
	index_file_path := filepath.Join(PathToLit, "./.bit/index")
	past_file_path := filepath.Join(PathToLit, "./.bit/past")

	// Get the commit id of the branch to checkout
	commit_hash, err := os.ReadFile(branch_pointer_file_path)
	if err != nil {
		fmt.Println(err)
	}

	// if branch is new and there is no commit object present
	// Change only the necessary info if the branch is new
	// No need to perform deletion and all
	if len(commit_hash) == 0 {
		fmt.Println("empty")
		// change active branch
		err = ChangeActiveBranch(PathToLit, BranchName)
		if err != nil {
			fmt.Println(err)
		}
		return
	}

	// if commit object is present
	// we will have to perform necessary creation and deletion of files accordingly

	// Get the hash value of tree object so that we can build an index out of it
	commit_object_file_path := filepath.Join(PathToLit, "/.bit/objects/", string(commit_hash[:2])+"/"+string(commit_hash[2:]))
	commit_object_data, _ := DecompressFile(commit_object_file_path)
	first_line := strings.Split(commit_object_data, "\n")[0]
	tree_object_hash := strings.Split(first_line, " ")[1]

	// Build new index for the specific branch
	New_Branch_Index, err := GenerateIndex(PathToLit, tree_object_hash)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Get the content of current index file
	Current_Branch_Index := GetIndexFileContent(index_file_path)

	// Write the new index file content to index file
	writeToIndex(New_Branch_Index, index_file_path)
	
	// write the new index file content to past file
	writeToIndex(New_Branch_Index, past_file_path)

	// Do the switching operation (creation and deltion of files)
	Switch(PathToLit, Current_Branch_Index, New_Branch_Index)

	err = ChangeActiveBranch(PathToLit, BranchName)
	if err != nil {
		fmt.Println(err)
	}

	// // After switching branches remove these ghost directories which contains no files
	// DeleteEmptyDir(PathToLit)
}

func ChangeActiveBranch(pathToLit, branchname string) error {

	HEAD_file_path := filepath.Join(pathToLit, "/.bit/HEAD")
	logs_HEAD_file_path := filepath.Join(pathToLit, "/.bit/logs/HEAD")

	prev_active_branch, err := CurrentActiveBranch(pathToLit)

	prev_active_branch_path := filepath.Join(pathToLit, "/.bit/refs/heads/"+prev_active_branch)
	current_active_branch_path := filepath.Join(pathToLit, "/.bit/refs/heads/"+branchname)

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

// generate index file for the new branch
func GenerateIndex(PathToLit, Tree_Object_Hash string) (map[string]FileInfo, error) {
	TreeQueue := Queue.Queue{}
	New_Branch_Index := map[string]FileInfo{}
	timeLayout := "2006-01-02 15:04:05.999999999 -0700 MST"

	// insert the main tree hash
	TreeQueue.Enqueue(Tree_Object_Hash)

	for !TreeQueue.IsEmpty() {
		// take the hash out and get the content of Tree Object
		val := TreeQueue.Dequeue()

		// Decompress the tree object and get the info
		main_tree_object_file_path := filepath.Join(PathToLit, "/.bit/objects/", string(val[:2])+"/"+string(val[2:]))
		Main_Tree, err := DecompressFile(main_tree_object_file_path)

		if err != nil {
			fmt.Println(err.Error())
			return nil, err
		}

		content := strings.Split(Main_Tree, "\n")

		for _, entry := range content {
			parameters := strings.Split(entry, " ")

			if parameters[0] == "Tree" {
				TreeQueue.Enqueue(parameters[3])
				continue
			}

			T, _ := time.Parse(timeLayout, parameters[1]+" "+parameters[2]+" "+parameters[3]+" "+parameters[4])

			size, _ := strconv.Atoi(parameters[8])
			perm, _ := strconv.Atoi(parameters[7])
			path := strings.Replace(parameters[9], "||", " ", -1)

			newEntry := FileInfo{
				FileName:       filepath.Base(path),
				FilePath:       path,
				FileSize:       uint64(size),
				FileModifiedAt: T,
				FilePerm:       uint32(perm),
				SHA1:           parameters[6],
				FileStatus:     "N",
				CommitStatus:   "C",
			}

			New_Branch_Index[path] = newEntry
		}
	}
	return New_Branch_Index, nil
}

// Performs actual switching operation of branch command
func Switch(dir string, currentB, newB map[string]FileInfo) error {
	fmt.Println("switch start")

	// Loop over new Branch
	for _, v := range newB {
		fileInfo, exists := currentB[v.FilePath]
		fmt.Println("switch mid", v.FilePath)

		if exists {
			delete(currentB, v.FilePath)
			// skip if file already exists and hash value is same
			if fileInfo.SHA1 == v.SHA1 {
				fmt.Println("switch cont", v.FilePath)
				delete(newB, v.FilePath)
				continue
			} else {
				// delete the existing file and new file will be created in the following code
				err := os.Remove(fileInfo.FilePath)
				fmt.Println("switch rm", v.FilePath)
				if err != nil {
					fmt.Println(err.Error())
					return err
				}
			}
		}

		// create folder if it doesn't exists, in-order to create file at that specific path
		DirectoyPath := filepath.Dir(v.FilePath)
		err := util.DoesExists(DirectoyPath)
		if err != nil {
			err = os.MkdirAll(DirectoyPath, 0777)
			if err != nil {
				fmt.Println("Error while creating directory", err)
			}
		}

		fmt.Println("switch creating", v.FilePath)

		// Now as we have our directory structure for the current file setup
		// we can proceed to create the file by decompressing it and storing it at a particular location
		inputFilePath := filepath.Join(dir, "/.bit/objects", string(v.SHA1[:2])+"/"+string(v.SHA1[2:]))

		fmt.Println(inputFilePath, v.FilePath)

		err = DecompressAndSaveFile(inputFilePath, v.FilePath)

		if err != nil {
			fmt.Println("this is it", err.Error())
		}
		delete(newB, v.FilePath)
	}

	// Remove the remaining files which were present in previous branch
	//  as they don't belong to the current branch
	for _, v := range currentB {
		err := os.Remove(v.FilePath)
		if err != nil {
			return err
		}

		parent := path.Dir(v.FilePath)

		empty, err := util.ISDirectoryEmpty(parent)
		if err != nil {
			return err
		}

		if empty {
			err = os.Remove(parent)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func DecompressFile(inputFilePath string) (string, error) {

	// Read the compressed data from the input file
	compressedData, err := os.ReadFile(inputFilePath)
	if err != nil {
		return "", err
	}

	// Create a zlib reader with default compression level
	reader, err := zlib.NewReader(bytes.NewReader(compressedData))
	if err != nil {
		return "", err
	}

	// Create a buffer to store the decompressed data
	var decompressedBuffer bytes.Buffer

	// Copy the data from the zlib reader to the decompressed buffer
	_, err = io.Copy(&decompressedBuffer, reader)
	if err != nil {
		return "", err
	}

	// Close the zlib reader
	reader.Close()

	// // Write the decompressed data to the output file
	// err = os.WriteFile(outputFilePath, decompressedBuffer.Bytes(), 0644)
	// if err != nil {
	// 	return err
	// }

	return string(decompressedBuffer.String()), nil
}

func DecompressAndSaveFile(inputFilePath, outputFilePath string) error {

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
