package main

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestReadDir(t *testing.T) {
	t.Run("simple start", func(t *testing.T) {
		realResult, _ := ReadDir("./testdata/env")
		expectedResult := make(map[string]EnvValue)
		expectedResult["BAR"] = EnvValue{"bar", false}
		expectedResult["EMPTY"] = EnvValue{"", false}
		expectedResult["FOO"] = EnvValue{"   foo\nwith new line", false}
		expectedResult["HELLO"] = EnvValue{"\"hello\"", false}
		expectedResult["UNSET"] = EnvValue{"", true}
		for key, value := range realResult {
			t.Run(key, func(t *testing.T) {
				if _, ok := expectedResult[key]; !ok {
					require.Failf(t, "key %v not exist", key)
				}
				require.Equalf(t, expectedResult[key], value, "expected value \"%v\" not equal real value \"%v\"", expectedResult[key], value)
			})
		}
	})

	t.Run("not specified directory", func(t *testing.T) {
		expected := 0
		data, _ := ReadDir("")
		require.Equal(t, expected, len(data), "map should be empty")
	})

	t.Run("file name with = symbol", func(t *testing.T) {
		data, _ := ReadDir("./testdata/env")
		_, exist := data["some=name"]
		require.False(t, exist, "file name with \"=\" doesnt exist in map")
	})

	t.Run("no such directory", func(t *testing.T) {
		expected := 0
		data, _ := ReadDir("./blabla")
		require.Equal(t, expected, len(data), "map should be empty")
	})
}
