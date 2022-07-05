package files

import (
	"bufio"
	"errors"
	"io/fs"
	"os"
	"testing"

	"github.com/curtismenmuir/go-file-diff/constants"
	"github.com/stretchr/testify/require"
)

const (
	fileName     string = "some-file.txt"
	errorMessage string = "Some Error"
)

type fileInfoMock struct {
	// Include FileInfo props to fulfill interface
	os.FileInfo
	// Set test props
	isDir bool
}

// Overwrite fileInfoMock IsDir() to consider test prop
func (m fileInfoMock) IsDir() bool { return m.isDir }

func TestDoesExists(t *testing.T) {
	t.Run("should return `true, nil` when file exists", func(t *testing.T) {
		// Mock
		getFileInfo = func(name string) (fs.FileInfo, error) {
			fileInfo := fileInfoMock{isDir: false}
			return fileInfo, nil
		}

		// Run
		result, err := doesExist(fileName, true)
		// Verify
		require.Equal(t, true, result)
		require.Equal(t, nil, err)
	})

	t.Run("should return `true, nil` when folder exists", func(t *testing.T) {
		// Mock
		getFileInfo = func(name string) (fs.FileInfo, error) {
			fileInfo := fileInfoMock{isDir: true}
			return fileInfo, nil
		}

		// Run
		result, err := doesExist(fileName, false)
		// Verify
		require.Equal(t, true, result)
		require.Equal(t, nil, err)
	})

	t.Run("should return `false, nil` when file does not exist", func(t *testing.T) {
		// Mock
		getFileInfo = func(name string) (fs.FileInfo, error) {
			return nil, errors.New(errorMessage)
		}

		checkNotExists = func(err error) bool {
			return true
		}

		// Run
		result, err := doesExist(fileName, true)
		// Verify
		require.Equal(t, false, result)
		require.Equal(t, nil, err)
	})

	t.Run("should return `false, error` when searching for file but found folder", func(t *testing.T) {
		// Setup
		expectedError := errors.New(constants.SearchingForFileButFoundDirError)
		// Mock
		getFileInfo = func(name string) (fs.FileInfo, error) {
			fileInfo := fileInfoMock{isDir: true}
			return fileInfo, nil
		}

		// Run
		result, err := doesExist(fileName, true)
		// Verify
		require.Equal(t, false, result)
		require.Equal(t, expectedError, err)
	})

	t.Run("should return `false, error` when error checking if file exists", func(t *testing.T) {
		// Setup
		expectedError := errors.New(constants.UnableToCheckFileFolderExistsError)
		// Mock
		getFileInfo = func(name string) (fs.FileInfo, error) {
			return nil, errors.New(errorMessage)
		}

		checkNotExists = func(err error) bool {
			return false
		}

		// Run
		result, err := doesExist(fileName, true)
		// Verify
		require.Equal(t, false, result)
		require.Equal(t, expectedError, err)
	})
}

func TestOpenFile(t *testing.T) {
	t.Run("should return file reader when successfully opened file", func(t *testing.T) {
		// Setup
		file := os.File{}
		expectedResult := bufio.NewReader(&file)
		// Mock
		getFileInfo = func(name string) (fs.FileInfo, error) {
			fileInfo := fileInfoMock{isDir: false}
			return fileInfo, nil
		}

		open = func(name string) (*os.File, error) {
			return &file, nil
		}

		// Run
		result, err := OpenFile(fileName)
		// Verify
		require.Equal(t, nil, err)
		require.Equal(t, expectedResult, result)
	})

	t.Run("should return error when unable to check if file exists", func(t *testing.T) {
		// Setup
		testError := errors.New(errorMessage)
		expectedResult := errors.New(constants.UnableToCheckFileFolderExistsError)
		// Mock
		getFileInfo = func(name string) (fs.FileInfo, error) {
			return nil, testError
		}

		checkNotExists = func(err error) bool {
			return false
		}

		open = func(name string) (*os.File, error) {
			return nil, testError
		}

		// Run
		_, err := OpenFile(fileName)
		// Verify
		require.Equal(t, expectedResult, err)
	})

	t.Run("should return `file does not exist` error when file does not exist", func(t *testing.T) {
		// Setup
		testError := errors.New(errorMessage)
		expectedResult := errors.New(constants.FileDoesNotExistError)
		// Mock
		getFileInfo = func(name string) (fs.FileInfo, error) {
			return nil, errors.New(errorMessage)
		}

		checkNotExists = func(err error) bool {
			return true
		}

		open = func(name string) (*os.File, error) {
			return nil, testError
		}

		// Run
		_, err := OpenFile(fileName)
		// Verify
		require.Equal(t, expectedResult, err)
	})

	t.Run("should return error when unable to open file", func(t *testing.T) {
		// Setup
		testError := errors.New(errorMessage)
		// Mock
		getFileInfo = func(name string) (fs.FileInfo, error) {
			fileInfo := fileInfoMock{isDir: false}
			return fileInfo, nil
		}

		open = func(name string) (*os.File, error) {
			return nil, testError
		}

		// Run
		_, err := OpenFile(fileName)
		// Verify
		require.Equal(t, testError, err)
	})
}
