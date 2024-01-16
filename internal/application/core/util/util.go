package util

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

func GetCurrentDirectory() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return dir, nil
}

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

func DoesExists(path string) error {
	_, err := os.Stat(path)

	// Check if the error is a "not exists" error which means it doesn't exist
	if os.IsNotExist(err) {
		return errors.New("The directory does not exist")
	} else if err != nil {
		// For other errors return the error
		return err
	}

	return nil
}

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
	return "", fmt.Errorf("Not a lit directory 🧯")
}
