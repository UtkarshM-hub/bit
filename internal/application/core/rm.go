package core

import (
	"log"
	"path"

	"github.com/UtkarshM-hub/bit/internal/application/core/util"
)

func RemoveFilesFromStagingArea(bitDirPath string, files []string) {

	indexFilePath := path.Join(bitDirPath, "./.bit/index")

	err := util.DoesExists(indexFilePath)
	if err != nil {
		log.Fatal("Index file does not exist!")
		return
	}

	mp := GetIndexFileContent(indexFilePath)

	var removables []string

	if len(files) > 0 {
		for _, v := range files {
			val, exists := mp[v]
			// remove entry from current index
			delete(mp, v)

			// remove object file from objects directory
			if exists {
				objectToRemove := path.Join(bitDirPath, "./.bit/objects/", val.SHA1[:2])
				removables = append(removables, objectToRemove)
			}
		}
	} else {
		// include all the files
		for _, v := range mp {
			objectToRemove := path.Join(bitDirPath, "./.bit/objects/", v.SHA1[:2])
			removables = append(removables, objectToRemove)
		}

		// make map empty inorder to make the index file empty
		mp = make(map[string]FileInfo)
	}

	util.RemoveDirectories(removables)

	writeToIndex(mp, indexFilePath)
}
