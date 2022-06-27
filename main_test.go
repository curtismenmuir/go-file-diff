package main

import (
	"errors"
	"fmt"
	"testing"

	"github.com/curtismenmuir/go-file-diff/constants"
	"github.com/stretchr/testify/require"
)

func TestGetSignature(t *testing.T) {
	t.Run("should throw `not implemented` error", func(t *testing.T) {
		// Setup
		expectedError := errors.New(constants.SignatureNotImplementedError)
		// Run
		err := getSignature(false, "some-file.txt", "some-file.txt")
		// Verify
		require.Equal(t, expectedError, err)
	})
}

func TestGetDelta(t *testing.T) {
	t.Run("should throw `not implemented` error", func(t *testing.T) {
		// Setup
		expectedError := errors.New(constants.DeltaNotImplementedError)
		// Run
		err := getDelta(false, "some-file.txt", "some-file.txt", "some-file.txt")
		// Verify
		require.Equal(t, expectedError, err)
	})
}

func TestMain(t *testing.T) {
	t.Run("should throw error when generating signature", func(t *testing.T) {
		// Setup
		logged := false
		loggedMessage := ""
		expectedError := fmt.Sprintf("Error: %s", constants.SignatureNotImplementedError)
		// Mock
		logger = func(message string, verbose bool) {
			logged = true
			loggedMessage = message
		}
		parseCMD = func() (bool, bool, bool, string, string, string, string) {
			return false, true, false, "some-file.txt", "some-file.txt", "some-file.txt", "some-file.txt"
		}
		verifyCMD = func(verbose bool, signatureMode bool, deltaMode bool, originalFile string, signatureFile string, updatedFile string, deltaFile string) bool {
			return true
		}
		// Run
		main()
		// Verify
		require.Equal(t, true, logged)
		require.Equal(t, expectedError, loggedMessage)
	})

	t.Run("should throw error when generating delta", func(t *testing.T) {
		// Setup
		logged := false
		loggedMessage := ""
		expectedError := fmt.Sprintf("Error: %s", constants.DeltaNotImplementedError)
		// Mock
		logger = func(message string, verbose bool) {
			logged = true
			loggedMessage = message
		}
		parseCMD = func() (bool, bool, bool, string, string, string, string) {
			return false, false, true, "some-file.txt", "some-file.txt", "some-file.txt", "some-file.txt"
		}
		verifyCMD = func(verbose bool, signatureMode bool, deltaMode bool, originalFile string, signatureFile string, updatedFile string, deltaFile string) bool {
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
		logged := false
		// Mock
		logger = func(message string, verbose bool) {
			logged = true
		}
		parseCMD = func() (bool, bool, bool, string, string, string, string) {
			return false, true, true, "some-file.txt", "some-file.txt", "some-file.txt", "some-file.txt"
		}
		verifyCMD = func(verbose bool, signatureMode bool, deltaMode bool, originalFile string, signatureFile string, updatedFile string, deltaFile string) bool {
			return false
		}
		// Run
		main()
		// Verify
		require.Equal(t, false, logged)
	})
}
