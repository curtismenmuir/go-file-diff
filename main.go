package main

import (
	"errors"
	"fmt"

	"github.com/curtismenmuir/go-file-diff/cmd"
	"github.com/curtismenmuir/go-file-diff/constants"
	"github.com/curtismenmuir/go-file-diff/models"
	"github.com/curtismenmuir/go-file-diff/utils"
)

var (
	logger    = utils.Logger
	parseCMD  = cmd.ParseCMD
	verifyCMD = cmd.VerifyCMD
)

func getSignature(cmd models.CMD) error {
	return errors.New(constants.SignatureNotImplementedError)
}

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
			logger(fmt.Sprintf("Error: %s", err.Error()), true)
			return
		}
	}

	// Generate Delta
	if cmd.DeltaMode {
		err := getDelta(cmd)
		if err != nil {
			logger(fmt.Sprintf("Error: %s", err.Error()), true)
			return
		}
	}
}
