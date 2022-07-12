package main

import (
	"bufio"
	"errors"
	"os"
	"testing"

	"github.com/curtismenmuir/go-file-diff/constants"
	"github.com/curtismenmuir/go-file-diff/models"
	"github.com/curtismenmuir/go-file-diff/sync"
	"github.com/stretchr/testify/require"
)

const (
	file         string = "some-file.txt"
	errorMessage string = "Some Error"
)

func TestGetSignature(t *testing.T) {
	t.Run("should throw `not implemented` error when Original file exists and write to Signature file succeeds", func(t *testing.T) {
		// Setup
		cmd := models.CMD{
			Verbose:       false,
			SignatureMode: true,
			DeltaMode:     true,
			OriginalFile:  file,
			SignatureFile: file,
			UpdatedFile:   file,
			DeltaFile:     file,
		}

		expectedError := errors.New(constants.SignatureNotImplementedError)
		// Mock
		openFile = func(fileName string) (*bufio.Reader, error) {
			file := os.File{}
			return bufio.NewReader(&file), nil
		}

		generateSignature = func(reader sync.Reader, verbose bool) ([]models.Signature, error) {
			return nil, nil
		}

		writeToFile = func(fileName string, output []byte) error {
			return nil
		}

		// Run
		err := getSignature(cmd)
		// Verify
		require.Equal(t, expectedError, err)
	})

	t.Run("should throw `Original File not exist` error when Original file cannot be found", func(t *testing.T) {
		// Setup
		cmd := models.CMD{
			Verbose:       false,
			SignatureMode: true,
			DeltaMode:     true,
			OriginalFile:  file,
			SignatureFile: file,
			UpdatedFile:   file,
			DeltaFile:     file,
		}

		expectedError := errors.New(constants.OriginalFileDoesNotExistError)
		// Mock
		openFile = func(fileName string) (*bufio.Reader, error) {
			return nil, errors.New(constants.FileDoesNotExistError)
		}

		writeToFile = func(fileName string, output []byte) error {
			return nil
		}

		// Run
		err := getSignature(cmd)
		// Verify
		require.Equal(t, expectedError, err)
	})

	t.Run("should throw `Original File is folder dir` error when user provides a folder dir instead of Original file", func(t *testing.T) {
		// Setup
		cmd := models.CMD{
			Verbose:       false,
			SignatureMode: true,
			DeltaMode:     true,
			OriginalFile:  file,
			SignatureFile: file,
			UpdatedFile:   file,
			DeltaFile:     file,
		}

		expectedError := errors.New(constants.OriginalFileIsFolderError)
		// Mock
		openFile = func(fileName string) (*bufio.Reader, error) {
			return nil, errors.New(constants.SearchingForFileButFoundDirError)
		}

		writeToFile = func(fileName string, output []byte) error {
			return nil
		}

		// Run
		err := getSignature(cmd)
		// Verify
		require.Equal(t, expectedError, err)
	})

	t.Run("should throw error when unable to open Original file", func(t *testing.T) {
		// Setup
		cmd := models.CMD{
			Verbose:       false,
			SignatureMode: true,
			DeltaMode:     true,
			OriginalFile:  file,
			SignatureFile: file,
			UpdatedFile:   file,
			DeltaFile:     file,
		}

		expectedError := errors.New(errorMessage)
		// Mock
		openFile = func(fileName string) (*bufio.Reader, error) {
			return nil, errors.New(errorMessage)
		}

		writeToFile = func(fileName string, output []byte) error {
			return nil
		}

		// Run
		err := getSignature(cmd)
		// Verify
		require.Equal(t, expectedError, err)
	})

	t.Run("should throw `Unable to generate Signature` error when fails to generate Signature of Original file", func(t *testing.T) {
		// Setup
		cmd := models.CMD{
			Verbose:       false,
			SignatureMode: true,
			DeltaMode:     true,
			OriginalFile:  file,
			SignatureFile: file,
			UpdatedFile:   file,
			DeltaFile:     file,
		}

		expectedError := errors.New(constants.UnableToGenerateSignature)
		// Mock
		openFile = func(fileName string) (*bufio.Reader, error) {
			file := os.File{}
			return bufio.NewReader(&file), nil
		}

		generateSignature = func(reader sync.Reader, verbose bool) ([]models.Signature, error) {
			return nil, expectedError
		}

		writeToFile = func(fileName string, output []byte) error {
			return nil
		}

		// Run
		err := getSignature(cmd)
		// Verify
		require.Equal(t, expectedError, err)
	})

	t.Run("should throw error when unable to write to Signature file", func(t *testing.T) {
		// Setup
		cmd := models.CMD{
			Verbose:       false,
			SignatureMode: true,
			DeltaMode:     true,
			OriginalFile:  file,
			SignatureFile: file,
			UpdatedFile:   file,
			DeltaFile:     file,
		}

		expectedError := errors.New(errorMessage)
		// Mock
		openFile = func(fileName string) (*bufio.Reader, error) {
			file := os.File{}
			return bufio.NewReader(&file), nil
		}

		generateSignature = func(reader sync.Reader, verbose bool) ([]models.Signature, error) {
			return nil, nil
		}

		writeToFile = func(fileName string, output []byte) error {
			return expectedError
		}

		// Run
		err := getSignature(cmd)
		// Verify
		require.Equal(t, expectedError, err)
	})
}

func TestGetDelta(t *testing.T) {
	t.Run("should throw `not implemented` error", func(t *testing.T) {
		// Setup
		cmd := models.CMD{
			Verbose:       false,
			SignatureMode: true,
			DeltaMode:     true,
			OriginalFile:  file,
			SignatureFile: file,
			UpdatedFile:   file,
			DeltaFile:     file,
		}

		expectedError := errors.New(constants.DeltaNotImplementedError)
		// Run
		err := getDelta(cmd)
		// Verify
		require.Equal(t, expectedError, err)
	})
}

func TestMain(t *testing.T) {
	t.Run("should throw error when generating signature", func(t *testing.T) {
		// Setup
		cmd := models.CMD{
			Verbose:       false,
			SignatureMode: true,
			DeltaMode:     true,
			OriginalFile:  file,
			SignatureFile: file,
			UpdatedFile:   file,
			DeltaFile:     file,
		}

		logged := false
		loggedMessage := ""
		expectedError := constants.SignatureNotImplementedError
		// Mock
		logger = func(message string, verbose bool) {
			logged = true
			loggedMessage = message
		}

		parseCMD = func() models.CMD {
			return cmd
		}

		verifyCMD = func(cmd models.CMD) bool {
			return true
		}

		openFile = func(fileName string) (*bufio.Reader, error) {
			file := os.File{}
			return bufio.NewReader(&file), nil
		}

		writeToFile = func(fileName string, output []byte) error {
			return nil
		}

		// Run
		main()
		// Verify
		require.Equal(t, true, logged)
		require.Equal(t, expectedError, loggedMessage)
	})

	t.Run("should throw error when generating delta", func(t *testing.T) {
		// Setup
		cmd := models.CMD{
			Verbose:       false,
			SignatureMode: false,
			DeltaMode:     true,
			OriginalFile:  file,
			SignatureFile: file,
			UpdatedFile:   file,
			DeltaFile:     file,
		}

		logged := false
		loggedMessage := ""
		expectedError := constants.DeltaNotImplementedError
		// Mock
		logger = func(message string, verbose bool) {
			logged = true
			loggedMessage = message
		}

		parseCMD = func() models.CMD {
			return cmd
		}

		verifyCMD = func(cmd models.CMD) bool {
			return true
		}

		// Run
		main()
		// Verify
		require.Equal(t, true, logged)
		require.Equal(t, expectedError, loggedMessage)
	})

	t.Run("should catch error with CMD args", func(t *testing.T) {
		// Setup
		cmd := models.CMD{
			Verbose:       false,
			SignatureMode: true,
			DeltaMode:     true,
			OriginalFile:  file,
			SignatureFile: file,
			UpdatedFile:   file,
			DeltaFile:     file,
		}

		logged := false
		// Mock
		logger = func(message string, verbose bool) {
			logged = true
		}

		parseCMD = func() models.CMD {
			return cmd
		}

		verifyCMD = func(cmd models.CMD) bool {
			return false
		}

		// Run
		main()
		// Verify
		require.Equal(t, false, logged)
	})
}
