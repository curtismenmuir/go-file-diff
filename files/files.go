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
	newDecoder       = gob.NewDecoder
	createNewDecoder = createDecoder
)

// Encoder interface for mocking gob.NewEncoder.
type Encoder interface {
	Encode(e any) error
}

// Decoder interface for mocking gob.NewDecoder.
type Decoder interface {
	Decode(e any) error
}

// Writer interface for mocking bufio.NewWriter.
type Writer interface {
	io.Writer
	WriteByte(c byte) error
	Flush() error
}

const outputDir string = "./Outputs/"

// createFolder() will attempt to create a folder based on provided folderName prop.
// Function will return `nil` when folder is created successfully.
// Function will return `unable to create folder` error when unable to create folder dir.
func createFolder(folderName string) error {
	if err := mkdir(folderName, os.ModePerm); err != nil {
		return errors.New(constants.UnableToCreateNewFolderError)
	}

	return nil
}

// createDecoder() will init and return a new gob file decoder.
// Returned file decoder will satisfy the `Decoder` interface.
func createDecoder(file *os.File) Decoder {
	return newDecoder(file)
}

// createEncoder() will init and return a new gob file encoder.
// Returned file encoder will satisfy the `Encoder` interface.
func createEncoder(file *os.File) Encoder {
	return newEncoder(file)
}

// createWriter() will init and return a new bufio file writer.
// Returned file writer will satisfy the `Writer` interface.
func createWriter(file *os.File) Writer {
	return newWriter(file)
}

// doesExist() checks if a file/folder exists and returns `true, nil` if specified file/folder is found.
// When checking existence of a file, set isFile to true.
// When checking existence of a folder dir, set isFile to false.
// Function will return `false, nil` when file/folder does not exist.
// Function will return `false, error` if an error is thrown.
// Function will return `found folder` error if searching for file but found a folder dir.
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

// OpenDelta() will attempt to open a local file and decode a Delta from it.
// Note: this will be used for the `patch` process.
// Function will return `Delta, nil` when successfully retrieve Delta from file.
// Function will return `emptyDelta, error` when unable to check existence of Delta file.
// Function will return `emptyDelta, DeltaFileDoesNotExistError` when Delta file not found.
// Function will return `emptyDelta, UnableToOpenDeltaFileError` when unable to open Delta file.
// Function will return `emptyDelta, UnableToDecodeDeltaFromFileError` when unable to decode Delta from file (EG invalid file).
func OpenDelta(fileName string, verbose bool) (models.Delta, error) {
	delta := models.Delta{}
	// Check if Delta file exists
	exists, err := doesExist(fileName, true)
	if err != nil {
		return delta, err
	} else if !exists {
		return delta, errors.New(constants.DeltaFileDoesNotExistError)
	}

	// Open Delta file
	file, err := open(fileName)
	if err != nil {
		return delta, errors.New(constants.UnableToOpenDeltaFileError)
	}

	defer file.Close()
	// Create new file decoder
	decoder := createNewDecoder(file)
	// Decode file to Delta struct
	err = decoder.Decode(&delta)
	if err != nil {
		return delta, errors.New(constants.UnableToDecodeDeltaFromFileError)
	}

	logger(fmt.Sprintf("File Delta: %+v\n", delta), verbose)
	return delta, nil
}

// OpenFile() will attempt to open a local file and will return a file reader when successful.
// Function will catch and return error when unable to access specified file.
// Function will return `file does not exist` error when specified file does not exist.
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

// OpenSignature() will attempt to open a local file and decode a Signature from the file.
// Function will return `Signature, nil` when successfully retrieve a Signature from file.
// Function will return `emptySignature, error` when unable to check existence of Signature file.
// Function will return `emptySignature, SignatureFileDoesNotExistError` when Signature file not found.
// Function will return `emptySignature, UnableToOpenSignatureFileError` when unable to open Signature file.
// Function will return `emptySignature, UnableToDecodeSignatureFromFileError` when unable to decode Signature from file (EG invalid signature file).
func OpenSignature(fileName string, verbose bool) (models.Signature, error) {
	signature := models.Signature{}
	// Check if Signature file exists
	exists, err := doesExist(fileName, true)
	if err != nil {
		return signature, err
	} else if !exists {
		return signature, errors.New(constants.SignatureFileDoesNotExistError)
	}

	// Open Signature file
	file, err := open(fileName)
	if err != nil {
		return signature, errors.New(constants.UnableToOpenSignatureFileError)
	}

	defer file.Close()
	// Create new file decoder
	decoder := createNewDecoder(file)
	// Decode file to Signature struct
	err = decoder.Decode(&signature)
	if err != nil {
		return signature, errors.New(constants.UnableToDecodeSignatureFromFileError)
	}

	logger(fmt.Sprintf("File Signature: %+v\n", signature), verbose)
	return signature, nil
}

// verifyOutputDirExists() will check for the existence of an `Outputs/` folder and will create if not exists.
// Function will return `nil` when folder already exists.
// Function will return `nil` when folder has been created successfully.
// Function will return `UnableToCreateOutputsFolderError` error when folder does not exist and unable to create.
// Function will return `error` when unable to verify if Outputs folder exists.
func verifyOutputDirExists() error {
	// Check if `Outputs` folder exists
	exists, err := doesExist(outputDir, false)
	if err != nil {
		return err
	} else if !exists {
		// Create folder if not exists
		err = createFolder(outputDir)
		if err != nil {
			return errors.New(constants.UnableToCreateOutputsFolderError)
		}
	}

	return nil
}

// WriteStructToFile() will create a file in Outputs folder (based on provided fileName), and encode provided struct before writing to file.
// Function will return `nil` when file has been created and written to successfully.
// Function will return `UnableToCreateFileError` error when unable to create file.
// Function will return `UnableToWriteToFileError` error when unable to write output to file after creation.
// Function will return `error` when unable to verify if Output folder exists.
func WriteStructToFile(model any, fileName string) error {
	// Verify `Outputs` folder exists
	err := verifyOutputDirExists()
	if err != nil {
		return err
	}

	// Create file
	file, err := createFile(outputDir + fileName)
	if err != nil {
		return errors.New(constants.UnableToCreateFileError)
	}

	defer file.Close()
	// Create encoder
	encoder := createNewEncoder(file)
	// Encode struct
	err = encoder.Encode(model)
	if err != nil {
		return errors.New(constants.UnableToWriteToFileError)
	}

	logger(fmt.Sprintf("%s created: %s%s\n", fileName, outputDir, fileName), true)
	return nil
}

// WriteToFile() will create a file in Outputs folder (based on provided fileName), and write the provided output to the file.
// Note: this will be used for the `patch` process.
// Function will return `nil` when file has been created and written to successfully.
// Function will return `UnableToCreateFileError` error when unable to create file.
// Function will return `UnableToWriteToFileError` error when unable to write output to file after creation.
// Function will return `error` when unable to verify if Output folder exists.
func WriteToFile(fileName string, output []byte) error {
	// Verify `Outputs` folder exists
	err := verifyOutputDirExists()
	if err != nil {
		return err
	}

	// Create file
	file, err := createFile(outputDir + fileName)
	if err != nil {
		return errors.New(constants.UnableToCreateFileError)
	}

	defer file.Close()
	fileWriter := createNewWriter(file)
	// Loop over output and write individual bytes
	for index := range output {
		err := fileWriter.WriteByte(output[index])
		if err != nil {
			return errors.New(constants.UnableToWriteToFileError)
		}
	}

	// Flush writer updates to file
	fileWriter.Flush()
	return nil
}
