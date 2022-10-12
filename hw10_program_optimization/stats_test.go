//go:build !bench
// +build !bench

package hw10programoptimization

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetDomainStat(t *testing.T) {
	data := `{"Id":1,"Name":"Howard Mendoza","Username":"0Oliver","Email":"aliquid_qui_ea@Browsedrive.gov","Phone":"6-866-899-36-79","Password":"InAQJvsq","Address":"Blackbird Place 25"}
			 {"Id":2,"Name":"Jesse Vasquez","Username":"qRichardson","Email":"mLynch@broWsecat.com","Phone":"9-373-949-64-00","Password":"SiZLeNSGn","Address":"Fulton Hill 80"}
			 {"Id":3,"Name":"Clarence Olson","Username":"RachelAdams","Email":"RoseSmith@Browsecat.com","Phone":"988-48-97","Password":"71kuz3gA5w","Address":"Monterey Park 39"}
			 {"Id":4,"Name":"Gregory Reid","Username":"tButler","Email":"5Moore@Teklist.net","Phone":"520-04-16","Password":"r639qLNu","Address":"Sunfield Park 20"}
			 {"Id":5,"Name":"Janice Rose","Username":"KeithHart","Email":"nulla@Linktype.com","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
			 {"Id":6,"Name":"Ivan Ivanovich","Username":"doubleIvan","Email":"invalidEmailcom","Phone":"+91155125126","Password":"1afds987a","Address":"Sheila Plaza, 320 Port Myron, MD 81178-3660"}
			 {"Id":7,"Name":"Gregory Reid","Username":"tButler","Email":"DuMpMailE.com@gmail","Phone":"520-04-16","Password":"r639qLNu","Address":"Sunfield Park 20"}`

	tableTest := []struct {
		input    string
		expected DomainStat
	}{
		{
			input: "com",
			expected: DomainStat{
				"browsecat.com": 2,
				"linktype.com":  1,
			},
		},

		{
			input: "gov",
			expected: DomainStat{
				"browsedrive.gov": 1,
			},
		},

		{
			input: "net",
			expected: DomainStat{
				"teklist.net": 1,
			},
		},

		{
			input: "cOm",
			expected: DomainStat{
				"browsecat.com": 2,
				"linktype.com":  1,
			},
		},

		{
			input: "gmail",
			expected: DomainStat{
				"gmail": 1,
			},
		},
	}

	for index, test := range tableTest {
		t.Run(fmt.Sprintf("case: %v", index+1), func(t *testing.T) {
			result, err := GetDomainStat(bytes.NewBufferString(data), test.input)
			require.NoError(t, err)
			require.Equal(t, test.expected, result)
		})
	}
}
