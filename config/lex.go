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
	itemSemicolon
	itemComma
	itemString   // quoted string (includes quotes)
	itemText     // plain text
	itemHexColor // hex color code
	itemCriteria
	// Keywords appear after all the rest.
	itemKeyword // used only to delimit the keywords
	// config options
	itemFont
	itemBindsym
	itemBindcode
	itemFloatingModifier
	itemFloatingMinimumSize
	itemFloatingMaximumSize
	itemDefaultOrientation
	itemWorkspaceLayout
	itemNewWindow
	itemNewFloat
	itemHideEdgeBorders
	itemForWindow
	itemNoFocus
	itemSet
	itemAssign
	itemFocusFollowsMouse
	itemMouseWarping
	itemWorkspaceAutoBackAndForth
	itemBar
	// bar subcommand
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
	// command
	itemFullscreen
	itemFloating
	itemExec
	itemKill
	itemRestart
	itemReload
	itemExit
	itemWorkspace
	itemMove
	itemSplit
	itemSplitv
	itemSplith
	itemSplitt
	itemLayout
	itemFocus
	itemSticky
	itemRename
	itemResize
	itemMode
	itemBorder
	itemScratchpad
	// argument
	itemYes
	itemNo
	itemNext
	itemPrev
	itemNextOnOutput
	itemPrevOnOutput
	itemBackAndForth
	itemNumber
	itemContainer
	itemWindow
	itemTo
	itemVertical
	itemHorizontal
	itemDefault
	itemTabbed
	itemStacking
	itemToggle
	itemAll
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
	itemEnable
	itemDisable
	itemGrow
	itemShrink
	itemHeight
	itemWidth
	itemShow
)

var key = map[string]itemType{
	// config options
	"font":                          itemFont,
	"bindsym":                       itemBindsym,
	"bindcode":                      itemBindcode,
	"floating_modifier":             itemFloatingModifier,
	"floating_minimum_size":         itemFloatingMinimumSize,
	"floating_maximum_size":         itemFloatingMaximumSize,
	"default_orientation":           itemDefaultOrientation,
	"workspace_layout":              itemWorkspaceLayout,
	"new_window":                    itemNewWindow,
	"new_float":                     itemNewFloat,
	"hide_edge_borders":             itemHideEdgeBorders,
	"for_window":                    itemForWindow,
	"no_focus":                      itemNoFocus,
	"set":                           itemSet,
	"assign":                        itemAssign,
	"focus_follows_mouse":           itemFocusFollowsMouse,
	"mouse_warping":                 itemMouseWarping,
	"workspace_auto_back_and_forth": itemWorkspaceAutoBackAndForth,
	"bar": itemBar,
	// bar subcommands
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
	"fullscreen": itemFullscreen,
	"floating":   itemFloating,
	"exec":       itemExec,
	"kill":       itemKill,
	"restart":    itemRestart,
	"reload":     itemReload,
	"exit":       itemExit,
	"workspace":  itemWorkspace,
	"move":       itemMove,
	"split":      itemSplit,
	"splitv":     itemSplitv,
	"splith":     itemSplith,
	"splitt":     itemSplitt,
	"layout":     itemLayout,
	"focus":      itemFocus,
	"sticky":     itemSticky,
	"rename":     itemRename,
	"resize":     itemResize,
	"mode":       itemMode,
	"border":     itemBorder,
	"scratchpad": itemScratchpad,
	// arguments
	"yes":            itemYes,
	"no":             itemNo,
	"next":           itemNext,
	"prev":           itemPrev,
	"next_on_output": itemNextOnOutput,
	"prev_on_output": itemPrevOnOutput,
	"back_and_forth": itemBackAndForth,
	"number":         itemNumber,
	"container":      itemContainer,
	"window":         itemWindow,
	"to":             itemTo,
	"vertical":       itemVertical,
	"horizontal":     itemHorizontal,
	"default":        itemDefault,
	"tabbed":         itemTabbed,
	"stacking":       itemStacking,
	"toggle":         itemToggle,
	"all":            itemAll,
	"left":           itemLeft,
	"right":          itemRight,
	"down":           itemDown,
	"up":             itemUp,
	"center":         itemCenter,
	"parent":         itemParent,
	"child":          itemChild,
	"tiling":         itemTiling,
	"mode_toggle":    itemModeToggle,
	"mouse":          itemMouse,
	"absolute":       itemAbsolute,
	"enable":         itemEnable,
	"disable":        itemDisable,
	"grow":           itemGrow,
	"shrink":         itemShrink,
	"height":         itemHeight,
	"width":          itemWidth,
	"show":           itemShow,
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
		case ' ', '\t':
			l.ignore()
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
				l.next()
				l.ignore()
			}
		default:
			return lexText
		}
	}
}

func lexText(l *lexer) stateFn {
	for {
		switch l.next() {
		case ' ', '\t', '$', '#', '"', '[', ',', ';', '\\', '\n', eof, '{', '}':
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
		case '\n', eof:
			l.backup()
			if !colorCheck && isHexColor(l.input[l.start:l.pos]) {
				l.emit(itemHexColor)
				return lexConfig
			}
			l.ignore()
			return lexConfig
		case ' ', '\t':
			if !colorCheck && isHexColor(l.input[l.start:l.pos-1]) {
				l.backup()
				l.emit(itemHexColor)
				return lexConfig
			}
			colorCheck = true
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

func isHexColor(color string) bool {
	if len(color) != 7 && len(color) != 9 {
		return false
	}

	if color[0] != '#' {
		return false
	}

	for _, c := range color[1:] {
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
