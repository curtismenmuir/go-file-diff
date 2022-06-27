package constants

// Error messages
const (
	ModeFlagMissingError            string = "Error: Must set at least one mode"
	SignatureFlagsMissingError      string = "Error: Must provide Original & Signature files when enabling Signature mode"
	DeltaFlagsMissingError          string = "Error: Must provide Signature, Updated & Delta files when enabling Delta mode"
	SignatureDeltaFlagsMissingError string = "Error: Must provide Updated & Delta files when enabling Signature & Delta modes"
	SignatureNotImplementedError    string = "signature mode not implemented, coming soon"
	DeltaNotImplementedError        string = "delta mode not implemented, coming soon"
)
