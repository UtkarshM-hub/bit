package core

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

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

func ChangeActiveBranch(pathToLit, branchname string) error {

	HEAD_file_path := filepath.Join(pathToLit, "/.lit/HEAD")

	// change the file content
	err := os.WriteFile(HEAD_file_path, []byte(fmt.Sprintf("ref: /refs/heads/%v", branchname)), 0644)
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

	// get current active branch
	HEAD_file_path := filepath.Join(pathToLit, "/.lit/HEAD")
	data, err := os.ReadFile(HEAD_file_path)
	if err != nil {
		return err
	}

	data_arr := strings.Split(string(data), "/")
	current_active_branch := data_arr[len(data_arr)-1]

	// print to terminal
	for _, v := range branches {
		if v == current_active_branch {
			color.Green.Printf("* %v\n", v)
			continue
		}
		fmt.Printf("  %v\n",v)
	}
	return nil
}
