package hw09structvalidator

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		meta   json.RawMessage
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}

	StructOfStruct struct {
		Field User `validate:"nested"`
	}

	NotSuchTag struct {
		Name string `validate:"letter:o|len:abs"`
		Age  int    `validate:"len:5|max:"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in: User{
				ID:     "111111111111111111111111111111111111",
				Name:   "Johnny",
				Age:    19,
				Email:  "sujadsga@gmail.com",
				Role:   "stuff",
				Phones: []string{"58246511346", "58241463463", "58241498075"},
				meta:   nil,
			},
			expectedErr: nil,
		},

		{
			in: StructOfStruct{
				Field: User{
					ID:     "111111111111111111111111111111111111",
					Name:   "Gyro",
					Age:    22,
					Email:  "sdagienv@gmail.com",
					Role:   "admin",
					Phones: []string{"11246511346", "16241463463", "13241498075"},
					meta:   nil,
				},
			},
			expectedErr: nil,
		},

		{
			in: App{
				Version: "6.2.5",
			},
			expectedErr: nil,
		},

		{
			in: Token{
				Header:    nil,
				Payload:   nil,
				Signature: nil,
			},
			expectedErr: nil,
		},

		{
			in: Response{
				Code: 200,
				Body: "bla",
			},
			expectedErr: nil,
		},

		{
			in: NotSuchTag{
				Name: "Robert",
				Age:  32,
			},
			expectedErr: ValidationErrors{
				{
					Field: "Name",
					Err:   fmt.Errorf(ErrorInvalidParam.Error(), "letter"),
				},
				{
					Field: "Name",
					Err:   fmt.Errorf(ErrorInvalidArg.Error(), "abs", "len"),
				},
				{
					Field: "Age",
					Err:   fmt.Errorf(ErrorInvalidParam.Error(), "len"),
				},
				{
					Field: "Age",
					Err:   fmt.Errorf(ErrorInvalidArg.Error(), "", "max"),
				},
			},
		},

		{
			in: App{
				Version: "681482",
			},
			expectedErr: ValidationErrors{
				{
					Field: "Version",
					Err:   fmt.Errorf(ErrorLen.Error(), 6, 5),
				},
			},
		},

		{
			in: Response{
				Code: 32,
				Body: "bla",
			},
			expectedErr: ValidationErrors{
				{
					Field: "Code",
					Err:   fmt.Errorf(ErrorIn.Error(), 32, "[200 404 500]"),
				},
			},
		},

		{
			in: User{
				ID:     "111111111111111111111111111111111111",
				Name:   "Johnny",
				Age:    19,
				Email:  "sujadsga@gmail,com",
				Role:   "stuf",
				Phones: []string{"582465113466", "582414636463", "158241498075"},
				meta:   nil,
			},
			expectedErr: ValidationErrors{
				{
					Field: "Email",
					Err:   fmt.Errorf(ErrorRegexp.Error(), "sujadsga@gmail,com", "^\\w+@\\w+\\.\\w+$"),
				},
				{
					Field: "Role",
					Err:   fmt.Errorf(ErrorIn.Error(), "stuf", "[admin stuff]"),
				},
				{
					Field: "Phones",
					Err: ValidationErrors{
						{
							Field: "Phones",
							Err:   fmt.Errorf(ErrorLen.Error(), 12, 11),
						},
						{
							Field: "Phones",
							Err:   fmt.Errorf(ErrorLen.Error(), 12, 11),
						},
						{
							Field: "Phones",
							Err:   fmt.Errorf(ErrorLen.Error(), 12, 11),
						},
					},
				},
			},
		},

		{
			in: StructOfStruct{
				Field: User{
					ID:     "2111111111111111111111111111111111111",
					Name:   "Gyro",
					Age:    22,
					Email:  "sdagienv@gmail.com",
					Role:   "person",
					Phones: []string{"11246511346", "16241463463", "13241498075"},
					meta:   nil,
				},
			},
			expectedErr: ValidationErrors{
				{
					Field: "Field",
					Err: ValidationErrors{
						{
							Field: "ID",
							Err:   fmt.Errorf(ErrorLen.Error(), 37, 36),
						},
						{
							Field: "Role",
							Err:   fmt.Errorf(ErrorIn.Error(), "person", "[admin stuff]"),
						},
					},
				},
			},
		},

		{
			in:          10,
			expectedErr: fmt.Errorf(ErrorInvalidInputData.Error()),
		},

		{
			in: struct {
				Apps []App `validate:"nested"`
			}{
				Apps: []App{
					{Version: "v1.34"},
					{Version: "v1.345"},
					{Version: "v1.3"},
				},
			},
			expectedErr: ValidationErrors{
				{
					Field: "Apps",
					Err: ValidationErrors{
						{
							Field: "Apps",
							Err: ValidationErrors{
								{
									Field: "Version",
									Err:   fmt.Errorf(ErrorLen.Error(), 6, 5),
								},
							},
						},
						{
							Field: "Apps",
							Err: ValidationErrors{
								{
									Field: "Version",
									Err:   fmt.Errorf(ErrorLen.Error(), 4, 5),
								},
							},
						},
					},
				},
			},
		},

		{
			in: struct {
				Number8  int8  `validate:"max:10"`
				Number16 int16 `validate:"min:10"`
				Number32 int32 `validate:"max:1000"`
				Number64 int64 `validate:"min:10"`
				Number   int   `validate:"min:10|max:20"`
			}{
				Number8:  9,
				Number16: 20,
				Number32: 999,
				Number64: 32,
				Number:   17,
			},
			expectedErr: nil,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i+1), func(t *testing.T) {
			tt := tt
			//t.Parallel()
			err := Validate(tt.in)
			if tt.expectedErr == nil {
				require.Nilf(t, err, "expected nil but received %q", err)
			} else if tt.expectedErr != nil {
				//fmt.Printf("%q %T\n%q %T\n", tt.expectedErr, tt.expectedErr, err, err)
				require.Equalf(t, tt.expectedErr, err, "expected %q but received %q", tt.expectedErr, err)

				require.True(t, reflect.DeepEqual(tt.expectedErr, err))
			}
		})
	}
}
