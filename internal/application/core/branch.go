package core

import (
	"fmt"
	"os"
	"path/filepath"
)

func CreateBranch(pathToLit,branchname string) error {
	// create file in logs/refs/heads/ folder
	branch_logfile_path:=filepath.Join(pathToLit,"/.lit/logs/refs/heads/"+branchname)
	file,err:=os.Create(branch_logfile_path)
	if err!=nil{
		return err
	}
	defer file.Close()
	// create file in refs/heads/ folder
	branch_reffile_path:=filepath.Join(pathToLit,"/.lit/refs/heads/"+branchname)
	file,err=os.Create(branch_reffile_path)
	if err!=nil{
		return err
	}
	defer file.Close()
	return nil
}

func ChangeActiveBranch(pathToLit,branchname string) error {

	HEAD_file_path:=filepath.Join(pathToLit,"/.lit/HEAD")

	// change the file content
	err:=os.WriteFile(HEAD_file_path,[]byte(fmt.Sprintf("ref: /refs/heads/%v",branchname)),0644)
	if err!=nil{
		return err
	}
	return nil
}