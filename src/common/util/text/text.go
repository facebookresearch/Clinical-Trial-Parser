// Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.

package text

import (
	"regexp"
	"sort"
	"strings"
	"unicode"

	"github.com/facebookresearch/Clinical-Trial-Parser/src/common/util/slice"

	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

var (
	reSpaceBoundary    = regexp.MustCompile(`[&\s\n,\.;:\-_\+\?!*"/\\\(\)\[\]]`)
	reNoSpaceBoundary  = regexp.MustCompile("['`“”]")
	reWhitespace       = regexp.MustCompile(`\s+`)
	reSentenceBoundary = regexp.MustCompile(`[\.\?!;]([\.\?!;\s\n]+|\z)`)

	reOr = regexp.MustCompile(`(\s*/\s*=\s*)`)

	reNumber = regexp.MustCompile(`\b(([<≤>≥]{1})\s*[0-9]*[\.,]?[0-9]+\s*(?:%|\w+(?:/\w+)?)) `)

	reHalfBoundedRelation = regexp.MustCompile(`(((?:(?:^| )[\(\)\w]+){1,8})\s*([<≤>≥=]{1,2})\s*([\d]*[\.,]?[\d]+\s*\S+(?:\s*(?:institutional|the)?\s*(?:uln|upper limit(?: of)?(?: institutional)? normal)|(?:\s*[\d^]*\s*\S*\s*/\s*\S+))?))`)

	reSlash = regexp.MustCompile(`\s*/\s*`)
)

func isBasic(r rune) bool {
	return r < 32 || r >= 127
}

func IsNumber(s string) bool {
	return -1 == strings.IndexFunc(s, func(r rune) bool {
		return r < '0' || r > '9'
	})
}

func StripCtlAndExtFromUnicode(s string) string {
	t := transform.Chain(norm.NFKD, transform.RemoveFunc(isBasic))
	s, _, _ = transform.String(t, s)
	return s
}

func NormalizeWhitespace(s string) string {
	return reWhitespace.ReplaceAllString(s, " ")
}

func NormalizeText(s string) string {
	norm := StripCtlAndExtFromUnicode(s)
	norm = strings.ToLower(norm)
	norm = strings.Replace(norm, `\n`, "\n", -1)
	norm = reSpaceBoundary.ReplaceAllString(norm, " ")
	norm = reNoSpaceBoundary.ReplaceAllString(norm, "")
	norm = reWhitespace.ReplaceAllString(norm, " ")
	norm = strings.TrimSpace(norm)
	return norm
}

func SplitWhitespace(s string) []string {
	return reWhitespace.Split(strings.TrimSpace(s), -1)
}

func SplitSlash(s string) []string {
	return reSlash.Split(strings.TrimSpace(s), -1)
}

// CustomizeSlash generates slash-space variations for variable and unit name matching.
func CustomizeSlash(s string) []string {
	values := SplitSlash(s)
	modified := []string{
		strings.Join(values, "/"),
		strings.Join(values, " / "),
		strings.Join(values, "/ "),
		strings.Join(values, " /"),
	}
	return modified
}

func SplitSentence(s string) []string {
	return reSentenceBoundary.Split(s, -1)
}

func ToName(s string) string {
	s = strings.ToLower(s)
	s = strings.Replace(s, ",", " ", -1)
	s = strings.Replace(s, "/", " ", -1)
	s = reWhitespace.ReplaceAllString(s, "_")
	return s
}

// IsYesNo returns true if v contains both 'yes' and 'no' strings.
func IsYesNo(v []string) bool {
	if len(v) != 2 {
		return false
	}
	slice.TrimSpace(v)
	sort.Sort(sort.Reverse(sort.StringSlice(v)))
	return v[0] == "yes" && v[1] == "no"
}

// Join joins the array elements using the separators sep1 and sep2.
// For example: ["0", "1", "2"], sep1=", ", sep2=" or " -> "0, 1 or 2"
func Join(v []string, sep1, sep2 string) string {
	switch len(v) {
	case 0:
		return ""
	case 1:
		return v[0]
	case 2:
		return v[0] + sep2 + v[1]
	default:
		s := v[0]
		for i := 1; i < len(v)-1; i++ {
			s += sep1 + v[i]
		}
		return s + sep2 + v[len(v)-1]
	}
}

// LetterPrefix extracts a letter prefix from a string.
func LetterPrefix(s string) string {
	for i, r := range s {
		if !unicode.IsLetter(r) {
			return s[:i]
		}
	}
	return s
}
