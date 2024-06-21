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
	CommitStatus   string
}

func GetIndexFileContent(path string) map[string]FileInfo {
	// Get the data from index file and store it in a map
	mp := make(map[string]FileInfo)

	file, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	// Get the data only if index file exist
	if err == nil {
		scanner := bufio.NewScanner(file)

		timeLayout := "2006-01-02 15:04:05.999999999 -0700 MST"

		for scanner.Scan() {
			line := scanner.Text()
			fields := strings.Split(line, " ")

			// util.go 2024-01-01 17:15:46.925855171 +0530 IST 1447 420 7baf657637a5c914f4e37cff2974941a5938ef9b /home/utkarsh/The Futher/UtkarshM-hub/Lit/internal/application/core/util/util.go N

			// replace || with space in filepath
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
				CommitStatus:   fields[10],
			}

			mp[key] = newEntry
		}
	}
	return mp
}

func GetFilesStatus(dir string) []FileInfo {
	var files []FileInfo

	// walk over the directory and get all the information about files
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return nil
		}
		if info.IsDir() {
			if info.Name() == ".git" || info.Name() == ".bit" {
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

		files = append(files, FileInfo{
			FileName:       info.Name(),
			FilePath:       path,
			FileSize:       uint64(info.Size()),
			FilePerm:       permissions,
			FileModifiedAt: modifiedTime,
			SHA1:           "",  // initially no values assigned
			FileStatus:     "N", // initially file is considered as new (N is used to represent new)
			CommitStatus:   "c", // initially small c represents un-commited files
		})

		return nil
	})
	if err != nil {
		fmt.Println(err.Error())
		return []FileInfo{}
	}
	return files
}

// Get the slice of tracked, modified, deleted and unracked files
// It takes the information about files in current director with index filepath
func GetStatus(files []FileInfo, path string) ([]FileInfo, []FileInfo, []FileInfo, []FileInfo, error) {
	// create four slices
	var untracked []FileInfo
	var modified []FileInfo
	var deleted []FileInfo
	var tracked []FileInfo

	mp := GetIndexFileContent(path)

	for _, v := range files {
		currentF, exists := mp[v.FilePath]
		// check if untracked
		if !exists {
			untracked = append(untracked, v)
			continue
		}

		// check time
		// if it is not commited then show it as tracked file
		if currentF.FileModifiedAt.Equal(v.FileModifiedAt) && currentF.CommitStatus != "C" {
			// In v all the values have status as 'N' so we'll have to take currentF
			tracked = append(tracked, currentF)
			delete(mp, v.FilePath)
			continue
		}

		// After commiting it should not show as tracked instead we want to show that the current directory is clean
		if currentF.FileModifiedAt.Equal(v.FileModifiedAt) && currentF.CommitStatus == "C" {
			delete(mp, v.FilePath)
			continue
		}

		// if time is different but we have to check the content itself inorder to  determine if the file is changed or not
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

		// if hash matches and not commited which means the content is not modified
		// but as it is not commited then it goes into tracked files
		if sha1HashWD == currentF.SHA1 && currentF.CommitStatus != "C" {
			tracked = append(tracked, currentF)
			delete(mp, v.FilePath)
			continue
		}

		// If the file is commited it should not be included in the file result to show the clean working tree
		if sha1HashWD == currentF.SHA1 && currentF.CommitStatus == "C" {
			delete(mp, v.FilePath)
			continue
		}

		// the purpose of adding this check was
		// we are ignoring the commited files from the starting so that they will not get included in the tracked files
		// but in that process they were getting added to modified instead
		// to prevent that behviour, following if condition is added for modified files

		if sha1HashWD != currentF.SHA1 {
			modified = append(modified, v)
			delete(mp, v.FilePath)
		}
	}

	// Get deleted
	for _, val := range mp {
		// remove those files which are deleted and are tracked or marked as deleted in index file
		if val.FileStatus == "D" {
			tracked = append(tracked, val)
			continue
		}

		// if not marked then include in deleted slice
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
