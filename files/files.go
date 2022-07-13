package files

import (
	"bufio"
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/curtismenmuir/go-file-diff/constants"
	"github.com/curtismenmuir/go-file-diff/models"
	"github.com/curtismenmuir/go-file-diff/utils"
)

var (
	open             = os.Open
	getFileInfo      = os.Stat
	checkNotExists   = os.IsNotExist
	mkdir            = os.Mkdir
	createFile       = os.Create
	logger           = utils.Logger
	newWriter        = bufio.NewWriter
	createNewWriter  = createWriter
	createNewEncoder = createEncoder
	newEncoder       = gob.NewEncoder
)

// Encoder interface for mocking gob.NewEncoder
type Encoder interface {
	Encode(e any) error
}

// Writer interface for mocking bufio.NewWriter
type Writer interface {
	io.Writer
	WriteByte(c byte) error
	Flush() error
}

const outputDir string = "./Outputs/"

// createFolder() will attempt to create a folder based on provided folderName prop
// Function will return `nil` when folder is created successfully
// Function will return `unable to create folder` error when unable to create folder dir
func createFolder(folderName string) error {
	if err := mkdir(folderName, os.ModePerm); err != nil {
		return errors.New(constants.UnableToCreateNewFolderError)
	}

	return nil
}

// createEncoder() will init and return a new gob file encoder
// Returned file encoder will satisfy the `Encoder` interface
func createEncoder(file *os.File) Encoder {
	return newEncoder(file)
}

// createWriter() will init and return a new bufio file writer
// Returned file writer will satisfy the `Writer` interface
func createWriter(file *os.File) Writer {
	return newWriter(file)
}

// doesExist() checks if a file/folder exists and returns `true, nil` if specified file/folder is found
// When checking existence of a file, set isFile to true
// When checking existence of a folder dir, set isFile to false
// Function will return `false, nil` when file/folder does not exist
// Function will return `false, error` if an error is thrown
// Function will return `found folder` error if searching for file but found a folder dir
func doesExist(path string, isFile bool) (bool, error) {
	// Attempt to get FileInfo
	fileInfo, err := getFileInfo(path)
	if err != nil {
		// Check if `not exists` error
		if checkNotExists(err) {
			return false, nil
		}

		return false, errors.New(constants.UnableToCheckFileFolderExistsError)
	}

	// If checking file, verify file is not folder dir
	if isFile && fileInfo.IsDir() {
		return false, errors.New(constants.SearchingForFileButFoundDirError)
	}

	return true, nil
}

// OpenFile() will attempt to open a local file and will return a file reader when successful
// Function will catch and return error when unable to access specified file
// Function will return `file does not exist` error when specified file does not exist
func OpenFile(fileName string) (*bufio.Reader, error) {
	// Check if file exists
	exists, err := doesExist(fileName, true)
	if err != nil {
		return nil, err
	} else if !exists {
		return nil, errors.New(constants.FileDoesNotExistError)
	}

	// Open file
	file, err := open(fileName)
	if err != nil {
		return nil, err
	}

	// Return file reader
	return bufio.NewReader(file), nil
}

// verifyOutputDirExists() will check for the existence of an `Outputs/` folder and will create if not exists
// Function will return `nil` when folder already exists
// Function will return `nil` when folder has been created successfully
// Function will return `unable to create Outputs dir` error when folder does not exist and unable to create
// Function will return error when unable to verify if Outputs folder exists
func verifyOutputDirExists() error {
	// Check if `Outputs` folder exists
	exists, err := doesExist(outputDir, false)
	if err != nil {
		return err
	} else if !exists {
		// Create folder if not exists
		err = createFolder(outputDir)
		if err != nil {
			return errors.New(constants.UnableToCreateOutputsFolder)
		}
	}

	return nil
}

// WriteSignatureToFile() will create a Signature file in Outputs folder (based on provided fileName), and encode Signature before writing to file
// Function will return `nil` when file has been created and written to successfully
// Function will return `unable to create Sig file` error when unable to create file
// Function will return `unable to write to Sig file` error when unable to write output to file after creation
// Function will return `error` when unable to verify if Output folder exists
func WriteSignatureToFile(signature []models.Signature, fileName string) error {
	// Verify `Outputs` folder exists
	err := verifyOutputDirExists()
	if err != nil {
		return err
	}

	// Create file
	file, err := createFile(outputDir + fileName)
	if err != nil {
		return errors.New(constants.UnableToCreateSignatureFile)
	}

	defer file.Close()
	// Create encoder
	encoder := createNewEncoder(file)
	// Encode Signature
	err = encoder.Encode(signature)
	if err != nil {
		return errors.New(constants.UnableToWriteToSignatureFile)
	}

	logger(fmt.Sprintf("Signature created: %s%s", outputDir, fileName), true)
	return nil
}

// WriteToFile() will create a file in Outputs folder (based on provided fileName), and write the provided output to the file
// Function will return `nil` when file has been created and written to successfully
// Function will return `unable to create Sig file` error when unable to create file
// Function will return `unable to write to Sig file` error when unable to write output to file after creation
// Function will return error when unable to verify if Output folder exists
func WriteToFile(fileName string, output []byte) error {
	// Verify `Outputs` folder exists
	err := verifyOutputDirExists()
	if err != nil {
		return err
	}

	// Create file
	file, err := createFile(outputDir + fileName)
	if err != nil {
		return errors.New(constants.UnableToCreateSignatureFile)
	}

	defer file.Close()
	fileWriter := createNewWriter(file)
	// Loop over output and write individual bytes
	for index := range output {
		err := fileWriter.WriteByte(output[index])
		if err != nil {
			return errors.New(constants.UnableToWriteToSignatureFile)
		}
	}

	// Flush writer updates to file
	fileWriter.Flush()
	return nil
}
