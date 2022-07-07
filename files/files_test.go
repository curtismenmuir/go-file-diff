package files

import (
	"bufio"
	"errors"
	"io"
	"io/fs"
	"os"
	"testing"

	"github.com/curtismenmuir/go-file-diff/constants"
	"github.com/stretchr/testify/require"
)

const (
	fileName     string = "some-file.txt"
	errorMessage string = "Some Error"
	testOutput   string = "Testing `write to file` for now.....\n!\"£$%^&*(){}:@~>?<,./;'#[]\n\nFile signature coming soon!\n"
)

// Mock for fs.FileInfo interface
type fileInfoMock struct {
	// Include FileInfo props to fulfill interface
	os.FileInfo
	// Set test props
	isDir bool
}

// Overwrite fileInfoMock.IsDir() to consider test prop
func (m fileInfoMock) IsDir() bool { return m.isDir }

// Mock for Writer interface
type writerMock struct {
	// Include io.Writer props to fulfill mock for bufio.NewWriter
	io.Writer
	// Set test props
	isError bool
}

// Overwrite writerMock.WriteByte() to consider test prop
func (w writerMock) WriteByte(c byte) error {
	// Throw error if isError set
	if w.isError {
		return errors.New(errorMessage)
	}

	return nil
}

// Implement writerMock.Flush()
func (w writerMock) Flush() error { return nil }

func TestCreateFolder(t *testing.T) {
	t.Run("should return `nil` when folder created successfully", func(t *testing.T) {
		// Mock
		mkdir = func(name string, perm fs.FileMode) error {
			return nil
		}

		// Run
		err := createFolder(fileName)
		// Verify
		require.Equal(t, nil, err)
	})

	t.Run("should return `error` when unable to create folder", func(t *testing.T) {
		// Setup
		expectedError := errors.New(constants.UnableToCreateNewFolderError)
		// Mock
		mkdir = func(name string, perm fs.FileMode) error {
			return errors.New(errorMessage)
		}

		// Run
		err := createFolder(fileName)
		// Verify
		require.Equal(t, expectedError, err)
	})
}

func TestCreateWriter(t *testing.T) {
	t.Run("should return file writer", func(t *testing.T) {
		// Setup
		file := os.File{}
		writer := writerMock{}
		fileWriter := bufio.NewWriter(writer)
		// Mock
		newWriter = func(w io.Writer) *bufio.Writer {
			return fileWriter
		}

		// Run
		result := createWriter(&file)
		// Verify
		require.Equal(t, fileWriter, result)
	})
}

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

func TestVerifyOutputDirExists(t *testing.T) {
	t.Run("should return `nil` when Outputs folder already exists", func(t *testing.T) {
		// Mock
		getFileInfo = func(name string) (fs.FileInfo, error) {
			fileInfo := fileInfoMock{isDir: false}
			return fileInfo, nil
		}

		// Run
		result := verifyOutputDirExists()
		// Verify
		require.Equal(t, nil, result)
	})

	t.Run("should return `nil` when successfully created Outputs folder", func(t *testing.T) {
		// Mock
		getFileInfo = func(name string) (fs.FileInfo, error) {
			return nil, errors.New(errorMessage)
		}

		checkNotExists = func(err error) bool {
			return true
		}

		mkdir = func(name string, perm fs.FileMode) error {
			return nil
		}

		// Run
		result := verifyOutputDirExists()
		// Verify
		require.Equal(t, nil, result)
	})

	t.Run("should return error when unable to verify if Output dir exists", func(t *testing.T) {
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
		result := verifyOutputDirExists()
		// Verify
		require.Equal(t, expectedError, result)
	})

	t.Run("should return `unable to create Output folder` error when unable to create folder dir", func(t *testing.T) {
		// Setup
		expectedError := errors.New(constants.UnableToCreateOutputsFolder)
		// Mock
		getFileInfo = func(name string) (fs.FileInfo, error) {
			return nil, errors.New(errorMessage)
		}

		checkNotExists = func(err error) bool {
			return true
		}

		mkdir = func(name string, perm fs.FileMode) error {
			return errors.New(errorMessage)
		}

		// Run
		result := verifyOutputDirExists()
		// Verify
		require.Equal(t, expectedError, result)
	})
}

func TestWriteToFile(t *testing.T) {
	t.Run("should return `nil` when Output dir exists and successfully written output to file", func(t *testing.T) {
		// Setup
		file := os.File{}
		output := []byte(testOutput)
		// Mock
		getFileInfo = func(name string) (fs.FileInfo, error) {
			fileInfo := fileInfoMock{isDir: false}
			return fileInfo, nil
		}

		createFile = func(name string) (*os.File, error) {
			return &file, nil
		}

		createNewWriter = func(file *os.File) Writer {
			writer := writerMock{isError: false}
			return writer
		}

		mkdir = func(name string, perm fs.FileMode) error {
			return nil
		}

		// Run
		result := WriteToFile(fileName, output)
		// Verify
		require.Equal(t, nil, result)
	})

	t.Run("should return `nil` when created Output dir and successfully written output to file", func(t *testing.T) {
		// Setup
		file := os.File{}
		output := []byte(testOutput)
		// Mock
		getFileInfo = func(name string) (fs.FileInfo, error) {
			return nil, errors.New(errorMessage)
		}

		checkNotExists = func(err error) bool {
			return true
		}

		mkdir = func(name string, perm fs.FileMode) error {
			return nil
		}

		createFile = func(name string) (*os.File, error) {
			return &file, nil
		}

		createNewWriter = func(file *os.File) Writer {
			writer := writerMock{isError: false}
			return writer
		}

		// Run
		result := WriteToFile(fileName, output)
		// Verify
		require.Equal(t, nil, result)
	})

	t.Run("should return error when unable to verify if Output dir exists", func(t *testing.T) {
		// Setup
		output := []byte(testOutput)
		expectedError := errors.New(constants.UnableToCheckFileFolderExistsError)
		// Mock
		getFileInfo = func(name string) (fs.FileInfo, error) {
			return nil, errors.New(errorMessage)
		}

		checkNotExists = func(err error) bool {
			return false
		}

		// Run
		result := WriteToFile(fileName, output)
		// Verify
		require.Equal(t, expectedError, result)
	})

	t.Run("should return `unable to create Sig file` error when unable to create file", func(t *testing.T) {
		// Setup
		file := os.File{}
		output := []byte(testOutput)
		expectedError := errors.New(constants.UnableToCreateSignatureFile)
		// Mock
		getFileInfo = func(name string) (fs.FileInfo, error) {
			fileInfo := fileInfoMock{isDir: false}
			return fileInfo, nil
		}

		createFile = func(name string) (*os.File, error) {
			return &file, errors.New(errorMessage)
		}

		// Run
		result := WriteToFile(fileName, output)
		// Verify
		require.Equal(t, expectedError, result)
	})

	t.Run("should return `unable to write to Sig file` error when unable to write to file", func(t *testing.T) {
		// Setup
		file := os.File{}
		output := []byte(testOutput)
		expectedError := errors.New(constants.UnableToWriteToSignatureFile)
		// Mock
		getFileInfo = func(name string) (fs.FileInfo, error) {
			fileInfo := fileInfoMock{isDir: false}
			return fileInfo, nil
		}

		createFile = func(name string) (*os.File, error) {
			return &file, nil
		}

		createNewWriter = func(file *os.File) Writer {
			writer := writerMock{isError: true}
			return writer
		}

		mkdir = func(name string, perm fs.FileMode) error {
			return nil
		}

		// Run
		result := WriteToFile(fileName, output)
		// Verify
		require.Equal(t, expectedError, result)
	})
}
