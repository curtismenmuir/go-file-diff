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

var (
	file          string             = "some-file.txt"
	errorMessage  string             = "Some Error"
	testSignature []models.Signature = []models.Signature{{Weak: 123, Strong: "some-hex"}}
)

func TestGetSignature(t *testing.T) {
	t.Run("should return `Signature, nil` when Signature generated successfully", func(t *testing.T) {
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

		// Mock
		openFile = func(fileName string) (*bufio.Reader, error) {
			file := os.File{}
			return bufio.NewReader(&file), nil
		}

		generateSignature = func(reader sync.Reader, verbose bool) ([]models.Signature, error) {
			return testSignature, nil
		}

		writeSigToFile = func(signature []models.Signature, fileName string) error {
			return nil
		}

		// Run
		signature, err := getSignature(cmd)
		// Verify
		require.Equal(t, nil, err)
		require.Equal(t, testSignature, signature)
	})

	t.Run("should return `EmptySignature, OriginalFileNotExistError` when Original file cannot be found", func(t *testing.T) {
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

		// Run
		signature, err := getSignature(cmd)
		// Verify
		require.Equal(t, expectedError, err)
		require.Equal(t, []models.Signature{}, signature)
	})

	t.Run("should return `EmptySignature, OriginalFileIsFolderError` when user provides a folder dir instead of Original file", func(t *testing.T) {
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

		// Run
		signature, err := getSignature(cmd)
		// Verify
		require.Equal(t, expectedError, err)
		require.Equal(t, []models.Signature{}, signature)
	})

	t.Run("should return `EmptySignature, Error` when unable to open Original file", func(t *testing.T) {
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

		// Run
		signature, err := getSignature(cmd)
		// Verify
		require.Equal(t, expectedError, err)
		require.Equal(t, []models.Signature{}, signature)
	})

	t.Run("should return `EmptySignature, UnableToGenerateSignatureError` when fails to generate Signature of Original file", func(t *testing.T) {
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

		expectedError := errors.New(constants.UnableToGenerateSignatureError)
		// Mock
		openFile = func(fileName string) (*bufio.Reader, error) {
			file := os.File{}
			return bufio.NewReader(&file), nil
		}

		generateSignature = func(reader sync.Reader, verbose bool) ([]models.Signature, error) {
			return nil, expectedError
		}

		// Run
		signature, err := getSignature(cmd)
		// Verify
		require.Equal(t, expectedError, err)
		require.Equal(t, []models.Signature{}, signature)
	})

	t.Run("should return `EmptySignature, UnableToWriteToSignatureFileError` when unable to write to Signature file", func(t *testing.T) {
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

		expectedError := errors.New(constants.UnableToWriteToSignatureFileError)
		// Mock
		openFile = func(fileName string) (*bufio.Reader, error) {
			file := os.File{}
			return bufio.NewReader(&file), nil
		}

		generateSignature = func(reader sync.Reader, verbose bool) ([]models.Signature, error) {
			return nil, nil
		}

		writeSigToFile = func(signature []models.Signature, fileName string) error {
			return errors.New(errorMessage)
		}

		// Run
		signature, err := getSignature(cmd)
		// Verify
		require.Equal(t, expectedError, err)
		require.Equal(t, []models.Signature{}, signature)
	})
}

func TestGetDelta(t *testing.T) {
	t.Run("should return `DeltaNotImplementedError` when successfully opened Updated file and generated Delta", func(t *testing.T) {
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

		expectedError := errors.New(constants.DeltaNotImplementedError)
		// Mock
		openFile = func(fileName string) (*bufio.Reader, error) {
			file := os.File{}
			return bufio.NewReader(&file), nil
		}

		generateDelta = func(reader sync.Reader, signature []models.Signature, verbose bool) error {
			return nil
		}

		// Run
		err := getDelta(cmd, testSignature)
		// Verify
		require.Equal(t, expectedError, err)
	})

	t.Run("should return `UpdatedFileDoesNotExistError` when unable to find Updated file", func(t *testing.T) {
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

		expectedError := errors.New(constants.UpdatedFileDoesNotExistError)
		// Mock
		openFile = func(fileName string) (*bufio.Reader, error) {
			return nil, errors.New(constants.FileDoesNotExistError)
		}

		// Run
		err := getDelta(cmd, testSignature)
		// Verify
		require.Equal(t, expectedError, err)
	})

	t.Run("should return `UpdatedFileIsFolderError` when found Updated file but it is a folder dir", func(t *testing.T) {
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

		expectedError := errors.New(constants.UpdatedFileIsFolderError)
		// Mock
		openFile = func(fileName string) (*bufio.Reader, error) {
			return nil, errors.New(constants.SearchingForFileButFoundDirError)
		}

		// Run
		err := getDelta(cmd, testSignature)
		// Verify
		require.Equal(t, expectedError, err)
	})

	t.Run("should return `error` when unable to open Updated file", func(t *testing.T) {
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

		expectedError := errors.New(errorMessage)
		// Mock
		openFile = func(fileName string) (*bufio.Reader, error) {
			return nil, expectedError
		}

		// Run
		err := getDelta(cmd, testSignature)
		// Verify
		require.Equal(t, expectedError, err)
	})

	t.Run("should return `UnableToGenerateDeltaError` when unable to generate Delta", func(t *testing.T) {
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

		expectedError := errors.New(constants.UnableToGenerateDeltaError)
		// Mock
		openFile = func(fileName string) (*bufio.Reader, error) {
			file := os.File{}
			return bufio.NewReader(&file), nil
		}

		generateDelta = func(reader sync.Reader, signature []models.Signature, verbose bool) error {
			return expectedError
		}

		// Run
		err := getDelta(cmd, testSignature)
		// Verify
		require.Equal(t, expectedError, err)
	})
}

func TestMain(t *testing.T) {
	t.Run("should not throw error when successfully generated Signature", func(t *testing.T) {
		// Setup
		cmd := models.CMD{
			Verbose:       false,
			SignatureMode: true,
			DeltaMode:     false,
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
			return true
		}

		openFile = func(fileName string) (*bufio.Reader, error) {
			file := os.File{}
			return bufio.NewReader(&file), nil
		}

		writeSigToFile = func(signature []models.Signature, fileName string) error {
			return nil
		}

		// Run
		main()
		// Verify
		require.Equal(t, false, logged)
	})

	t.Run("should throw `UnableToWriteToSignatureFileError` when unable to generate Signature", func(t *testing.T) {
		// Setup
		cmd := models.CMD{
			Verbose:       false,
			SignatureMode: true,
			DeltaMode:     false,
			OriginalFile:  file,
			SignatureFile: file,
			UpdatedFile:   file,
			DeltaFile:     file,
		}

		logged := false
		loggedMessage := ""
		expectedError := constants.UnableToWriteToSignatureFileError
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

		writeSigToFile = func(signature []models.Signature, fileName string) error {
			return errors.New(errorMessage)
		}

		// Run
		main()
		// Verify
		require.Equal(t, true, logged)
		require.Equal(t, expectedError, loggedMessage)
	})

	t.Run("should throw `DeltaNotImplementedError` when attempting to generate Delta after opening Signature from file", func(t *testing.T) {
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

		openSignature = func(fileName string, verbose bool) ([]models.Signature, error) {
			return testSignature, nil
		}

		openFile = func(fileName string) (*bufio.Reader, error) {
			file := os.File{}
			return bufio.NewReader(&file), nil
		}

		generateDelta = func(reader sync.Reader, signature []models.Signature, verbose bool) error {
			return nil
		}

		// Run
		main()
		// Verify
		require.Equal(t, true, logged)
		require.Equal(t, expectedError, loggedMessage)
	})

	t.Run("should throw `DeltaNotImplementedError` when attempting to generate Delta after generating Signature", func(t *testing.T) {
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

		openFile = func(fileName string) (*bufio.Reader, error) {
			file := os.File{}
			return bufio.NewReader(&file), nil
		}

		writeSigToFile = func(signature []models.Signature, fileName string) error {
			return nil
		}

		generateDelta = func(reader sync.Reader, signature []models.Signature, verbose bool) error {
			return nil
		}

		// Run
		main()
		// Verify
		require.Equal(t, true, logged)
		require.Equal(t, expectedError, loggedMessage)
	})

	t.Run("should throw error when generating Delta and unable to open Signature file", func(t *testing.T) {
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
		expectedError := constants.UnableToOpenSignatureFileError
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

		openSignature = func(fileName string, verbose bool) ([]models.Signature, error) {
			return nil, errors.New(expectedError)
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
