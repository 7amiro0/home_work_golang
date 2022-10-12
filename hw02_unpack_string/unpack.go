package hw02unpackstring

import (
	"errors"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func dataValidation(letter rune, shield *bool,
	prevLetter *rune, resultText *string, shieldPrevLetter *bool) bool {
	skip := false
	if letter == '\\' && !*shield {
		*shield = true
		*shieldPrevLetter = true
		skip = true
	} else if unicode.IsDigit(letter) && !*shield && *prevLetter != 0 {
		*resultText += strings.Repeat(string(*prevLetter), int(letter-'1'))
		*prevLetter = letter
		*shieldPrevLetter = false
		skip = true
	}
	return skip
}

func Unpack(text string) (string, error) {
	resultText := ""
	shield := false
	prevLetter := rune(0)
	shieldPrevLetter := false
	for _, letter := range text {

		if unicode.IsDigit(prevLetter) && unicode.IsDigit(letter) && !shieldPrevLetter {
			return "", ErrInvalidString
		}
		if shield && unicode.IsLetter(letter) {
			return "", ErrInvalidString
		}

		if dataValidation(letter, &shield, &prevLetter, &resultText, &shieldPrevLetter) {
			continue
		}

		resultText += string(letter)
		prevLetter = letter
		shield = false
	}
	return resultText, nil
}
