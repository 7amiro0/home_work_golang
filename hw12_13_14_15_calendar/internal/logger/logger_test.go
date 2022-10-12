package logger

import (
	"bufio"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

type testTable struct {
	name     string
	level    string
	text     string
	function func(*Logger, ...any)
}

func createLogger(level string, out string) *Logger {
	config := configLogger(level)
	config.OutputPaths = []string{out}
	config.ErrorOutputPaths = []string{out}
	log, err := config.Build()
	if err != nil {
		panic(err)
	}

	sugar := log.Sugar()

	return &Logger{
		logger: sugar,
	}
}

func TestLogger(t *testing.T) {
	tests := []testTable{
		{
			name:     "info",
			level:    "info",
			text:     "some info",
			function: (*Logger).Info,
		},

		{
			name:     "debug",
			level:    "debug",
			text:     "some debug",
			function: (*Logger).Debug,
		},

		{
			name:     "warn",
			level:    "warn",
			text:     "some warn",
			function: (*Logger).Warn,
		},

		{
			name:     "error",
			level:    "error",
			text:     "some error",
			function: (*Logger).Error,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			file, err := os.CreateTemp("", "")
			require.Nil(t, err)
			defer os.RemoveAll(file.Name())

			logger := createLogger(test.level, file.Name())

			test.function(logger, test.text)

			scan := bufio.NewScanner(file)
			for scan.Scan() {
				require.Equalf(t, test.text, scan.Text(), "expected %s, but got %s", test.text, scan.Text())
			}
		},
		)
	}
}
