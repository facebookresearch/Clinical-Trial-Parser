// Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.

package parser

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/facebookresearch/Clinical-Trial-Parser/src/common/util/text"
)

const (
	emitSpace = false

	comparisonChars  = "<≤>≥="
	punctuationChars = ".,;"
	identifierChars  = "%^/-"
)

var key = map[string]tokenType{
	"and":     tokenConjunction,
	"or":      tokenConjunction,
	"and/or":  tokenConjunction,
	"but":     tokenConjunction,
	"no":      tokenNegation,
	"not":     tokenNegation,
	"less":    tokenLessComparison,
	"below":   tokenLessComparison,
	"under":   tokenLessComparison,
	"younger": tokenLessComparison,
	"above":   tokenGreaterComparison,
	"greater": tokenGreaterComparison,
	"higher":  tokenGreaterComparison,
	"more":    tokenGreaterComparison,
	"over":    tokenGreaterComparison,
	"longer":  tokenGreaterComparison,
	"older":   tokenGreaterComparison,
	"between": tokenComparison,
	"at":      tokenComparison,
	"least":   tokenComparison,
	"than":    tokenComparison,
}

// stateFn represents the state of the lexer as a function that returns the next state.
type stateFn func(*Lexer) stateFn

// Lexer defines the struct that converts a string to tokens.
// Rob Pike's presentation has been a great inspiration for this design.
// https://www.youtube.com/watch?v=HxaD_trXwRE
type Lexer struct {
	input      string      // the string being scanned
	pos        Pos         // current position in the input
	start      Pos         // start position of this Token
	width      Pos         // width of last rune read from input
	tokens     chan *Token // channel of scanned tokens
	parenDepth int         // nesting depth of ( )
}

// NewLexer creates a new lexer for the input string.
func NewLexer(input string) *Lexer {
	l := &Lexer{
		input:  input,
		tokens: make(chan *Token),
	}
	go l.run()
	return l
}

// NextToken returns the next token from the input string.
func (l *Lexer) NextToken() *Token {
	return <-l.tokens
}

// Drain drains the output to the slice of tokens.
func (l *Lexer) Drain() Tokens {
	tokens := NewTokens()
	for token := l.NextToken(); token.typ != tokenEOF; token = l.NextToken() {
		tokens = append(tokens, token)
	}
	return tokens
}

// run runs the state machine for the lexer.
func (l *Lexer) run() {
	for state := lexAction; state != nil; {
		state = state(l)
	}
	close(l.tokens)
}

// next returns the next rune in the input string.
func (l *Lexer) next() rune {
	if int(l.pos) >= len(l.input) {
		l.width = 0
		return eos
	}
	r, w := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = Pos(w)
	l.pos += l.width
	return r
}

// peek returns but does not consume the next rune in the input string.
func (l *Lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

// current returns the current rune in the input string.
func (l *Lexer) current() rune {
	if l.pos < l.width {
		return ' '
	}
	r, _ := utf8.DecodeRuneInString(l.input[l.pos-l.width:])
	return r
}

// backup steps back one rune. Can only be called once per call of next.
func (l *Lexer) backup() {
	l.pos -= l.width
}

// emit passes a token back to the client.
func (l *Lexer) emit(t tokenType) {
	l.tokens <- NewToken(t, l.start, l.input[l.start:l.pos])
	l.start = l.pos
}

// swallow skips over the pending input before this point.
func (l *Lexer) swallow() {
	l.start = l.pos
}

// accept consumes the next rune if it's from the valid set.
func (l *Lexer) accept(valid string) bool {
	if strings.ContainsRune(valid, l.next()) {
		return true
	}
	l.backup()
	return false
}

// acceptRun consumes a run of runes from the valid set.
func (l *Lexer) acceptRun(valid string) {
	for strings.ContainsRune(valid, l.next()) {
	}
	l.backup()
}

// errorf returns an error token and terminates the scan by passing
// back a nil pointer that will be the next state, terminating l.nexttoken.
func (l *Lexer) errorf(format string, args ...interface{}) stateFn {
	l.tokens <- NewToken(tokenError, l.start, fmt.Sprintf(format, args...))
	return nil
}

// State functions:

// lexAction scans the elements from the input string.
func lexAction(l *Lexer) stateFn {
	q := l.current()
	switch r := l.next(); {
	case r == eos:
		l.emit(tokenEOF)
		return nil
	case isSpaceChar(r):
		return lexSpace
	case unicode.IsLetter(r):
		l.backup()
		return lexIdentifier
	case isComparisonChar(r):
		return lexComparison
	case r == '%':
		l.emit(tokenUnit)
	case r == '-':
		if q == ' ' && unicode.IsDigit(l.peek()) {
			l.backup()
			return lexNumber
		}
		l.emit(tokenDash)
	case r == '+' || ('0' <= r && r <= '9'):
		l.backup()
		return lexNumber
	case r == '/':
		l.emit(tokenSlash)
	case r == '(':
		l.emit(tokenLeftParenthesis)
		l.parenDepth++
	case r == ')':
		l.emit(tokenRightParenthesis)
		l.parenDepth--
		if l.parenDepth < 0 {
			l.swallow()
		}
	case isPunctuationChar(r) && r != ',':
		l.emit(tokenPunctuation)
	default:
		l.emit(tokenChar)
	}
	return lexAction
}

// lexSpace scans a run of space characters.
// One space has already been seen.
func lexSpace(l *Lexer) stateFn {
	for isSpaceChar(l.peek()) {
		l.next()
	}
	switch {
	case emitSpace:
		l.emit(tokenSpace)
	default:
		l.swallow()
	}
	return lexAction
}

// lexIdentifier scans an alphanumeric word.
func lexIdentifier(l *Lexer) stateFn {
Loop:
	for {
		switch r := l.next(); {
		case isIdentifierChar(r):
			// absorb.
		default:
			l.backup()
			word := l.input[l.start:l.pos]
			switch {
			case key[word] > tokenKeyword:
				l.emit(key[word])
			case text.IsRomanNumeral(word):
				l.emit(tokenNumber)
			default:
				l.emit(tokenIdentifier)
			}
			break Loop
		}
	}
	return lexAction
}

// lexNumber scans a number, which can be int, decimal or scientific.
func lexNumber(l *Lexer) stateFn {
	if !l.scanNumber() {
		return l.errorf("bad number syntax: %q", l.input[l.start:l.pos])
	}
	l.emit(tokenNumber)
	return lexAction
}

func (l *Lexer) scanNumber() bool {
	l.accept("+-")
	l.acceptRun("0123456789,.")
	if l.accept("eE") {
		l.accept("+-")
		l.acceptRun("0123456789")
	} else if l.current() != ',' {
		pos := l.pos
		if !l.scanScientificMultiplier() {
			l.pos = pos
		}
	}
	r, w := utf8.DecodeRuneInString(l.input[l.pos-1:])
	if isPunctuationChar(r) {
		l.pos -= Pos(w)
	}
	return true
}

// scanScientificMultiplier scans the scientific multiplier ('x 10^exp') of the scientific number.
// The multiply symbol (x|⨯) is optional. Returns true if a valid multiplier is found.
func (l *Lexer) scanScientificMultiplier() bool {
	l.acceptRun(" ")
	r := l.peek()
	if r == 'x' || r == '×' {
		l.next()
		l.acceptRun(" ")
		r = l.peek()
	}
	if !unicode.IsDigit(r) {
		return false
	}
	pos := l.pos
	l.acceptRun("0123456789")
	if l.accept("^") || l.accept("e") {
		if l.accept("0123456789") {
			l.acceptRun("0123456789")
		} else {
			return false
		}
	}
	// The multiplier (10^exp) should contain more than 2 runes.
	return l.pos-pos > 2
}

// lexComparison scans a run of comparison characters.
// One comparison has already been seen.
func lexComparison(l *Lexer) stateFn {
	l.accept("/")
	l.accept(comparisonChars)
	l.emit(tokenComparison)
	return lexAction
}

// isSpaceChar reports whether r is a space character.
func isSpaceChar(r rune) bool {
	return r == ' ' || r == '\t' || r == '\r' || r == '\n'
}

// isIdentifierChar reports whether r is a valid identifier rune.
func isIdentifierChar(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsNumber(r) || strings.ContainsRune(identifierChars, r)
}

// isComparisonChar reports whether r is a valid comparison.
func isComparisonChar(r rune) bool {
	return strings.ContainsRune(comparisonChars, r)
}

// isPunctuationChar reports whether r is a valid punctuation.
func isPunctuationChar(r rune) bool {
	return strings.ContainsRune(punctuationChars, r)
}
