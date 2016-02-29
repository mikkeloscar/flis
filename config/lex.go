// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// based on the lexer from: src/pkg/text/template/parse/lex.go (golang source)

package config

import (
	"fmt"
	"strings"
	"unicode"
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

func (i item) String() string {
	switch {
	case i.typ == itemEOF:
		return "EOF"
	case i.typ == itemError:
		return i.val
	case i.typ > itemKeyword:
		return fmt.Sprintf("<%s>", i.val)
	case len(i.val) > 10:
		return fmt.Sprintf("%.10q...", i.val)
	}
	return fmt.Sprintf("%q", i.val)
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
	itemNewline          // newline, for the most part, indicates end of command.
	itemBackslash        // backslash at the end of a line makes the command continue to the next line.
	itemSemicolon
	itemComma
	itemString   // quoted string (includes quotes)
	itemText     // plain text
	itemVariable // variable starting with '$', such as '$1' or '$hello'
	itemHexColor // hex color code
	itemCriteria
	// Keywords appear after all the rest.
	itemKeyword // used only to delimit the keywords
	itemFont
	itemYes
	itemNo
	itemBindsym
	itemBindcode
	itemFloatingModifier
	itemSet
	itemAssign
	itemExec
	itemWorkspace
	itemNext
	itemPrev
	itemNextOnOutput
	itemPrevOnOutput
	itemBackAndForth
	itemNumber
	itemWorkspaceAutoBackAndForth
	itemBar
	itemStatusCommand
	itemBarCommand
	itemModifier
	itemBarID
	itemPosition
	itemOutput
	itemColors
	itemBackground
	itemStatusLine
	itemSeparator
	itemFocusedWorkspace
	itemActiveWorkspace
	itemInactiveWorkspace
	itemUrgentWorkspace
	itemBindingMode
	itemWorkspaceButtons
	itemStripWorkspaceNumbers
	itemBindingModeIndicator
	itemMove
	itemContainer
	itemWindow
	itemTo
	itemKill
	itemSplit
	itemVertical
	itemHorizontal
	itemLayout
	itemDefault
	itemTabbed
	itemStacking
	itemSplitv
	itemSplith
	itemToggle
	itemAll
	itemFullscreen
	itemFloating
	itemFocus
	itemLeft
	itemRight
	itemDown
	itemUp
	itemCenter
	itemParent
	itemChild
	itemTiling
	itemModeToggle
	itemMouse
	itemAbsolute
	itemSticky
	itemEnable
	itemDisable
	itemRename
	itemResize
	itemGrow
	itemShrink
	itemMode
	itemHeight
	itemWidth
	itemRestart
	itemReload
	itemExit
	itemScratchpad
	itemShow
)

var key = map[string]itemType{
	"font":                          itemFont,
	"yes":                           itemYes,
	"no":                            itemNo,
	"bindsym":                       itemBindsym,
	"bindcode":                      itemBindcode,
	"floating_modifier":             itemFloatingModifier,
	"set":                           itemSet,
	"assign":                        itemAssign,
	"exec":                          itemExec,
	"workspace":                     itemWorkspace,
	"next":                          itemNext,
	"prev":                          itemPrev,
	"next_on_output":                itemNextOnOutput,
	"prev_on_output":                itemPrevOnOutput,
	"back_and_forth":                itemBackAndForth,
	"number":                        itemNumber,
	"workspace_auto_back_and_forth": itemWorkspaceAutoBackAndForth,
	"bar":                     itemBar,
	"status_command":          itemStatusCommand,
	"bar_command":             itemBarCommand,
	"modifier":                itemModifier,
	"id":                      itemBarID,
	"position":                itemPosition,
	"output":                  itemOutput,
	"colors":                  itemColors,
	"background":              itemBackground,
	"statusline":              itemStatusLine,
	"separator":               itemSeparator,
	"focused_workspace":       itemFocusedWorkspace,
	"active_workspace":        itemActiveWorkspace,
	"inactive_workspace":      itemInactiveWorkspace,
	"urgent_workspace":        itemUrgentWorkspace,
	"binding_mode":            itemBindingMode,
	"workspace_buttons":       itemWorkspaceButtons,
	"strip_workspace_numbers": itemStripWorkspaceNumbers,
	"binding_mode_indicator":  itemBindingModeIndicator,
	// commands
	"move":        itemMove,
	"container":   itemContainer,
	"window":      itemWindow,
	"to":          itemTo,
	"kill":        itemKill,
	"split":       itemSplit,
	"vertical":    itemVertical,
	"horizontal":  itemHorizontal,
	"layout":      itemLayout,
	"default":     itemDefault,
	"tabbed":      itemTabbed,
	"stacking":    itemStacking,
	"splitv":      itemSplitv,
	"splith":      itemSplith,
	"toggle":      itemToggle,
	"all":         itemAll,
	"fullscreen":  itemFullscreen,
	"floating":    itemFloating,
	"focus":       itemFocus,
	"left":        itemLeft,
	"right":       itemRight,
	"down":        itemDown,
	"up":          itemUp,
	"center":      itemCenter,
	"parent":      itemParent,
	"child":       itemChild,
	"tiling":      itemTiling,
	"mode_toggle": itemModeToggle,
	"mouse":       itemMouse,
	"absolute":    itemAbsolute,
	"sticky":      itemSticky,
	"enable":      itemEnable,
	"disable":     itemDisable,
	"rename":      itemRename,
	"resize":      itemResize,
	"grow":        itemGrow,
	"shrink":      itemShrink,
	"mode":        itemMode,
	"height":      itemHeight,
	"width":       itemWidth,
	"restart":     itemRestart,
	"reload":      itemReload,
	"exit":        itemExit,
	"scratchpad":  itemScratchpad,
	"show":        itemShow,
}

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

// accept consumes the next rune if it's from the valid set.
func (l *lexer) accept(valid string) bool {
	if strings.IndexRune(valid, l.next()) >= 0 {
		return true
	}
	l.backup()
	return false
}

// acceptRun consumes a run of runes from the valid set.
func (l *lexer) acceptRun(valid string) {
	for strings.IndexRune(valid, l.next()) >= 0 {
	}
	l.backup()
}

// lineNumber reports which line we're on, based on the position of
// the previous item returned by nextItem. Doing it this way
// means we don't have to worry about peek double counting.
func (l *lexer) lineNumber() int {
	return 1 + strings.Count(l.input[:l.lastPos], "\n")
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
		case '\n':
			l.emit(itemNewline)
		case ';':
			l.emit(itemSemicolon)
		case ',':
			l.emit(itemComma)
		case ' ':
			l.ignore()
		case '$':
			return lexVariable
		case '#':
			return lexComment
		case '"':
			return lexQuote
		case '[':
			return lexCriteria
		case '{':
			l.emit(itemLeftCurlBracket)
		case '}':
			l.emit(itemRightCurlBracket)
		case '\\':
			if l.peek() == '\n' {
				l.emit(itemBackslash)
			}
		default:
			return lexText
		}
	}
}

func lexText(l *lexer) stateFn {
	for {
		switch l.next() {
		case ' ', '$', '#', '"', '[', ',', ';', '\\', '\n', eof, '{', '}':
			l.backup()
			keyword := l.input[l.start:l.pos]

			if v, ok := key[keyword]; ok {
				l.emit(v)
				return lexConfig
			}
			l.emit(itemText)
			return lexConfig
		default:
			// absorb
		}
	}
}

func lexComment(l *lexer) stateFn {
	colorCheck := false
	for {
		switch l.next() {
		case '\n':
			if !colorCheck && isHexColor(l.input[l.start:l.pos-1]) {
				l.backup()
				l.emit(itemHexColor)
				return lexConfig
			}
			l.ignore()
			return lexConfig
		case ' ':
			if !colorCheck && isHexColor(l.input[l.start:l.pos-1]) {
				l.backup()
				l.emit(itemHexColor)
				return lexConfig
			}
			colorCheck = true
		case eof:
			l.emit(itemEOF)
			return nil
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
	l.emit(itemString)
	return lexConfig
}

// lexCriteria scans a criteria definition.
func lexCriteria(l *lexer) stateFn {
Loop:
	for {
		switch l.next() {
		case eof, '\n':
			return l.errorf("unterminated criteria")
		case ']':
			break Loop
		}
	}
	l.emit(itemCriteria)
	return lexConfig
}

func lexVariable(l *lexer) stateFn {
	for {
		switch r := l.next(); {
		case isAlphaNumericUnderscore(r):
			// absorb
		case r == ' ' || r == '\n' || r == eof:
			l.backup()
			if l.pos-l.start < 2 {
				return l.errorf("$ can't stand by itself")
			}
			l.emit(itemVariable)
			return lexConfig
		default:
			pattern := l.input[l.start:l.pos]
			return l.errorf("invalid variable: %s", pattern)
		}
	}

}

func isHexColor(color string) bool {
	if len(color) != 6 && len(color) != 8 {
		return false
	}

	for _, c := range color {
		if !isHexChar(c) {
			return false
		}
	}

	return true
}

func isHexChar(r rune) bool {
	return unicode.IsDigit(r) || ('a' <= r && r <= 'f') || ('A' <= r && r <= 'F')
}

func isAlphaNumericUnderscore(r rune) bool {
	return r == '_' || unicode.IsDigit(r) || (unicode.IsLetter(r) && unicode.IsLower(r))
}
