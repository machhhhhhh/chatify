package utils

import (
	"errors"
	"os"
	"strings"
)

func FileExists(file_path string) (bool, error) {
	_, err := os.Stat(file_path)
	if os.IsNotExist(err) == true {
		return false, nil
	}

	return err == nil, err
}

func CheckFileExist(file_path string) (string, error) {
	is_exist, err := FileExists(strings.TrimSpace(file_path))
	if err != nil {
		return file_path, err
	}

	if is_exist != true {
		return file_path, errors.New("File not exist : " + file_path)
	}

	file_path = file_path[1:]
	return file_path, nil
}
