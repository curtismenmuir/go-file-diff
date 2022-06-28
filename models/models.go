package models

type CMD struct {
	Verbose       bool   `json:"verbose"`
	SignatureMode bool   `json:"signatureMode"`
	DeltaMode     bool   `json:"deltaMode"`
	OriginalFile  string `json:"originalFile"`
	SignatureFile string `json:"signatureFile"`
	UpdatedFile   string `json:"updatedFile"`
	DeltaFile     string `json:"deltaFile"`
}