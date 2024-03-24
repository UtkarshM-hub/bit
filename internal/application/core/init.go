package core

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// initializes the directory as bit directory
func Init(path string) bool {
	var wg sync.WaitGroup
	litPath := filepath.Join(path, ".bit")

	// Create the directory with 0700 permissions (read, write, execute for the owner only)
	// 4=read 2=write 1=execute=7
	// sequence is like owner group and other so 700 and first 0 is sticky bit
	// sticky bit: if a file is created inside that directory it gets permission same as its parent if the sticky bit is set or get's permissions of group by default
	err := os.Mkdir(litPath, 0700)
	if err != nil {
		return false
	}

	// Create other sub-directories
	directories := []string{"branches", "info", "logs", "hooks"}

	wg.Add(1)
	go func(directories []string){
		for _, v := range directories {
			dirPath := filepath.Join(litPath, v)
			err = os.Mkdir(dirPath, 0700)
			if err != nil {
				return 
			}
		}
		wg.Done()
	}(directories)

	subfiles:=[]string{"objects/info","objects/pack","refs/heads","refs/tags","logs/refs/heads"}
	go func(subfiles []string){
		for _,v:=range subfiles{
			fileName:=filepath.Join(litPath,v)
			err:=os.MkdirAll(fileName,0700)
			if err!=nil{
				return
			}
		}
	}(subfiles)

	wg.Wait()
	// Create file in .bit folder
	files:=[]string{"config","description","HEAD","index","packed-refs","/refs/heads/master","/logs/refs/heads/master"}
	for _,v:=range files{
		fileName:=filepath.Join(litPath,v)
		file, err := os.Create(fileName)
		
		// CHANGE THIS IN FUTURE AND CREATE A SEPARATE BLOCK OF CODE
		if v=="HEAD"{
			file.WriteString("ref: refs/heads/master")
		}

		if err != nil {
			fmt.Println(err.Error())
		}
		defer file.Close()
	}
	return true
}
