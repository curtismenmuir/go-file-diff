package main

import (
	"errors"

	"github.com/curtismenmuir/go-file-diff/cmd"
	"github.com/curtismenmuir/go-file-diff/constants"
	"github.com/curtismenmuir/go-file-diff/files"
	"github.com/curtismenmuir/go-file-diff/models"
	"github.com/curtismenmuir/go-file-diff/sync"
	"github.com/curtismenmuir/go-file-diff/utils"
)

var (
	logger            = utils.Logger
	parseCMD          = cmd.ParseCMD
	verifyCMD         = cmd.VerifyCMD
	openFile          = files.OpenFile
	writeSigToFile    = files.WriteSignatureToFile
	generateSignature = sync.GenerateSignature
)

// getSignature() will generate a Signature of a specified file and write the Signature output to a file
// Function returns `nil` when successful
// Function returns `Original File does not exist` error when Original file cannot be found
// Function returns `Original File is a folder dir` error when found a folder dir instead of Original file
// Function returns `Unable to generate Signature` error when unable to generate file Signature
// Function returns `Unable to create Signature` error when unable to create Signature output file
// Function returns `Unable to write Signature` error when unable to write Signature to output file
func getSignature(cmd models.CMD) error {
	// Read Original file
	reader, err := openFile(cmd.OriginalFile)
	if err != nil {
		// Replace generic `file not exist` error with specific Original File error
		if err.Error() == constants.FileDoesNotExistError {
			return errors.New(constants.OriginalFileDoesNotExistError)
		}

		// Replace generic `file is folder dir` error with specific Original File error
		if err.Error() == constants.SearchingForFileButFoundDirError {
			return errors.New(constants.OriginalFileIsFolderError)
		}

		return err
	}

	// Generate Signature
	signature, err := generateSignature(reader, cmd.Verbose)
	if err != nil {
		return errors.New(constants.UnableToGenerateSignature)
	}

	// Write Signature to file
	err = writeSigToFile(signature, cmd.SignatureFile)
	if err != nil {
		logger(err.Error(), true)
		return err
	}

	return nil
}

// getDelta is a placeholder which returns "not implemented" error
func getDelta(cmd models.CMD) error {
	return errors.New(constants.DeltaNotImplementedError)
}

func main() {
	// Parse CMD flags
	cmd := parseCMD()
	// Verify valid CMD flags provided
	if !verifyCMD(cmd) {
		return
	}

	// Generate Signature
	if cmd.SignatureMode {
		err := getSignature(cmd)
		if err != nil {
			logger(err.Error(), true)
			return
		}
	}

	// Generate Delta
	if cmd.DeltaMode {
		err := getDelta(cmd)
		if err != nil {
			logger(err.Error(), true)
			return
		}
	}
}
