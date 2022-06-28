package cmd

import (
	"testing"

	"github.com/curtismenmuir/go-file-diff/models"
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
		cmd := ParseCMD()
		// Verify
		require.Equal(t, true, cmd.Verbose)
		require.Equal(t, true, cmd.SignatureMode)
		require.Equal(t, true, cmd.DeltaMode)
		require.Equal(t, file, cmd.OriginalFile)
		require.Equal(t, file, cmd.SignatureFile)
		require.Equal(t, file, cmd.UpdatedFile)
		require.Equal(t, file, cmd.DeltaFile)
	})
}

func TestVerifyCMD(t *testing.T) {
	t.Run("should return true when signature mode set with correct files", func(t *testing.T) {
		// Setup
		cmd := models.CMD{
			Verbose:       false,
			SignatureMode: true,
			DeltaMode:     false,
			OriginalFile:  file,
			SignatureFile: file,
			UpdatedFile:   "",
			DeltaFile:     "",
		}

		// Run
		result := VerifyCMD(cmd)
		// Verify
		require.Equal(t, true, result)
	})

	t.Run("should return true when delta mode set with correct files", func(t *testing.T) {
		// Setup
		cmd := models.CMD{
			Verbose:       false,
			SignatureMode: false,
			DeltaMode:     true,
			OriginalFile:  "",
			SignatureFile: file,
			UpdatedFile:   file,
			DeltaFile:     file,
		}

		// Run
		result := VerifyCMD(cmd)
		// Verify
		require.Equal(t, true, result)
	})

	t.Run("should return true when signature & delta modes set with correct files", func(t *testing.T) {
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

		// Run
		result := VerifyCMD(cmd)
		// Verify
		require.Equal(t, true, result)
	})

	t.Run("should return false when no mode set", func(t *testing.T) {
		// Setup
		cmd := models.CMD{
			Verbose:       false,
			SignatureMode: false,
			DeltaMode:     false,
			OriginalFile:  "",
			SignatureFile: "",
			UpdatedFile:   "",
			DeltaFile:     "",
		}

		// Run
		result := VerifyCMD(cmd)
		// Verify
		require.Equal(t, false, result)
	})

	t.Run("should return false when signature mode set but missing original file", func(t *testing.T) {
		// Setup
		cmd := models.CMD{
			Verbose:       false,
			SignatureMode: true,
			DeltaMode:     false,
			OriginalFile:  "",
			SignatureFile: file,
			UpdatedFile:   "",
			DeltaFile:     "",
		}

		// Run
		result := VerifyCMD(cmd)
		// Verify
		require.Equal(t, false, result)
	})

	t.Run("should return false when signature mode set but missing signature file", func(t *testing.T) {
		// Setup
		cmd := models.CMD{
			Verbose:       false,
			SignatureMode: true,
			DeltaMode:     false,
			OriginalFile:  file,
			SignatureFile: "",
			UpdatedFile:   "",
			DeltaFile:     "",
		}

		// Run
		result := VerifyCMD(cmd)
		// Verify
		require.Equal(t, false, result)
	})

	t.Run("should return false when delta mode set but missing signature file", func(t *testing.T) {
		// Setup
		cmd := models.CMD{
			Verbose:       false,
			SignatureMode: false,
			DeltaMode:     true,
			OriginalFile:  "",
			SignatureFile: "",
			UpdatedFile:   file,
			DeltaFile:     file,
		}

		// Run
		result := VerifyCMD(cmd)
		// Verify
		require.Equal(t, false, result)
	})

	t.Run("should return false when delta mode set but missing update file", func(t *testing.T) {
		// Setup
		cmd := models.CMD{
			Verbose:       false,
			SignatureMode: false,
			DeltaMode:     true,
			OriginalFile:  "",
			SignatureFile: file,
			UpdatedFile:   "",
			DeltaFile:     file,
		}

		// Run
		result := VerifyCMD(cmd)
		// Verify
		require.Equal(t, false, result)
	})

	t.Run("should return false when delta mode set but missing delta file", func(t *testing.T) {
		// Setup
		cmd := models.CMD{
			Verbose:       false,
			SignatureMode: false,
			DeltaMode:     true,
			OriginalFile:  "",
			SignatureFile: file,
			UpdatedFile:   file,
			DeltaFile:     "",
		}

		// Run
		result := VerifyCMD(cmd)
		// Verify
		require.Equal(t, false, result)
	})

	t.Run("should return false when signature & delta modes set but missing update file", func(t *testing.T) {
		// Setup
		cmd := models.CMD{
			Verbose:       false,
			SignatureMode: true,
			DeltaMode:     true,
			OriginalFile:  file,
			SignatureFile: file,
			UpdatedFile:   "",
			DeltaFile:     file,
		}

		// Run
		result := VerifyCMD(cmd)
		// Verify
		require.Equal(t, false, result)
	})

	t.Run("should return false when signature & delta modes set but missing delta file", func(t *testing.T) {
		// Setup
		cmd := models.CMD{
			Verbose:       false,
			SignatureMode: true,
			DeltaMode:     true,
			OriginalFile:  file,
			SignatureFile: file,
			UpdatedFile:   file,
			DeltaFile:     "",
		}

		// Run
		result := VerifyCMD(cmd)
		// Verify
		require.Equal(t, false, result)
	})
}
