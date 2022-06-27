package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHello(t *testing.T) {
	t.Run("should return `Hello, world!` string", func(t *testing.T) {
		result := Hello()
		expected := "Hello, world!"
		require.Equal(t, expected, result)
	})
}

func TestMain(t *testing.T) {
	t.Run("should call log function", func(t *testing.T) {
		invoked := false
		log = func(a ...interface{}) (n int, err error) {
			invoked = true
			return 0, nil
		}
		main()
		require.Equal(t, true, invoked)
	})
}