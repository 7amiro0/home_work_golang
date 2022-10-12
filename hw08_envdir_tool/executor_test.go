package main

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRunCmd(t *testing.T) {
	t.Run("simple start", func(t *testing.T) {
		data := make(map[string]EnvValue)
		data["BAR"] = EnvValue{"bar", false}
		data["EMPTY"] = EnvValue{"", false}
		cmd := []string{"./testdata/echo.sh"}
		realResult := RunCmd(cmd, data)
		expectedResult := 0
		require.Equalf(t, expectedResult, realResult, "expect result %v not equal real result %v", expectedResult, realResult)
	})

	t.Run("empty data", func(t *testing.T) {
		data := make(map[string]EnvValue)
		cmd := []string{"/testdata/echo.sh"}
		realResult := RunCmd(cmd, data)
		expectedResult := 0
		require.Equalf(t, expectedResult, realResult, "expect result %v not equal real result %v", expectedResult, realResult)
	})

	t.Run("empty command", func(t *testing.T) {
		data := make(map[string]EnvValue)
		data["BAR"] = EnvValue{"bar", false}
		data["EMPTY"] = EnvValue{"", false}
		cmd := []string{}
		realResult := RunCmd(cmd, data)
		expectedResult := 0
		require.Equalf(t, expectedResult, realResult, "expect result %v not equal real result %v", expectedResult, realResult)
	})

	t.Run("without input data", func(t *testing.T) {
		returnCode := RunCmd([]string{"cd", "$HOME", "&&", "pwd"}, Environment{"HOME": {Value: "/home/samir/Downloads", NeedRemove: false}})

		require.Equal(t, 0, returnCode)
	})
}
