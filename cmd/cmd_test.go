package cmd

import (
	"testing"

	"github.com/stretchr/testify/require"
)

const file string = "some-file.txt"

func TestParseCMD(t *testing.T) {
	t.Run("should return correct types for CMD args", func(t *testing.T) {
		// Mock
		defineBool = func(name string, value bool, usage string) *bool {
			result := true
			return &result
		}
		defineString = func(name, value, usage string) *string {
			result := file
			return &result
		}
		// Run
		verbose, signatureMode, deltaMode, originalFile, signatureFile, updatedFile, deltaFile := ParseCMD()
		// Verify
		require.Equal(t, true, verbose)
		require.Equal(t, true, signatureMode)
		require.Equal(t, true, deltaMode)
		require.Equal(t, file, originalFile)
		require.Equal(t, file, signatureFile)
		require.Equal(t, file, updatedFile)
		require.Equal(t, file, deltaFile)
	})
}

func TestVerifyCMD(t *testing.T) {
	t.Run("should return true when signature mode set with correct files", func(t *testing.T) {
		// Setup
		verbose := false
		signatureMode := true
		deltaMode := false
		originalFile := file
		signatureFile := file
		updatedFile := ""
		deltaFile := ""
		// Run
		result := VerifyCMD(verbose, signatureMode, deltaMode, originalFile, signatureFile, updatedFile, deltaFile)
		// Verify
		require.Equal(t, true, result)
	})

	t.Run("should return true when delta mode set with correct files", func(t *testing.T) {
		// Setup
		verbose := false
		signatureMode := false
		deltaMode := true
		originalFile := ""
		signatureFile := file
		updatedFile := file
		deltaFile := file
		// Run
		result := VerifyCMD(verbose, signatureMode, deltaMode, originalFile, signatureFile, updatedFile, deltaFile)
		// Verify
		require.Equal(t, true, result)
	})

	t.Run("should return true when signature & delta modes set with correct files", func(t *testing.T) {
		// Setup
		verbose := false
		signatureMode := true
		deltaMode := true
		originalFile := file
		signatureFile := file
		updatedFile := file
		deltaFile := file
		// Run
		result := VerifyCMD(verbose, signatureMode, deltaMode, originalFile, signatureFile, updatedFile, deltaFile)
		// Verify
		require.Equal(t, true, result)
	})

	t.Run("should return false when no mode set", func(t *testing.T) {
		// Setup
		verbose := false
		signatureMode := false
		deltaMode := false
		originalFile := ""
		signatureFile := ""
		updatedFile := ""
		deltaFile := ""
		// Run
		result := VerifyCMD(verbose, signatureMode, deltaMode, originalFile, signatureFile, updatedFile, deltaFile)
		// Verify
		require.Equal(t, false, result)
	})

	t.Run("should return false when signature mode set but missing original file", func(t *testing.T) {
		// Setup
		verbose := false
		signatureMode := true
		deltaMode := false
		originalFile := ""
		signatureFile := file
		updatedFile := ""
		deltaFile := ""
		// Run
		result := VerifyCMD(verbose, signatureMode, deltaMode, originalFile, signatureFile, updatedFile, deltaFile)
		// Verify
		require.Equal(t, false, result)
	})

	t.Run("should return false when signature mode set but missing signature file", func(t *testing.T) {
		// Setup
		verbose := false
		signatureMode := true
		deltaMode := false
		originalFile := file
		signatureFile := ""
		updatedFile := ""
		deltaFile := ""
		// Run
		result := VerifyCMD(verbose, signatureMode, deltaMode, originalFile, signatureFile, updatedFile, deltaFile)
		// Verify
		require.Equal(t, false, result)
	})

	t.Run("should return false when delta mode set but missing signature file", func(t *testing.T) {
		// Setup
		verbose := false
		signatureMode := false
		deltaMode := true
		originalFile := ""
		signatureFile := ""
		updatedFile := file
		deltaFile := file
		// Run
		result := VerifyCMD(verbose, signatureMode, deltaMode, originalFile, signatureFile, updatedFile, deltaFile)
		// Verify
		require.Equal(t, false, result)
	})

	t.Run("should return false when delta mode set but missing update file", func(t *testing.T) {
		// Setup
		verbose := false
		signatureMode := false
		deltaMode := true
		originalFile := ""
		signatureFile := file
		updatedFile := ""
		deltaFile := file
		// Run
		result := VerifyCMD(verbose, signatureMode, deltaMode, originalFile, signatureFile, updatedFile, deltaFile)
		// Verify
		require.Equal(t, false, result)
	})

	t.Run("should return false when delta mode set but missing delta file", func(t *testing.T) {
		// Setup
		verbose := false
		signatureMode := false
		deltaMode := true
		originalFile := ""
		signatureFile := file
		updatedFile := file
		deltaFile := ""
		// Run
		result := VerifyCMD(verbose, signatureMode, deltaMode, originalFile, signatureFile, updatedFile, deltaFile)
		// Verify
		require.Equal(t, false, result)
	})

	t.Run("should return false when signature & delta modes set but missing update file", func(t *testing.T) {
		// Setup
		verbose := false
		signatureMode := true
		deltaMode := true
		originalFile := file
		signatureFile := file
		updatedFile := ""
		deltaFile := file
		// Run
		result := VerifyCMD(verbose, signatureMode, deltaMode, originalFile, signatureFile, updatedFile, deltaFile)
		// Verify
		require.Equal(t, false, result)
	})

	t.Run("should return false when signature & delta modes set but missing delta file", func(t *testing.T) {
		// Setup
		verbose := false
		signatureMode := true
		deltaMode := true
		originalFile := file
		signatureFile := file
		updatedFile := file
		deltaFile := ""
		// Run
		result := VerifyCMD(verbose, signatureMode, deltaMode, originalFile, signatureFile, updatedFile, deltaFile)
		// Verify
		require.Equal(t, false, result)
	})
}
