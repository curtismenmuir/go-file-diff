package cmd

import (
	"flag"
	"fmt"

	"github.com/curtismenmuir/go-file-diff/constants"
	"github.com/curtismenmuir/go-file-diff/utils"
)

var (
	logger       = utils.Logger
	defineBool   = flag.Bool
	defineString = flag.String
)

func ParseCMD() (bool, bool, bool, string, string, string, string) {
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

	// Log variables
	logger(fmt.Sprintf("Verbose Logging: %t", *verbose), *verbose)
	logger(fmt.Sprintf("Signature Mode: %t", *signatureMode), *verbose)
	logger(fmt.Sprintf("Delta Mode: %t", *deltaMode), *verbose)
	logger(fmt.Sprintf("Original File: %s", *originalFile), *verbose)
	logger(fmt.Sprintf("Signature File: %s", *signatureFile), *verbose)
	logger(fmt.Sprintf("Updated File: %s", *updatedFile), *verbose)
	logger(fmt.Sprintf("Delta File: %s", *deltaFile), *verbose)

	return *verbose, *signatureMode, *deltaMode, *originalFile, *signatureFile, *updatedFile, *deltaFile
}

func VerifyCMD(verbose bool, signatureMode bool, deltaMode bool, originalFile string, signatureFile string, updatedFile string, deltaFile string) bool {
	// Verify mode set
	if !signatureMode && !deltaMode {
		logger(constants.ModeFlagMissingError, true)
		return false
	}

	// Verify files set for Signature mode
	if signatureMode && (originalFile == "" || signatureFile == "") {
		logger(constants.SignatureFlagsMissingError, true)
		return false
	}

	// Verify files set for Delta mode
	if deltaMode {
		if signatureMode && (updatedFile == "" || deltaFile == "") {
			logger(constants.SignatureDeltaFlagsMissingError, true)
			return false
		} else if !signatureMode && (signatureFile == "" || updatedFile == "" || deltaFile == "") {
			logger(constants.DeltaFlagsMissingError, true)
			return false
		}
	}

	return true
}
