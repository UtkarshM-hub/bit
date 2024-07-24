package core

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/UtkarshM-hub/bit/internal/application/core/util"
)

// Restores the current state of the file with previous state
func Restore(paths []string, rootDirPath string, staged bool) error {
	past_file_path := path.Join(rootDirPath, "./.bit/past")
	index_file_path := path.Join(rootDirPath, "./.bit/index")

	if err := util.DoesExists(past_file_path); err != nil && !staged {
		fmt.Println(past_file_path, ": commit object does not exist")
		return err
	}

	index_mp := GetIndexFileContent(index_file_path)
	past_mp := GetIndexFileContent(past_file_path)

	// range over the paths provided to the restore command
	for _, v := range paths {
		// check if it exists in the previous commit or not
		pastEntry, existPast := past_mp[v]
		if !existPast && !staged {
			continue
		}

		// if no blob object is present for the current file in the previous commit
		// and wants to unstage the changes
		if !existPast && staged {
			delete(index_mp, v)
			continue
		}

		// Just in case if the user has deleted the folder structure of the file
		// it will create the required folder
		var temp FileInfo = index_mp[v]

		// if the use wants to discard the changed
		if !staged {

			// create directory just in case it doesn't exist
			DirectoyPath := filepath.Dir(pastEntry.FilePath)
			err := util.DoesExists(DirectoyPath)
			if err != nil {
				err = os.MkdirAll(DirectoyPath, 0777)
				if err != nil {
					fmt.Println("Error while creating directory", err)
				}
			}

			// write file and change the entry in the index as well
			objectFilePath := path.Join(rootDirPath, "./.bit/objects/", pastEntry.SHA1[:2], pastEntry.SHA1[2:])

			// decompress and save the file at that particular location
			err = DecompressAndSaveFile(objectFilePath, pastEntry.FilePath)
			if err != nil {
				log.Fatal(err)
			}

			// modify the entries in current index file
			// mark it as committed as the changes are discarded
			temp.CommitStatus = "C"
			timeLayout := "2006-01-02 15:04:05.999999999 -0700 MST"
			formatedTime, _ := time.Parse(timeLayout, time.Now().Format(timeLayout))
			temp.FileModifiedAt = formatedTime
			index_mp[v] = temp
		} else {

			// mark as uncommited as the changes are still there and
			// it is just unstaged
			temp.CommitStatus = "c"

			// maintaining the past entry file status
			// keeping these things same as previous entry is important because
			// the status command is implemented in such a way that marking following values to updated will affect the expected results
			temp.FileModifiedAt = pastEntry.FileModifiedAt
			temp.FileStatus = pastEntry.FileStatus
			temp.SHA1 = pastEntry.SHA1
			index_mp[v] = temp
		}
	}

	// write the changes to current index
	writeToIndex(index_mp, index_file_path)

	return nil
}
