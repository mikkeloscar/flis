// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// based on the lexer tests from: src/pkg/text/template/parse/lex_test.go (golang source)

package config

import "testing"

type lexTest struct {
	name  string
	input string
	items []item
}

var (
	tEOF = item{itemEOF, 0, ""}
	tNL  = item{itemNewline, 0, "\n"}
)

var lexTests = []lexTest{
	{"empty", "", []item{tEOF}},
	{"text", `random text`, []item{
		{itemText, 0, "random"},
		{itemText, 0, "text"},
		tEOF,
	}},
	{"assignment", `set $mod value`, []item{
		{itemSet, 0, "set"},
		{itemVariable, 0, "$mod"},
		{itemText, 0, "value"},
		tEOF,
	}},
	{"string", `"this is a string"`, []item{
		{itemString, 0, `"this is a string"`},
		tEOF,
	}},
	{"comma", `,;`, []item{
		{itemComma, 0, ","},
		{itemSemicolon, 0, ";"},
		tEOF,
	}},
	{"bar", `bar {`, []item{
		{itemBar, 0, "bar"},
		{itemLeftCurlBracket, 0, "{"},
		tEOF,
	}},
	{"space", "  \t", []item{tEOF}},
	{"bindsym", `bindsym $mod+q kill`, []item{
		{itemBindsym, 0, "bindsym"},
		{itemVariable, 0, "$mod"},
		{itemText, 0, "+q"},
		{itemKill, 0, "kill"},
		tEOF,
	}},
	{"bindsym no vars", `bindsym Mod1+q kill`, []item{
		{itemBindsym, 0, "bindsym"},
		{itemText, 0, "Mod1+q"},
		{itemKill, 0, "kill"},
		tEOF,
	}},
	{"comment", `# this is a comment`, []item{tEOF}},
	{"comment with newline", "# this is a comment\n", []item{tNL, tEOF}},
	{"color", `#ffffff`, []item{
		{itemHexColor, 0, "#ffffff"},
		tEOF,
	}},
	{"color alpha", `#ffffff00`, []item{
		{itemHexColor, 0, "#ffffff00"},
		tEOF,
	}},
	{"color comment", "\t\t#ffffff # this is a color", []item{
		{itemHexColor, 0, "#ffffff"},
		tEOF,
	}},
	{"backslash newline", "\\\n", []item{
		tEOF,
	}},
	{"newline", "\n", []item{
		tNL,
		tEOF,
	}},
}

// collect gathers the emitted items into a slice.
func collect(t *lexTest) (items []item) {
	l := lex(t.input)
	for {
		item := l.nextItem()
		items = append(items, item)
		if item.typ == itemEOF || item.typ == itemError {
			break
		}
	}
	return
}

func equal(i1, i2 []item, checkPos bool) bool {
	if len(i1) != len(i2) {
		return false
	}
	for k := range i1 {
		if i1[k].typ != i2[k].typ {
			return false
		}
		if i1[k].val != i2[k].val {
			return false
		}
		if checkPos && i1[k].pos != i2[k].pos {
			return false
		}
	}
	return true
}

func TestLex(t *testing.T) {
	for _, test := range lexTests {
		items := collect(&test)
		if !equal(items, test.items, false) {
			t.Errorf("%s: got\n\t%+v\nexpected\n\t%v", test.name, items, test.items)
		}
	}
}
