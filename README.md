# go-file-diff

## :collision: Important

- This project has been built using Go version: `go1.18.3`

## :arrow_up: How to Setup Project

**Step 1:** git clone this repo

**Step 2:** Ensure Go is installed & configured on machine

**Step 3:** Download deps: `go mod download`

## :arrow_forward: How to Run Project

**Step 1:** Complete Setup instructions above

**Step 2:** Run app: `go run . <CMD Args>`

- NOTE: See `CMD Commands` section below for more details

## :cop: Linting

- This project uses [golangci-lint](https://github.com/golangci/golangci-lint) for linting
- Once linter is installed, run tool with: `golangci-lint run` 
- [Local installation Guide](https://golangci-lint.run/usage/install/#local-installation)
- [Quick Start Guide](https://golangci-lint.run/usage/quick-start)

## :rotating_light: Unit Tests

- Run `go test ./...` from the root directory

- For test coverage, run `go test ./... -coverprofile cp.out` from the root directory
  - `go tool cover -html=cp.out` can then be run (root dir) to open detailed coverage breakdown in internet browser

## :computer: CMD Commands

- Signature Mode: `go run . -signatureMode -original=a.txt -signature=b.txt`
- Delta Mode: `go run . -deltaMode -signature=a.txt -updated=b.txt -delta=c.txt`
- Signature + Delta Mode: `go run . -signatureMode -deltaMode -original=a.txt -signature=b.txt -updated=c.txt -delta=d.txt`
