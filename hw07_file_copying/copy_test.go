package main

import (
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"math"
	"os"
	"testing"
)

func TestCopy(t *testing.T) {
	t.Run("Test: this is work?", func(t *testing.T) {
		err := Copy("testdata/input.txt", "out.txt", 0, 0)
		require.NoErrorf(t, err, "Fail: \"%v\" not equal \"nil\"", err)

		fromFile, _ := os.Open("testdata/input.txt")
		outFile, _ := os.Open("out.txt")
		fromBytes, _ := ioutil.ReadAll(fromFile)
		outBytes, _ := ioutil.ReadAll(outFile)

		require.Equalf(t, string(fromBytes), string(outBytes), "Fail: \"out.txt\" not equal \"testdata/input.txt\"\n")
		fromFile.Close()
		outFile.Close()
		os.Remove("out.txt")
	})

	t.Run("Test: offset exceeds file size", func(t *testing.T) {
		expectedError := ErrOffsetExceedsFileSize
		maxInt64 := int64(math.MaxInt64)
		realError := Copy("testdata/input.txt", "out.txt", maxInt64, 0)
		require.Equalf(t, expectedError, realError, "Fail: error \"%v\" not equal error \"%v\"", realError, expectedError)
	})

	t.Run("Test: copy jpg file", func(t *testing.T) {
		realError := Copy("testdata/123.jpg", "out.jpg", 0, 0)
		require.NoErrorf(t, realError, "Fail: error \"%v\" not equal error \"nil\"", realError)

		fromFile, _ := os.Open("testdata/123.jpg")
		outFile, _ := os.Open("out.jpg")
		fromBytes, _ := ioutil.ReadAll(fromFile)
		outBytes, _ := ioutil.ReadAll(outFile)

		require.Equalf(t, string(fromBytes), string(outBytes), "Fail: \"out.txt\" not equal \"testdata/input.txt\"\n")
		fromFile.Close()
		outFile.Close()
		os.Remove("out.jpg")
	})

	t.Run("Test: unsupported open file", func(t *testing.T) {
		expectedError := ErrUnsupportedFile
		realError := Copy("", "out.txt", 0, 0)
		require.Equalf(t, expectedError, realError, "Fail: error \"%v\" not equal error \"%v\"", realError, expectedError)
	})

	t.Run("Test: unsupported create file", func(t *testing.T) {
		expectedError := ErrUnsupportedFile
		realError := Copy("testdata/input.txt", "", 0, 0)
		require.Equalf(t, expectedError, realError, "Fail: error \"%v\" not equal error \"%v\"", realError, expectedError)
	})

	t.Run("Test: copy dev/null", func(t *testing.T) {
		expectedError := ErrOffsetExceedsFileSize
		maxInt64 := int64(math.MaxInt64)
		realError := Copy("../../../../../dev/null", "out.txt", maxInt64, 0)
		require.Equalf(t, expectedError, realError, "Fail: error \"%v\" not equal error \"%v\"", realError, expectedError)
	})
}
