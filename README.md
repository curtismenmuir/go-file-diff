# go-file-diff

## :collision: Important

- This project has been built using Go version: `go1.18.3`
- This project will diff an Original + Updated version of a file to produce a changeset on how to update the Original version to sync latest changes.
- This project implements a 16-byte rolling hash algorithm for evaluating differences between the 2 files.
  - Rolling hash algorithm is based on the `Rabin–Karp algorithm`.
- This project is based on the [rdiff](https://linux.die.net/man/1/rdiff) application.
- Delta changeset will evaluate:
  - Chunk changes and/or additions
  - Chunk removals
  - Additions between chunks with shifted original chunks

## :memo: Description

This project can be used to compare 2 versions of a file, establish what has changed, and produce a `Delta` changeset of how the Original version can be patched to sync the latest changes.
- NOTE: `patch` functionality is out of scope for now, but will be added in the future!

This can be used with 2 files on the same machine, or used to update files across different machines.

### Example flow for distributed files

- Machine 1 and Machine 2 both have a copy of the same file
- Machine 1 has local updates to the file and these should be synced with Machine 2
- Machine 2 generates a `Signature` of their original copy of the file: `./go-file-diff -signatureMode -original=original.txt -signature=sig.txt`
- Machine 2 sends `Signature` file to Machine 1
- Machine 1 generates a `Delta` of their local changes using provided `Signature` file: `./go-file-diff -deltaMode -signature=sig.txt -updated=updated.txt -delta=delta.txt`
- Machine 1 returns `Delta` file to Machine 2
- Machine 2 uses the `Delta` file to `Patch` their original version of the file to sync latest changes
    - NOTE: `Patch` functionality coming soon!

## :soon: Future Improvements

- Add `Dockerfile` 
  - Use `docker-compose` for mounting host volume into container?
- Implement `Patch` functionality
- Performance testing
- Setup CI pipeline
  - CircleCI free account?
- Setup go channels for processing Signature Weak + Strong hashes concurrently?

## :arrow_up: How to Setup Project

**Step 1:** git clone this repo

**Step 2:** Ensure Go `v1.18.3` is installed & configured on machine

- NOTE: Go should be installed with [gvm](https://github.com/moovweb/gvm) for managing multiple go versions

**Step 3:** Download deps: `go mod download`

## :arrow_forward: How to Run Project for Development

**Step 1:** Complete Setup instructions above

**Step 2:** Run app: `go run . <CMD Args>`

- NOTE: See `CMD Commands` section below for more details

## :rocket: How to Run Project for Release

**Step 1:** Complete Setup instructions above

**Step 2:** Build release app: `go build`

**Step 3:** Run release app: `./go-file-diff <CMD Args>`

- EG `./go-file-diff -signatureMode -original=original.txt -signature=sig.txt -v`

## :bulb: CMD Commands

| Command        | Example usage             | Description   | 
| -------------- | ------------------------- | ------------- |
| -signatureMode | `-signatureMode`          | Enables Signature generation. |
| -deltaMode     | `-deltaMode`              | Enables Delta generation. |
| -original      | `-original=SomeFile.txt`  | Name of Original file used for Signature generation. |
| -signature     | `-signature=SomeFile.txt` | Name of Signature file. In Signature mode, this will be used as Output file. In Delta mode, this will be used as an input file. |
| -updated       | `-updated=SomeFile.txt`   | Name of Updated file used for Delta generation. |
| -delta         | `-delta=SomeFile.txt`     | Name of Delta file. In Delta mode, this will be used as an Output file. |
| -v             | `-v`                      | Enables verbose logging. |

**NOTE:** Relative file paths should be used to access files in different folders from the application. EG:

- `./SomeFolder/SomeFile.txt`
- `../../AnotherFile.txt`

## :computer: Example Usage

- Signature Mode: `./go-file-diff -signatureMode -original=original.txt -signature=sig.txt -v`
- Delta Mode: `./go-file-diff -deltaMode -signature=Outputs/sig.txt -updated=updated.txt -delta=delta.txt -v`
- Signature + Delta Mode: `./go-file-diff -signatureMode -deltaMode -original=original.txt -signature=sig.txt -updated=updated.txt -delta=delta.txt -v`

## :rotating_light: Unit Tests

- Run `go test ./...` from the root directory

- For test coverage, run `go test ./... -coverprofile cp.out` from the root directory
  - `go tool cover -html=cp.out` can then be run (root dir) to open detailed coverage breakdown in internet browser

## :cop: Linting
- Run linter: `golangci-lint run` 
- NOTE: This project uses [golangci-lint](https://github.com/golangci/golangci-lint) for linting
  - [Local installation Guide](https://golangci-lint.run/usage/install/#local-installation)
  - [Quick Start Guide](https://golangci-lint.run/usage/quick-start)