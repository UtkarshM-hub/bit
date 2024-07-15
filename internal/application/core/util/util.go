package util

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// get the pat of current directory
func GetCurrentDirectory() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return dir, nil
}

// join the path
func GetJoinedPaths(relativePath string) (string, error) {
	// Join the base path and relative path
	basePath, err := GetCurrentDirectory()
	if err != nil {
		return "", err
	}
	joinedPath := filepath.Join(basePath, relativePath)

	// Clean and normalize the path
	joinedPath, err = filepath.Abs(joinedPath)
	if err != nil {
		return "", err
	}

	return joinedPath, nil
}

// check if directory exists or not
func DoesExists(path string) error {
	_, err := os.Stat(path)

	// Check if the error is a "not exists" error which means it doesn't exist
	if os.IsNotExist(err) {
		return errors.New("the directory does not exist")
	} else if err != nil {
		// For other errors return the error
		return err
	}

	return nil
}

// Find the directory by moving towards parent
func FindDirectory(targetDir string) (string, error) {
	currentDir, err := os.Getwd() //get current working directory
	if err != nil {
		return "", err
	}

	for {
		// check if file exists in current dir
		if _, err := os.Stat(filepath.Join(currentDir, targetDir)); err == nil {
			return currentDir, nil
		}

		// Move to the parent directory
		parentDir := filepath.Dir(currentDir)
		// If we are already at the root directory, break the loop
		if parentDir == currentDir {
			break
		}

		currentDir = parentDir
	}
	return "", fmt.Errorf("not a bit directory 🧯")
}

// read file and return the string data
func ReadFile(filepath string) (string, error) {
	data, err := os.ReadFile(filepath)
	return string(data), err
}

// check if directory is empty
func ISDirectoryEmpty(dirPath string) (bool, error) {
	f, err := os.Open(dirPath)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1)
	if err == nil {
		// Directory is not empty
		return false, nil
	}
	if err == io.EOF {
		// Directory is empty
		return true, nil
	}
	return false, err
}

// remove folder with childrens
func RemoveDirectories(paths []string) error {

	for _, v := range paths {
		if err := DoesExists(v); err != nil {
			return fmt.Errorf("folder does not exists")
		}
		err := os.RemoveAll(v)
		if err != nil {
			return err
		}
	}
	return nil
}
