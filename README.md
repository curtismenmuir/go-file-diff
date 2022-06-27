# go-file-diff

## :collision: Important

- This project is built with Go Version: `go1.16.15`
- This project uses `Go Mod`

## :arrow_up: How to Setup Project

**Step 1:** git clone this repo

**Step 2:** Ensure Go is installed & configured on machine

**Step 5:** Download deps: `go mod download`

## :arrow_forward: How to Run Project

**Step 1:** Complete Setup instructions above

**Step 2:** Run app: `go run .`

## :rotating_light: Unit Tests

- Run `go test ./...` from the root directory

- For test coverage, run `go test ./... -coverprofile cp.out` from the root directory
  - `go tool cover -html=cp.out` can then be run (root dir) to open detailed coverage breakdown in internet browser
