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
	writeToFile       = files.WriteToFile
	generateSignature = sync.GenerateSignature
)

// getSignature is a placeholder which returns "not implemented" error when provided a valid Original file (CMD flags)
// Function will catch and return any errors if they occur
// Function returns `Original File does not exist` error when Original file cannot be found
// Function returns `Original File is a folder dir` error when found a folder dir instead of Original file
func getSignature(cmd models.CMD) error {
	// Read Original file
	reader, err := openFile(cmd.OriginalFile)
	if err != nil {
		// Replace `file not exist` error with specific Original File error
		if err.Error() == constants.FileDoesNotExistError {
			return errors.New(constants.OriginalFileDoesNotExistError)
		}

		// Replace `file is folder dir` error with specific Original File error
		if err.Error() == constants.SearchingForFileButFoundDirError {
			return errors.New(constants.OriginalFileIsFolderError)
		}

		return err
	}

	// Generate Signature
	_ = generateSignature(reader, cmd.Verbose)
	signature := []byte("Testing `write to file` for now.....\n!\"£$%^&*(){}:@~>?<,./;'#[]\n\nFile signature coming soon!\n")

	// Write Signature to file
	err = writeToFile(cmd.SignatureFile, signature)
	if err != nil {
		logger(err.Error(), true)
		return err
	}

	return errors.New(constants.SignatureNotImplementedError)
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
