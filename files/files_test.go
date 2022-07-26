package files

import (
	"bufio"
	"encoding/gob"
	"errors"
	"io"
	"io/fs"
	"os"
	"testing"

	"github.com/curtismenmuir/go-file-diff/constants"
	"github.com/curtismenmuir/go-file-diff/models"
	"github.com/stretchr/testify/require"
)

const (
	fileName     string = "some-file.txt"
	errorMessage string = "Some Error"
	testOutput   string = "Testing `write to file` for now.....\n!\"£$%^&*(){}:@~>?<,./;'#[]\n\nFile signature coming soon!\n"
)

// Mock for Encoder interface
type encoderMock struct {
	// Set test props
	isError bool
}

// Overwrite encoderMock.Encode() to consider test prop
func (encoder encoderMock) Encode(e any) error {
	if encoder.isError {
		return errors.New(errorMessage)
	}

	return nil
}

// Mock for Decoder interface
type decoderMock struct {
	// Set test props
	isError bool
}

// Overwrite decoderMock.Decode() to consider test prop
func (decoder decoderMock) Decode(e any) error {
	if decoder.isError {
		return errors.New(errorMessage)
	}

	return nil
}

// Mock for fs.FileInfo interface
type fileInfoMock struct {
	// Include FileInfo props to fulfill interface
	os.FileInfo
	// Set test props
	isDir bool
}

// Overwrite fileInfoMock.IsDir() to consider test prop
func (m fileInfoMock) IsDir() bool { return m.isDir }

// Mock for io.Reader interface
type readerMock struct{}

// Overwrite readerMock.Read()
func (r readerMock) Read(p []byte) (n int, err error) {
	return 1, nil
}

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

func TestCreateDecoder(t *testing.T) {
	t.Run("should return file decoder", func(t *testing.T) {
		// Setup
		file := os.File{}
		reader := readerMock{}
		decoder := gob.NewDecoder(reader)
		// Mock
		newDecoder = func(r io.Reader) *gob.Decoder {
			return decoder
		}

		// Run
		result := createDecoder(&file)
		// Verify
		require.Equal(t, decoder, result)
	})
}

func TestCreateEncoder(t *testing.T) {
	t.Run("should return file encoder", func(t *testing.T) {
		// Setup
		file := os.File{}
		writer := writerMock{}
		encoder := gob.NewEncoder(writer)
		// Mock
		newEncoder = func(w io.Writer) *gob.Encoder {
			return encoder
		}

		// Run
		result := createEncoder(&file)
		// Verify
		require.Equal(t, encoder, result)
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

func TestOpenDelta(t *testing.T) {
	t.Run("should return `delta, nil` when successfully read Delta from file", func(t *testing.T) {
		// Setup
		file := os.File{}
		decoder := decoderMock{isError: false}
		// Decoder mock doesn't update provided pointer, so use empty struct for now
		// NOTE: Function will only return `err` as `nil` when successful
		expectedDelta := models.Delta{}
		var expectedError error = nil

		// Mock
		getFileInfo = func(name string) (fs.FileInfo, error) {
			fileInfo := fileInfoMock{isDir: false}
			return fileInfo, nil
		}

		open = func(name string) (*os.File, error) {
			return &file, nil
		}

		createNewDecoder = func(file *os.File) Decoder {
			return decoder
		}

		// Run
		delta, err := OpenDelta(fileName, false)
		// Verify
		require.Equal(t, expectedError, err)
		require.Equal(t, expectedDelta, delta)
	})

	t.Run("should return `emptyDelta, error` when unable to check if Delta file exists", func(t *testing.T) {
		// Setup
		testError := errors.New(errorMessage)
		expectedError := errors.New(constants.UnableToCheckFileFolderExistsError)
		expectedDelta := models.Delta{}
		// Mock
		getFileInfo = func(name string) (fs.FileInfo, error) {
			return nil, testError
		}

		checkNotExists = func(err error) bool {
			return false
		}

		// Run
		delta, err := OpenDelta(fileName, false)
		// Verify
		require.Equal(t, expectedError, err)
		require.Equal(t, expectedDelta, delta)
	})

	t.Run("should return `emptyDelta, DeltaFileDoesNotExistError` when Delta file does not exist", func(t *testing.T) {
		// Setup
		testError := errors.New(errorMessage)
		expectedError := errors.New(constants.DeltaFileDoesNotExistError)
		expectedDelta := models.Delta{}
		// Mock
		getFileInfo = func(name string) (fs.FileInfo, error) {
			return nil, testError
		}

		checkNotExists = func(err error) bool {
			return true
		}

		// Run
		delta, err := OpenDelta(fileName, false)
		// Verify
		require.Equal(t, expectedError, err)
		require.Equal(t, expectedDelta, delta)
	})

	t.Run("should return `emptyDelta, UnableToOpenDeltaFileError` when unable to open file", func(t *testing.T) {
		// Setup
		testError := errors.New(errorMessage)
		expectedError := errors.New(constants.UnableToOpenDeltaFileError)
		expectedDelta := models.Delta{}
		// Mock
		getFileInfo = func(name string) (fs.FileInfo, error) {
			fileInfo := fileInfoMock{isDir: false}
			return fileInfo, nil
		}

		open = func(name string) (*os.File, error) {
			return nil, testError
		}

		// Run
		delta, err := OpenDelta(fileName, false)
		// Verify
		require.Equal(t, expectedError, err)
		require.Equal(t, expectedDelta, delta)
	})

	t.Run("should return `emptyDelta, UnableToDecodeDeltaFromFileError` when unable to decode Delta from file", func(t *testing.T) {
		// Setup
		file := os.File{}
		decoder := decoderMock{isError: true}
		expectedError := errors.New(constants.UnableToDecodeDeltaFromFileError)
		expectedDelta := models.Delta{}
		// Mock
		getFileInfo = func(name string) (fs.FileInfo, error) {
			fileInfo := fileInfoMock{isDir: false}
			return fileInfo, nil
		}

		open = func(name string) (*os.File, error) {
			return &file, nil
		}

		createNewDecoder = func(file *os.File) Decoder {
			return decoder
		}

		// Run
		delta, err := OpenDelta(fileName, false)
		// Verify
		require.Equal(t, expectedError, err)
		require.Equal(t, expectedDelta, delta)
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

func TestOpenSignature(t *testing.T) {
	t.Run("should return `signature, nil` when successfully read Signature from file", func(t *testing.T) {
		// Setup
		file := os.File{}
		decoder := decoderMock{isError: false}
		// Decoder mock doesn't update provided pointer, so use empty struct for now
		// NOTE: Function will only return `err` as `nil` when successful
		expectedSignature := models.Signature{}
		var expectedError error = nil

		// Mock
		getFileInfo = func(name string) (fs.FileInfo, error) {
			fileInfo := fileInfoMock{isDir: false}
			return fileInfo, nil
		}

		open = func(name string) (*os.File, error) {
			return &file, nil
		}

		createNewDecoder = func(file *os.File) Decoder {
			return decoder
		}

		// Run
		signature, err := OpenSignature(fileName, false)
		// Verify
		require.Equal(t, expectedError, err)
		require.Equal(t, expectedSignature, signature)
	})

	t.Run("should return `emptySignature, error` when unable to check if Signature file exists", func(t *testing.T) {
		// Setup
		testError := errors.New(errorMessage)
		expectedError := errors.New(constants.UnableToCheckFileFolderExistsError)
		expectedSignature := models.Signature{}
		// Mock
		getFileInfo = func(name string) (fs.FileInfo, error) {
			return nil, testError
		}

		checkNotExists = func(err error) bool {
			return false
		}

		// Run
		signature, err := OpenSignature(fileName, false)
		// Verify
		require.Equal(t, expectedError, err)
		require.Equal(t, expectedSignature, signature)
	})

	t.Run("should return `emptySignature, SignatureFileDoesNotExistError` when Signature file does not exist", func(t *testing.T) {
		// Setup
		testError := errors.New(errorMessage)
		expectedError := errors.New(constants.SignatureFileDoesNotExistError)
		expectedSignature := models.Signature{}
		// Mock
		getFileInfo = func(name string) (fs.FileInfo, error) {
			return nil, testError
		}

		checkNotExists = func(err error) bool {
			return true
		}

		// Run
		signature, err := OpenSignature(fileName, false)
		// Verify
		require.Equal(t, expectedError, err)
		require.Equal(t, expectedSignature, signature)
	})

	t.Run("should return `emptySignature, UnableToOpenSignatureFileError` when unable to open file", func(t *testing.T) {
		// Setup
		testError := errors.New(errorMessage)
		expectedError := errors.New(constants.UnableToOpenSignatureFileError)
		expectedSignature := models.Signature{}
		// Mock
		getFileInfo = func(name string) (fs.FileInfo, error) {
			fileInfo := fileInfoMock{isDir: false}
			return fileInfo, nil
		}

		open = func(name string) (*os.File, error) {
			return nil, testError
		}

		// Run
		signature, err := OpenSignature(fileName, false)
		// Verify
		require.Equal(t, expectedError, err)
		require.Equal(t, expectedSignature, signature)
	})

	t.Run("should return `emptySignature, UnableToDecodeSignatureFromFileError` when unable to decode Signature from file", func(t *testing.T) {
		// Setup
		file := os.File{}
		decoder := decoderMock{isError: true}
		expectedError := errors.New(constants.UnableToDecodeSignatureFromFileError)
		expectedSignature := models.Signature{}
		// Mock
		getFileInfo = func(name string) (fs.FileInfo, error) {
			fileInfo := fileInfoMock{isDir: false}
			return fileInfo, nil
		}

		open = func(name string) (*os.File, error) {
			return &file, nil
		}

		createNewDecoder = func(file *os.File) Decoder {
			return decoder
		}

		// Run
		signature, err := OpenSignature(fileName, false)
		// Verify
		require.Equal(t, expectedError, err)
		require.Equal(t, expectedSignature, signature)
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

	t.Run("should return `UnableToCreateOutputsFolderError` error when unable to create folder dir", func(t *testing.T) {
		// Setup
		expectedError := errors.New(constants.UnableToCreateOutputsFolderError)
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

func TestWriteStructToFile(t *testing.T) {
	t.Run("should return `nil` when successfully written Signature to output file", func(t *testing.T) {
		// Setup
		file := os.File{}
		encoder := encoderMock{isError: false}
		signature := models.Signature{}
		// Mock
		getFileInfo = func(name string) (fs.FileInfo, error) {
			fileInfo := fileInfoMock{isDir: false}
			return fileInfo, nil
		}

		createFile = func(name string) (*os.File, error) {
			return &file, nil
		}

		createNewEncoder = func(file *os.File) Encoder {
			return encoder
		}

		// Run
		result := WriteStructToFile(signature, fileName)
		// Verify
		require.Equal(t, nil, result)
	})

	t.Run("should return `nil` when successfully created Outputs folder and written Signature to output file", func(t *testing.T) {
		// Setup
		file := os.File{}
		encoder := encoderMock{isError: false}
		signature := models.Signature{}
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

		createNewEncoder = func(file *os.File) Encoder {
			return encoder
		}

		// Run
		result := WriteStructToFile(signature, fileName)
		// Verify
		require.Equal(t, nil, result)
	})

	t.Run("should return `error` when unable to verify if Output dir exists", func(t *testing.T) {
		// Setup
		signature := models.Signature{}
		expectedError := errors.New(constants.UnableToCheckFileFolderExistsError)
		// Mock
		getFileInfo = func(name string) (fs.FileInfo, error) {
			return nil, errors.New(errorMessage)
		}

		checkNotExists = func(err error) bool {
			return false
		}

		// Run
		result := WriteStructToFile(signature, fileName)
		// Verify
		require.Equal(t, expectedError, result)
	})

	t.Run("should return `UnableToCreateFileError` error when unable to create file", func(t *testing.T) {
		// Setup
		file := os.File{}
		signature := models.Signature{}
		expectedError := errors.New(constants.UnableToCreateFileError)
		// Mock
		getFileInfo = func(name string) (fs.FileInfo, error) {
			fileInfo := fileInfoMock{isDir: false}
			return fileInfo, nil
		}

		createFile = func(name string) (*os.File, error) {
			return &file, errors.New(errorMessage)
		}

		// Run
		result := WriteStructToFile(signature, fileName)
		// Verify
		require.Equal(t, expectedError, result)
	})

	t.Run("should return `UnableToWriteToFileError` error when unable to write to file", func(t *testing.T) {
		// Setup
		file := os.File{}
		encoder := encoderMock{isError: true}
		signature := models.Signature{}
		expectedError := errors.New(constants.UnableToWriteToFileError)
		// Mock
		getFileInfo = func(name string) (fs.FileInfo, error) {
			fileInfo := fileInfoMock{isDir: false}
			return fileInfo, nil
		}

		createFile = func(name string) (*os.File, error) {
			return &file, nil
		}

		createNewEncoder = func(file *os.File) Encoder {
			return encoder
		}

		// Run
		result := WriteStructToFile(signature, fileName)
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

	t.Run("should return `UnableToCreateFileError` error when unable to create file", func(t *testing.T) {
		// Setup
		file := os.File{}
		output := []byte(testOutput)
		expectedError := errors.New(constants.UnableToCreateFileError)
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

	t.Run("should return `UnableToWriteToFileError` error when unable to write to file", func(t *testing.T) {
		// Setup
		file := os.File{}
		output := []byte(testOutput)
		expectedError := errors.New(constants.UnableToWriteToFileError)
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
