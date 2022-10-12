package hw02unpackstring

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUnpack(t *testing.T) {
	testTable := []struct {
		str      string
		expected string
		err      error
	}{
		{
			str:      "hel2o2",
			expected: "helloo",
			err:      nil,
		},
		{
			str:      `hel\2o`,
			expected: `hel2o`,
			err:      nil,
		},
		{
			str:      `sla\\2sh`,
			expected: `sla\\sh`,
			err:      nil,
		},
		{
			str:      `hel\2o`,
			expected: `hel2o`,
			err:      nil,
		},
		{
			str:      `45`,
			expected: ``,
			err:      ErrInvalidString,
		},
		{
			str:      `a45`,
			expected: ``,
			err:      ErrInvalidString,
		},
		//{
		//	str:      `asdf\`,
		//	expected: ``,
		//	err:      ErrInvalidString,
		//},
		{
			str:      `hel\\\2o`,
			expected: `hel\2o`,
			err:      nil,
		},
		{
			str:      `hel1o`,
			expected: `helo`,
			err:      nil,
		},
		{
			str:      ``,
			expected: ``,
			err:      nil,
		},
		{
			str:      "gap here ->\n5.",
			expected: "gap here ->\n\n\n\n\n.",
			err:      nil,
		},
		{
			str:      `qw\er`,
			expected: ``,
			err:      ErrInvalidString,
		},
		{
			str:      `qw\43er`,
			expected: `qw444er`,
			err:      ErrInvalidString,
		},
	}

	for _, testCase := range testTable {
		t.Run("start test", func(t *testing.T) {
			realResult, _ := Unpack(testCase.str)

			require.Equal(t, testCase.expected, realResult)

		})
	}
}
