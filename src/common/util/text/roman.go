// Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.

package text

var romanToArabic = map[string]string{
	"i":   "1",
	"ii":  "2",
	"iii": "3",
	"iv":  "4",
	"v":   "5",
	"vi":  "6",
}

// IsRomanNumeral returns true if s is a roman numeral between i and vi.
func IsRomanNumeral(s string) bool {
	_, ok := romanToArabic[s]
	return ok
}

// RomanToArabicNumberals converts roman numerals to arabic numerals.
// Its purpose is to help parse NYHA class and Fitzpatrick skin type.
func RomanToArabicNumerals(s string) string {
	if a, ok := romanToArabic[s]; ok {
		return a
	}
	return s
}
