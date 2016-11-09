// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// based on the lexer from: src/pkg/text/template/parse/lex.go (golang source)

package config

import (
	"fmt"
	"unicode/utf8"
)

// pos is a position in input being scanned.
type pos int

// item represents a token or text string returned from the scanner.
type item struct {
	typ itemType // The type of this item.
	pos pos      // The starting position, in bytes, of this item in the input string.
	val string   // The value of this item.
}

// itemType identifies the type of lex items.
type itemType int

const (
	itemError itemType = iota // error occurred; value is text of error
	itemEOF
	itemLeftCurlBracket  // '{'
	itemRightCurlBracket // '}'
	itemLeftBracket      // '['
	itemRightBracket     // ']'
	itemSemicolon
	itemComma
	itemString // string
	// itemText     // plain text
	itemHexColor // hex color code
	itemCriteria
)

const eof = -1

// stateFn represents the state of the scanner as a function that returns the next state.
type stateFn func(*lexer) stateFn

// lexer holds the state of the scanner.
type lexer struct {
	input      string    // the string being scanned
	state      stateFn   // the next lexing function to enter
	pos        pos       // current position in the input
	start      pos       // start position of this item
	width      pos       // width of last rune read from input
	lastPos    pos       // position of most recent item returned by nextItem
	items      chan item // channel of scanned items
	parenDepth int       // nesting depth of ( ) exprs
}

// next returns the next rune in the input.
func (l *lexer) next() rune {
	if int(l.pos) >= len(l.input) {
		l.width = 0
		return eof
	}
	r, w := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = pos(w)
	l.pos += l.width
	return r
}

// peek returns but does not consume the next rune in the input.
func (l *lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

// backup steps back one rune. Can only be called once per call of next.
func (l *lexer) backup() {
	l.pos -= l.width
}

// emit passes an item back to the client.
func (l *lexer) emit(t itemType) {
	l.items <- item{t, l.start, l.input[l.start:l.pos]}
	l.start = l.pos
}

// ignore skips over the pending input before this point.
func (l *lexer) ignore() {
	l.start = l.pos
}

// errorf returns an error token and terminates the scan by passing
// back a nil pointer that will be the next state, terminating l.nextItem.
func (l *lexer) errorf(format string, args ...interface{}) stateFn {
	l.items <- item{itemError, l.start, fmt.Sprintf(format, args...)}
	return nil
}

// nextItem returns the next item from the input.
// Called by the parser, not in the lexing goroutine.
func (l *lexer) nextItem() item {
	item := <-l.items
	l.lastPos = item.pos
	return item
}

// drain drains the output so the lexing goroutine will exit.
// Called by the parser, not in the lexing goroutine.
func (l *lexer) drain() {
	for range l.items {
	}
}

// lex creates a new scanner for the input string.
func lex(input string) *lexer {
	l := &lexer{
		input: input,
		items: make(chan item),
	}
	go l.run()
	return l
}

// run runs the state machine for the lexer.
func (l *lexer) run() {
	for l.state = lexConfig; l.state != nil; {
		l.state = l.state(l)
	}
	close(l.items)
}

// state functions

func lexConfig(l *lexer) stateFn {
	for {
		switch l.next() {
		case eof:
			l.emit(itemEOF)
			return nil
		case ';':
			l.emit(itemSemicolon)
		case ',':
			l.emit(itemComma)
		case ' ', '\t':
			l.ignore()
		case '"':
			return lexQuote
		case '[':
			return lexCriteria
		case '\\':
			switch l.peek() {
			case ',', ';', '\\', ' ', '\t':
				l.next()
				return lexText
			default:
				return l.errorf("invalid escape char: '\\%x'", l.next())
			}
		default:
			return lexText
		}
	}
}

func lexText(l *lexer) stateFn {
	for {
		switch l.next() {
		case ' ', '\t', eof:
			l.backup()
			l.emit(itemString)
			return lexConfig
		default:
			// absorb
		}
	}
}

// lexQuote scans a quoted string.
func lexQuote(l *lexer) stateFn {
Loop:
	for {
		switch l.next() {
		case '\\':
			if r := l.next(); r != eof && r != '\n' {
				break
			}
			fallthrough
		case eof, '\n':
			return l.errorf("unterminated quoted string")
		case '"':
			break Loop
		}
	}
	l.backup()
	l.start++
	l.emit(itemString)
	l.next()
	l.ignore()
	return lexConfig
}

// lexCriteria scans a criteria definition.
func lexCriteria(l *lexer) stateFn {
Loop:
	for {
		switch l.next() {
		case eof:
			return l.errorf("unterminated criteria")
		case ']':
			break Loop
		}
	}
	l.emit(itemCriteria)
	return lexConfig
}
