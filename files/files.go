package files

import (
	"bufio"
	"errors"
	"os"

	"github.com/curtismenmuir/go-file-diff/constants"
	"github.com/curtismenmuir/go-file-diff/utils"
)

var (
	logger         = utils.Logger
	open           = os.Open
	getFileInfo    = os.Stat
	checkNotExists = os.IsNotExist
)

// doesExist() checks if a file/folder exists and returns `true, nil` if specified file/folder is found
// When checking existence of a file, set isFile to true
// When checking existence of a folder dir, set isFile to false
// Function will return `false, nil` when file/folder does not exist
// Function will return `false, error` if an error is thrown
// Function will return `found folder` error if searching for file but found a folder dir
func doesExist(path string, isFile bool) (bool, error) {
	fileInfo, err := getFileInfo(path)
	if err != nil {
		if checkNotExists(err) {
			return false, nil
		}

		return false, errors.New(constants.UnableToCheckFileFolderExistsError)
	}

	if isFile && fileInfo.IsDir() {
		return false, errors.New(constants.SearchingForFileButFoundDirError)
	}

	return true, nil
}

// OpenFile() will attempt to open a local file and will return a file Reader when successful
// Function will catch and return error when unable to access specified file
// Function will return `file does not exist` error when specified file does not exist
func OpenFile(fileName string) (*bufio.Reader, error) {
	exists, err := doesExist(fileName, true)
	if err != nil {
		return nil, err
	} else if !exists {
		return nil, errors.New(constants.FileDoesNotExistError)
	}

	file, err := open(fileName)
	if err != nil {
		return nil, err
	}

	return bufio.NewReader(file), nil
}
