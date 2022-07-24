package sync

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"math"
	"math/big"

	"github.com/curtismenmuir/go-file-diff/constants"
	"github.com/curtismenmuir/go-file-diff/models"
	"github.com/curtismenmuir/go-file-diff/utils"
)

var (
	logger                 = utils.Logger
	initialiseBuffer       = populateBuffer
	rollBuffer             = roll
	chunk            int64 = 16           // 16 (bytes) is max chunk size for seed == 11
	seed             int64 = 11           // Prime number
	mod              int64 = 100000000009 // 10^11 + 9
)

// FileReader interface for mocking bufio.Reader.
type Reader interface {
	Read(p []byte) (int, error)
	ReadByte() (byte, error)
}

// compareChecksums() will search for a Weak hash in provided Signature.
// When match is found with Weak hash, function will generate Strong hash and compare against Signature item.
// Function will return `true, item.Head, item.Tail` when successfully found block in Signature (EG When Weak & Strong hashes match Signature item).
// Function will return `false, -1, -1` when unable to find block in Signature.
func compareChecksums(signature models.Signature, buffer []byte, weakHash int64, verbose bool) (bool, int, int) {
	// Search Signature for Weak hash
	if item, exists := signature[weakHash]; exists {
		// Generate Strong hash of buffer
		strongHash := generateStrongHash(buffer, chunk)
		logger(fmt.Sprintf("Strong hash = %s", strongHash), verbose)
		// Verify if Strong hash also matches Signature item
		if strongHash == item.Hash {
			logger("Block found\n", verbose)
			return true, item.Head, item.Tail
		}
	}

	logger("Block missing\n", verbose)
	return false, -1, -1
}

// GenerateDelta() will create a Delta changeset of how to update a provided file Signature to match an updated version of the file.
// Delta will contain a list of reusable blocks from the original file, and where they should be added to match the Updated file.
// Delta will also contain a list of new blocks which can be applied to the file to sync latest modifications.
// Function will return `delta, nil` when generated Delta successfully.
// Function will return `emptyDelta, UpdatedFileHasNoChangesError` when Updated file has no changes from Original.
// Function will return `emptyDelta, error` when unable to populate buffer from file.
// Function will return `emptyDelta, error` when unable to read data from file to roll buffer.
func GenerateDelta(reader Reader, signature models.Signature, verbose bool) (models.Delta, error) {
	blockHead := 0
	deltaHead := 0
	deltaTail := int(chunk) - 1
	delta := make(models.Delta)
	initialBlockMatches := true
	var block models.Block
	// Create buffer based on chunk size
	buffer, err := initialiseBuffer(reader, chunk)
	if err != nil {
		return models.Delta{}, err
	}

	logger(fmt.Sprintf("Initial Buffer = %q", buffer[:]), verbose)
	// Generate Weak hash of initial buffer
	weakHash := generateWeakHash(buffer, chunk)
	logger(fmt.Sprintf("Weak hash = %d", weakHash), verbose)
	// Search Signature for match on initial buffer
	exists, head, tail := compareChecksums(signature, buffer, weakHash, verbose)
	if exists {
		// Create new matched block
		block = models.Block{Head: head, Tail: tail, IsModified: !exists, Value: []byte{}}
	} else {
		// Create new missing block and record initial block does not match
		block = models.Block{Head: deltaHead, Tail: deltaHead, IsModified: !exists, Value: []byte{buffer[0]}}
		initialBlockMatches = false
	}

	// Loop until EOF
	for {
		var initialByte, nextByte byte
		var rollExists bool
		var rollHead, rollTail int
		// Roll buffer to next position
		buffer, initialByte, nextByte, err = rollBuffer(reader, buffer)
		if err != nil {
			// Break loop when EOF returned
			if err.Error() == constants.EndOfFileError {
				// Add final block to Delta
				delta[blockHead] = block
				logger(fmt.Sprintf("Final Block added to Delta: %+v\n", block), verbose)
				if block.IsModified {
					logger(fmt.Sprintf("Final Block Value = %q\n", block.Value[:]), verbose)
				}

				break
			}

			// Handle errors
			return models.Delta{}, err
		}

		logger(fmt.Sprintf("Rolled Buffer = %q", buffer[:]), verbose)
		// Increment Delta position
		deltaHead++
		deltaTail++
		// Roll Weak hash
		weakHash = rollWeakHash(weakHash, initialByte, nextByte, chunk)
		logger(fmt.Sprintf("Rolled hash = %d", weakHash), verbose)
		// Search Signature for match on rolled buffer
		rollExists, rollHead, rollTail = compareChecksums(signature, buffer, weakHash, verbose)
		if rollExists {
			// Match found in Signature, generate matched block
			block, blockHead, initialBlockMatches = generateMatchedBlock(delta, block, exists, initialBlockMatches, blockHead, deltaHead, rollHead, rollTail, rollExists, verbose)
		} else {
			// No match found in Signature, generate missing block
			block, blockHead = generateMissingBlock(delta, block, exists, initialBlockMatches, blockHead, nextByte, buffer, verbose)
		}

		// Record if match found for next iteration
		exists = rollExists
	}

	logger(fmt.Sprintf("Delta: %+v\n", delta), verbose)

	// Verify if Delta contains any modifications for Original file
	if len(delta) == 1 && !delta[0].IsModified {
		return models.Delta{}, errors.New(constants.UpdatedFileHasNoChangesError)
	}

	return delta, nil
}

// generateMatchedBlock() will generate a new matched block after adding previous missing block to Delta (only added to delta when applicable).
// If previous roll was a match, then function will increase blocks tail position.
// If previous roll was a missing block at the start of the file, then function will add provided block to Delta and return a new matched block.
// Note: Missing initial block will be found at start of buffer (EG not rolled in).
// If previous roll was a missing block but not found at beginning of file, then function will reduce block to remove any matched bytes, add block to Delta, and return a new matched block.
// Note: Function reduces block as final roll will include 15 bytes of next match (EG rolling 16 byte buffer).
// Function returns `block, blockHead, initialBlockMatches` upon completion.
// Note: Function will update original instance of provided `Delta` as maps are reference types.
func generateMatchedBlock(delta models.Delta, block models.Block, exists bool, initialBlockMatches bool, blockHead int, deltaHead int, rollHead int, rollTail int, rollExists bool, verbose bool) (models.Block, int, bool) {
	// Verify if previous block matched
	if exists {
		// Increase blocks tail position when rolled buffer still matches
		block.Tail++
	} else {
		// Verify if updating initial missing block
		if !initialBlockMatches {
			// If initial block is missing then block will contain only updated values
			initialBlockMatches = true
		} else {
			// Reduce block to remove following matched characters
			// EG last 15 characters of buffer will contain start of next matched block due to rolling function (EG buffer size == 16)
			block.Tail = block.Tail + 1 - int(chunk)
			missingValues := make([]byte, 0)
			missingValues = append(missingValues, block.Value[0:block.Tail+1]...)
			block.Value = missingValues
		}

		// Add missing block to Delta
		delta[blockHead] = block
		logger(fmt.Sprintf("Missing Block added to Delta: %+v", block), verbose)
		logger(fmt.Sprintf("Missing Block Position: %d", blockHead), verbose)
		logger(fmt.Sprintf("Missing Block Value = %q\n", block.Value[:]), verbose)
		// Update position for next matching block
		blockHead = deltaHead
		// Create new matching block
		block = models.Block{Head: rollHead, Tail: rollTail, IsModified: !rollExists, Value: []byte{}}

	}

	return block, blockHead, initialBlockMatches
}

// generateMissingBlock() will generate a new missing block after adding previous matched block to Delta (only added to delta when applicable).
// If previous roll was a match, then function will add matched block to Delta, update block head to new position, and return a new missing block.
// If previous roll was a missing block at the start of the file, the function will add byte from beginning of buffer to block Value & increment block Tail position.
// Note: Use buffer first item as missing initial block will be found at start of buffer (EG not rolled in).
// If previous roll was a missing block but not at beginning of file, the function will add next rolled byte to block Value & increment block Tail position.
// Note: Use nextByte as missing block will be added to end of buffer (EG rolling 16 byte buffer).
// Function returns `block, blockHead` upon completion.
// Note: Function will update original instance of provided `Delta` as maps are reference types.
func generateMissingBlock(delta models.Delta, block models.Block, exists bool, initialBlockMatches bool, blockHead int, nextByte byte, buffer []byte, verbose bool) (models.Block, int) {
	// Verify if previous block matched
	if exists {
		// Add matching block to Delta
		delta[blockHead] = block
		logger(fmt.Sprintf("Matched Block added to Delta: %+v\n", block), verbose)
		// Update position for next missing block
		blockHead = blockHead + block.Tail - block.Head + 1
		// Create new missing block
		block = models.Block{Head: 0, Tail: 0, IsModified: exists, Value: []byte{nextByte}}
	} else {
		// Verify if updating initial missing block
		if !initialBlockMatches {
			// Use initial buffer position when first block does not match
			block.Value = append(block.Value, buffer[0])
		} else {
			// Add new byte (EG rolled value)
			block.Value = append(block.Value, nextByte)
		}

		// Increase blocks last position when rolled buffer still missing
		block.Tail++
	}

	return block, blockHead
}

// GenerateSignature() will create a file Signature from a provided file reader.
// Signature will contain a `weak` rolling hash of the file in 16 byte chunks.
// Signature will also contain a strong hash of each chunk to avoid collisions when generating Delta.
// Function returns `Signature, nil` when successful.
// Function returns `emptySignature, error` when unsuccessful.
func GenerateSignature(reader Reader, verbose bool) (models.Signature, error) {
	head := 0
	tail := int(chunk) - 1
	signature := make(models.Signature, 0)
	// Create buffer based on chunk size
	buffer, err := initialiseBuffer(reader, chunk)
	if err != nil {
		return models.Signature{}, err
	}

	logger(fmt.Sprintf("Initial Buffer = %q", buffer[:]), verbose)
	// Generate Weak hash of initial buffer
	weakHash := generateWeakHash(buffer, chunk)
	logger(fmt.Sprintf("Weak hash = %d", weakHash), verbose)
	// Generate Strong hash of buffer
	strongHash := generateStrongHash(buffer, chunk)
	logger(fmt.Sprintf("Strong hash = %s\n", strongHash), verbose)
	// Store values in Signature
	signature[weakHash] = models.StrongSignature{Hash: strongHash, Head: head, Tail: tail}
	// Loop until EOF
	for {
		var initialByte byte
		var nextByte byte
		head++
		tail++
		// Roll buffer to next position
		buffer, initialByte, nextByte, err = rollBuffer(reader, buffer)
		if err != nil {
			// Break loop when EOF returned
			if err.Error() == constants.EndOfFileError {
				break
			}

			// Handle errors
			return models.Signature{}, err
		}

		logger(fmt.Sprintf("Rolled Buffer = %q", buffer[:]), verbose)
		// Roll Weak hash
		weakHash = rollWeakHash(weakHash, initialByte, nextByte, chunk)
		logger(fmt.Sprintf("Rolled hash = %d", weakHash), verbose)
		// Generate Strong hash of updated buffer
		strongHash = generateStrongHash(buffer, chunk)
		logger(fmt.Sprintf("Strong hash = %s\n", strongHash), verbose)
		// Add hashes to Signature
		signature[weakHash] = models.StrongSignature{Hash: strongHash, Head: head, Tail: tail}
	}

	logger(fmt.Sprintf("Signature: %+v\n", signature), verbose)
	return signature, nil
}

// generateStrongHash() will hash a provided buffer with SHA-256.
// Function returns final `hash` value encoded as a hex string.
func generateStrongHash(buffer []byte, chunkSize int64) string {
	sha := sha256.New()
	sha.Write(buffer)
	return hex.EncodeToString(sha.Sum(nil))
}

// generateWeakHash() will generate a `weak` hash of a byte array based on the Rabin–Karp algorithm.
// EG hash = ((array[0] * seed^n-1) + (array[1] * seed^n-2) + ... + (array[n] * seed^0)) % mod;
// Hash is classed as `weak` as there is potential for collisions.
// Function returns `hash`.
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

// modulo() will run a mod operation on 2 numbers and return the result.
// math/big is used over the built-in mod operator as `%` does not implement Euclidean modulus.
// Function returns `result` -> EG x % y;
func modulo(x int64, y int64) int64 {
	return new(big.Int).Mod(big.NewInt(x), big.NewInt(y)).Int64()
}

// pop() will remove the first item from a provided buffer.
// Function returns `updatedBuffer, initialByte`.
// Note: initialByte is the item popped from buffer.
func pop(buffer []byte) ([]byte, byte) {
	// Create new slice
	buf := make([]byte, 0)
	// Get initial byte
	initialByte := buffer[0]
	// Fill slice with buffer items from position 1 (pop)
	buf = append(buf, buffer[1:]...)
	return buf, initialByte
}

// populateBuffer() will create a new buffer and populate it, based on `chuck` size, from the provided file reader.
// Function will return `buffer, nil` when successful.
// Function will return `emptyBuffer, EOF` error when reader reaches end of file.
// Function will return `emptyBuffer, error` when unable to read from file.
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

// push() will append the provided byte to the end of the provided buffer.
// Function returns `updatedBuffer`.
func push(buffer []byte, item byte) []byte {
	return append(buffer, item)
}

// roll() will move the rolling hash function to the next position.
// This will include: read next item from file; popping 1st item from buffer; pushing new item to end of buffer;
// Function will return `updatedBuffer, initialByte, nextByte, nil` when successful.
// Note: initialByte = byte popped from first position.
// Note: nextByte = byte pushed onto end of buffer.
// Function will return `emptyBuffer, 0, 0, EOF_error` when file reader reaches EOF.
// Function will return `emptyBuffer, 0, 0, error` when unable to read byte from file.
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

// rollWeakHash() will roll a hash value to the next position based on initial byte of hash + new byte to roll in.
// EG newHash = ((((hash - ((initialByte * seed^n-1) % mod)) * seed) % mod) + nextByte) % mod;
// This function will return `updatedHash` once complete.
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
