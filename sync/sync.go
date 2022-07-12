package sync

import (
	"errors"
	"fmt"
	"io"
	"math"
	"math/big"

	"github.com/curtismenmuir/go-file-diff/constants"
	"github.com/curtismenmuir/go-file-diff/utils"
)

var (
	logger           = utils.Logger
	rollBuffer       = roll
	chunk      int64 = 16           // 16 (bytes) is max chunk size for seed == 11
	seed       int64 = 11           // Prime number
	mod        int64 = 100000000009 // 10^11 + 9
)

// FileReader interface for mocking bufio.Reader
type Reader interface {
	Read(p []byte) (int, error)
	ReadByte() (byte, error)
}

// GenerateSignature() will create a rolling "hash" buffer and loop through file until it reaches EOF
// NOTE: Signature generation still to come
// Function returns `nil` when successful
// Function returns `error` when unsuccessful
func GenerateSignature(reader Reader, verbose bool) error {
	// Create buffer based on chunk size
	buffer, err := populateBuffer(reader, chunk)
	if err != nil {
		return err
	}

	logger(fmt.Sprintf("Initial Buffer = %q", buffer[:]), verbose)
	// Generate Weak hash of initial buffer
	hash := generateWeakHash(buffer, chunk)
	logger(fmt.Sprintf("Initial hash = %d\n", hash), verbose)
	// TODO Generate Strong hash
	// TODO Store values in Signature Table

	// Loop until EOF
	for {
		var initialByte byte
		var nextByte byte
		// Roll buffer to next position
		buffer, initialByte, nextByte, err = rollBuffer(reader, buffer)
		if err != nil {
			// Break loop when EOF returned
			if err.Error() == constants.EndOfFileError {
				break
			}

			// Handle errors
			return err
		}

		logger(fmt.Sprintf("Rolled Buffer = %q", buffer[:]), verbose)

		// Roll Weak hash
		hash = rollWeakHash(hash, initialByte, nextByte, chunk)
		logger(fmt.Sprintf("Rolled hash = %d\n", hash), verbose)
		// TODO Generate Strong hash
		// TODO Store values in Signature Table
	}

	// TODO return Signature table
	return nil
}

// generateWeakHash() will generate a `weak` hash of a byte array based on the Rabin–Karp algorithm
// EG hash = ((array[0] * seed^n-1) + (array[1] * seed^n-2) + ... + (array[n] * seed^0)) % mod
// Hash is classed as `weak` as there is potential for collisions
// This function will return `hash` once complete
func generateWeakHash(buffer []byte, chunkSize int64) int64 {
	multiplier := chunkSize - 1
	var hash int64 = 0
	for index := range buffer {
		// Generate hash value for buffer item -> (buffer[i] * (seed^multiplier))
		value := int64(buffer[index]) * int64(math.Pow(float64(seed), float64(multiplier)))
		// Add value to hash
		hash = hash + value
		// Reduce multiplier for next iteration
		multiplier--
	}

	// Mod output for final hash
	hash = modulo(hash, mod)
	return hash
}

// pop() will remove the first item from a provided buffer
// Function returns `updatedBuffer, initialByte`
// Note: initialByte is the item popped from buffer
func pop(buffer []byte) ([]byte, byte) {
	// Create new slice
	buf := make([]byte, 0)
	// Get initial byte
	initialByte := buffer[0]
	// Fill slice with buffer items from position 1 (pop)
	buf = append(buf, buffer[1:]...)
	return buf, initialByte
}

// modulo() will run a mod operation on 2 numbers and return the result
// math/big is used over the built-in mod operator as `%` does not implement Euclidean modulus
// Function returns `result` -> EG x % y
func modulo(x int64, y int64) int64 {
	return new(big.Int).Mod(big.NewInt(x), big.NewInt(y)).Int64()
}

// populateBuffer() will create a new buffer and populate it, based on `chuck` size, from the provided file reader
// Function will return `buffer, nil` when successful
// Function will return `emptyBuffer, EOF` error when reader reaches end of file
// Function will return `emptyBuffer, error` when unable to read from file
func populateBuffer(reader Reader, chunkSize int64) ([]byte, error) {
	// Create buffer based on chunk size
	buffer := make([]byte, chunkSize)
	// Fill buffer from file reader
	n, err := reader.Read(buffer)
	if err != nil {
		// Handle EOF error
		if err == io.EOF {
			return []byte{}, errors.New(constants.EndOfFileError)
		}

		return []byte{}, err
	}

	if n == 0 {
		// Handle EOF error
		return []byte{}, errors.New(constants.EndOfFileError)
	}

	return buffer, nil
}

// push() will append the provided byte to the end of the provided buffer
// Function returns `updatedBuffer`
func push(buffer []byte, item byte) []byte {
	return append(buffer, item)
}

// roll() will move the rolling hash function to the next position
// This will include: read next item from file; popping 1st item from buffer; pushing new item to end of buffer;
// Function will return `updatedBuffer, initialByte, nextByte, nil` when successful
// Note: initialByte = byte popped from first position
// Note: nextByte = byte pushed onto end of buffer
// Function will return `emptyBuffer, 0, 0, EOL` error when read EOF from file reader
// Function will return `emptyBuffer, 0, 0, error` when unable to read byte from file
func roll(reader Reader, buffer []byte) ([]byte, byte, byte, error) {
	// Read a byte from file reader
	nextByte, err := reader.ReadByte()
	if err != nil {
		// Handle EOF error
		if err == io.EOF {
			return []byte{}, 0, 0, errors.New(constants.EndOfFileError)
		}

		return []byte{}, 0, 0, err
	}

	// Pop initial byte
	buf, initialByte := pop(buffer)
	// Push new byte
	buf = push(buf, nextByte)
	return buf, initialByte, nextByte, nil
}

// rollWeakHash() will roll a hash value to the next position based on initial byte of hash + new byte to roll in
// EG newHash = ((((hash - ((initialByte * seed^n-1) % mod)) * seed) % mod) + nextByte) % mod
// This function will return `updatedHash` once complete
func rollWeakHash(hash int64, initialByte byte, nextByte byte, chunkSize int64) int64 {
	// Hash initialByte -> initialByte * seed^n-1
	hashedInitialByte := int64(initialByte) * int64(math.Pow(float64(seed), float64(chunkSize-1)))
	// Mod hashedInitialByte and remove from hash -> hash - (hashedInitialByte % mod)
	updatedHash := hash - modulo(hashedInitialByte, mod)
	// Multiply seed -> result * seed
	updatedHash = updatedHash * seed
	// Mod + add new byte -> result % mod + int64(nextByte)
	updatedHash = modulo(updatedHash, mod) + int64(nextByte)
	// Mod output to get final updated hash -> result % mod
	return modulo(updatedHash, mod)
}
