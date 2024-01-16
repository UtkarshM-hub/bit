package core

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"
)

type FileInfo struct {
	FileName       string
	FilePath       string
	FileSize       uint64
	FilePerm       uint32
	FileModifiedAt time.Time
	SHA1           string
	FileStatus     string
}

func GetFilesStatus(dir string) []FileInfo {
	var files []FileInfo
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return nil
		}
		if info.IsDir() {
			if info.Name() == ".git" || info.Name() == ".lit" {
				return filepath.SkipDir
			}
			// fmt.Println("Directory:", info.Name())
			return nil
		}

		// Type assert info.Sys() to *syscall.Stat_t to get access to specific info like permissions
		stat, ok := info.Sys().(*syscall.Stat_t)
		if !ok {
			return fmt.Errorf("Sys() did not return *syscall.Stat_t")
		}

		// Extract file permissions and modified date from *syscall.Stat_t
		permissions := uint32(os.ModePerm) & uint32(stat.Mode)
		modifiedTime := info.ModTime()

		// STORE ENTIES IN SORTED ORDER
		files = append(files, FileInfo{FileName: info.Name(), FilePath: path, FileSize: uint64(info.Size()), FilePerm: permissions, FileModifiedAt: modifiedTime, SHA1: "", FileStatus: "N"})

		return nil
	})
	if err != nil {
		fmt.Println(err.Error())
		return []FileInfo{}
	}
	return files
}

// Get the slice of modified, deleted and unracked files
func GetStatus(files []FileInfo, path string) ([]FileInfo, []FileInfo, []FileInfo, []FileInfo, error) {
	// create three arrays
	var untracked []FileInfo
	var modified []FileInfo
	var deleted []FileInfo
	var tracked []FileInfo

	// Get the data from index file and store it in a map
	mp := make(map[string]FileInfo)

	file, err := os.Open(path)
	defer file.Close()

	// Get the data only if index file exist
	if err == nil {
		scanner := bufio.NewScanner(file)

		timeLayout := "2024-01-01 17:15:46.925855171 +0530 IST"

		for scanner.Scan() {
			line := scanner.Text()
			fields := strings.Split(line, " ")

			// util.go 2024-01-01 17:15:46.925855171 +0530 IST 1447 420 7baf657637a5c914f4e37cff2974941a5938ef9b /home/utkarsh/The Futher/UtkarshM-hub/Lit/internal/application/core/util/util.go N

			// fmt.Printf("Name: %v\n Time: %v %v %v %v\nSize: %v\nPerm: %v\nSHA1: %v\nPath: %v\n",fields[0],fields[1],fields[2],fields[3],fields[4],fields[5],fields[6],fields[7],fields[8])

			key := strings.Replace(fields[8], "||", " ", -1)

			T, _ := time.Parse(timeLayout, fields[1]+" "+fields[2]+" "+fields[3]+" "+fields[4])

			size, _ := strconv.Atoi(fields[5])
			perm, _ := strconv.Atoi(fields[6])
			newEntry := FileInfo{
				FileName:       fields[0],
				FilePath:       strings.Replace(fields[8], "||", " ", -1),
				FileSize:       uint64(size),
				FileModifiedAt: T,
				FilePerm:       uint32(perm),
				SHA1:           fields[7],
				FileStatus:     fields[9],
			}

			mp[key] = newEntry
		}
	}

	for _, v := range files {
		currentF, exists := mp[v.FilePath]
		// check if untracked
		if !exists {
			untracked = append(untracked, v)
			continue
		}
		// check time
		if currentF.FileModifiedAt.Equal(v.FileModifiedAt) {
			tracked = append(tracked, v)
			delete(mp, v.FilePath)
			continue
		}

		// check hash value
		content, err := os.ReadFile(v.FilePath)
		if err != nil {
			fmt.Println(err.Error())
			return []FileInfo{}, []FileInfo{}, []FileInfo{}, []FileInfo{}, err
		}

		// Construct the header with object type+filesize+\0
		header := "blob" + " " + strconv.Itoa(int(v.FileSize)) + "\\0"

		// Find sha1 hash of the header+content
		sha1HashWD := calculateSHA1(header + string(content))

		if sha1HashWD == currentF.SHA1 {
			tracked = append(tracked, v)
			delete(mp, v.FilePath)
			continue
		}
		modified = append(modified, currentF)
		delete(mp, v.FilePath)
	}

	// Get deleted
	for _, val := range mp {
		deleted = append(deleted, val)
	}
	// fmt.Println("-----")
	// fmt.Println(deleted)
	// fmt.Println("-----")
	// fmt.Println(untracked)
	// fmt.Println("-----")
	// fmt.Println(modified)
	// fmt.Println("-----")
	// fmt.Println(tracked)
	return tracked, untracked, modified, deleted, nil
}
