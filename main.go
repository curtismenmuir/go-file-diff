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
	writeStructToFile = files.WriteStructToFile
	generateSignature = sync.GenerateSignature
	openSignature     = files.OpenSignature
	generateDelta     = sync.GenerateDelta
)

// getSignature() will generate a Signature of a specified file and write the Signature output to a file.
// Function returns `Signature, nil` when successful.
// Function returns `EmptySignature, OriginalFileNotExistError` when Original file cannot be found.
// Function returns `EmptySignature, OriginalFileIsFolderError` when found a folder dir instead of Original file.
// Function returns `EmptySignature, UnableToGenerateSignatureError` when unable to generate file Signature.
// Function returns `EmptySignature, UnableToWriteToSignatureFileError` when unable to write Signature to output file.
func getSignature(cmd models.CMD) (models.Signature, error) {
	// Create FileReader for Original file
	reader, err := openFile(cmd.OriginalFile)
	if err != nil {
		// Replace generic `file not exist` error with specific Original File error
		if err.Error() == constants.FileDoesNotExistError {
			return models.Signature{}, errors.New(constants.OriginalFileDoesNotExistError)
		}

		// Replace generic `file is folder dir` error with specific Original File error
		if err.Error() == constants.SearchingForFileButFoundDirError {
			return models.Signature{}, errors.New(constants.OriginalFileIsFolderError)
		}

		return models.Signature{}, err
	}

	// Generate Signature
	signature, err := generateSignature(reader, cmd.Verbose)
	if err != nil {
		return models.Signature{}, errors.New(constants.UnableToGenerateSignatureError)
	}

	// Write Signature to file
	err = writeStructToFile(signature, cmd.SignatureFile)
	if err != nil {
		// Replace generic `UnableToCreateFileError` error with specific Signature File error
		if err.Error() == constants.UnableToCreateFileError {
			return models.Signature{}, errors.New(constants.UnableToCreateSignatureFileError)
		}

		return models.Signature{}, errors.New(constants.UnableToWriteToSignatureFileError)
	}

	return signature, nil
}

// getDelta() will attempt to generate a Delta changeset for syncing 2 files.
// Delta changeset can be applied to the Original file to sync latest updates.
// Delta generation will use a Signature of the original file to compare against Updated file.
// Function returns `delta, nil` when successful.
// Function returns `emptyDelta, UpdatedFileDoesNotExistError` when unable to find Updated file.
// Function returns `emptyDelta, UpdatedFileIsFolderError` when found a folder dir instead of Updated file.
// Function returns `emptyDelta, UpdatedFileHasNoChangesError` when Delta generation finds no changes in Updated file.
// Function returns `emptyDelta, UnableToGenerateDeltaError` when unable to generate Delta.
// Function returns `emptyDelta, UnableToCreateDeltaFileError` when unable to create Delta file.
// Function returns `emptyDelta, UnableToWriteToDeltaFileError` when unable to write to Delta file.
func getDelta(cmd models.CMD, signature models.Signature) (models.Delta, error) {
	// Create FileReader for Updated file
	reader, err := openFile(cmd.UpdatedFile)
	if err != nil {
		// Replace generic `file not exist` error with specific Updated File error
		if err.Error() == constants.FileDoesNotExistError {
			return models.Delta{}, errors.New(constants.UpdatedFileDoesNotExistError)
		}

		// Replace generic `file is folder dir` error with specific Updated File error
		if err.Error() == constants.SearchingForFileButFoundDirError {
			return models.Delta{}, errors.New(constants.UpdatedFileIsFolderError)
		}

		return models.Delta{}, err
	}

	// Generate Delta
	delta, err := generateDelta(reader, signature, cmd.Verbose)
	if err != nil {
		// Return err when no changes detected in Updated file
		if err.Error() == constants.UpdatedFileHasNoChangesError {
			return models.Delta{}, err
		}

		// Return generic unable to generate Delta error
		return models.Delta{}, errors.New(constants.UnableToGenerateDeltaError)
	}

	// Write Delta to file
	err = writeStructToFile(delta, cmd.DeltaFile)
	if err != nil {
		// Replace generic `UnableToCreateFileError` error with specific Delta File error
		if err.Error() == constants.UnableToCreateFileError {
			return models.Delta{}, errors.New(constants.UnableToCreateDeltaFileError)
		}

		return models.Delta{}, errors.New(constants.UnableToWriteToDeltaFileError)
	}

	return delta, nil
}

func main() {
	// Parse CMD flags
	cmd := parseCMD()
	// Verify valid CMD flags provided
	if !verifyCMD(cmd) {
		return
	}

	var signature models.Signature
	var err error

	if cmd.SignatureMode {
		// Generate Signature
		signature, err = getSignature(cmd)
		if err != nil {
			logger(err.Error(), true)
			return
		}
	}

	if cmd.DeltaMode {
		// Get signature from file when running delta mode only
		if !cmd.SignatureMode {
			signature, err = openSignature(cmd.SignatureFile, cmd.Verbose)
			if err != nil {
				logger(err.Error(), true)
				return
			}
		}

		// Generate Delta
		_, err = getDelta(cmd, signature)
		if err != nil {
			logger(err.Error(), true)
			return
		}
	}
}
