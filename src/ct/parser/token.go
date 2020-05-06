// Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.

package parser

import (
	"fmt"
)

// end-of-string marker
const eos = -1

// tokenType defines the type of lexer tokens.
type tokenType int

const (
	tokenError             tokenType = iota // error; value is text of error
	tokenEOF                                // end of string token
	tokenChar                               // printable ASCII character; escape hatch for unspecified tokens.
	tokenSpace                              // spaces
	tokenIdentifier                         // alphanumeric identifier not starting with '.'
	tokenNumber                             // simple number, including imaginary
	tokenUnit                               // unit token
	tokenLeftParenthesis                    // '('
	tokenRightParenthesis                   // ')'
	tokenDash                               // '-'
	tokenSlash                              // '/'
	tokenPunctuation                        // punctuations
	tokenKeyword                            // keywords specified after this
	tokenConjunction                        // conjunction: 'and', 'or'
	tokenNegation                           // negation: 'no', 'not'
	tokenComparison                         // comparison token
	tokenLessComparison                     // less than comparison token
	tokenGreaterComparison                  // greater than comparison token
)

// Pos is the rune position of the token in the string.
type Pos int

// Token defines a token or text string that is returned from the lexer.
type Token struct {
	typ tokenType // The type of token.
	pos Pos       // The starting position, in bytes, of this Token in the input string.
	val string    // The value of token.
}

// NewToken creates a new token.
func NewToken(typ tokenType, pos Pos, val string) *Token {
	return &Token{typ: typ, pos: pos, val: val}
}

// String returns a string representation of the token.
func (t *Token) String() string {
	return fmt.Sprintf("{id:%d,value:%q}", t.typ, t.val)
}

// Tokens defines a slice of tokens.
type Tokens []*Token

// NewTokens creates a slice of tokens.
func NewTokens() Tokens {
	return make(Tokens, 0)
}
