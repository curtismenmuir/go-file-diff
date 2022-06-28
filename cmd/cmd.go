package cmd

import (
	"flag"
	"fmt"

	"github.com/curtismenmuir/go-file-diff/constants"
	"github.com/curtismenmuir/go-file-diff/models"
	"github.com/curtismenmuir/go-file-diff/utils"
)

var (
	logger       = utils.Logger
	defineBool   = flag.Bool
	defineString = flag.String
)

func ParseCMD() models.CMD {
	// Define CMD flags
	verbose := defineBool("v", false, "Enable extended logging")
	signatureMode := defineBool("signatureMode", false, "Enable Signature mode")
	deltaMode := defineBool("deltaMode", false, "Enable Delta mode")
	originalFile := defineString("original", "", "Original file")
	signatureFile := defineString("signature", "", "Signature file")
	updatedFile := defineString("updated", "", "Updated file")
	deltaFile := defineString("delta", "", "Delta file")

	// Parse CMD flags
	flag.Parse()

	// Format CMD flags
	cmd := models.CMD{
		Verbose:       *verbose,
		SignatureMode: *signatureMode,
		DeltaMode:     *deltaMode,
		OriginalFile:  *originalFile,
		SignatureFile: *signatureFile,
		UpdatedFile:   *updatedFile,
		DeltaFile:     *deltaFile,
	}

	logger(fmt.Sprintf("CMD: %+v", cmd), *verbose)
	return cmd
}

func VerifyCMD(cmd models.CMD) bool {
	// Verify mode set
	if !cmd.SignatureMode && !cmd.DeltaMode {
		logger(constants.ModeFlagMissingError, true)
		return false
	}

	// Verify files set for Signature mode
	if cmd.SignatureMode && (cmd.OriginalFile == "" || cmd.SignatureFile == "") {
		logger(constants.SignatureFlagsMissingError, true)
		return false
	}

	// Verify files set for Delta mode
	if cmd.DeltaMode {
		if cmd.SignatureMode && (cmd.UpdatedFile == "" || cmd.DeltaFile == "") {
			logger(constants.SignatureDeltaFlagsMissingError, true)
			return false
		} else if !cmd.SignatureMode && (cmd.SignatureFile == "" || cmd.UpdatedFile == "" || cmd.DeltaFile == "") {
			logger(constants.DeltaFlagsMissingError, true)
			return false
		}
	}

	return true
}
