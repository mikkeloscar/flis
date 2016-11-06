// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// based on the lexer tests from: src/pkg/text/template/parse/lex_test.go (golang source)

package commands

import (
	"fmt"
	"testing"
)

type lexTest struct {
	name  string
	input string
	items []item
}

var tEOF = item{itemEOF, 0, ""}

var lexTests = []lexTest{
	{"empty", "", []item{tEOF}},
	{"text", `random    text`, []item{
		{itemString, 0, "random"},
		{itemString, 0, "text"},
		tEOF,
	}},
	{"quoted string", `"this is a string"`, []item{
		{itemString, 0, `this is a string`},
		tEOF,
	}},
	{"unterminated quoted string", `"this is a string`, []item{
		{itemError, 0, "unterminated quoted string"},
	}},
	{"valid escaped quoted string", `"this is a \\string"`, []item{
		{itemString, 0, `this is a \\string`},
		tEOF,
	}},
	{"invalid escaped quoted string", `"this is a string\"`, []item{
		{itemError, 0, "unterminated quoted string"},
	}},
	{"invalid escaped quoted string2", `"this is a string\`, []item{
		{itemError, 0, "unterminated quoted string"},
	}},
	{"comma", `,;`, []item{
		{itemComma, 0, ","},
		{itemSemicolon, 0, ";"},
		tEOF,
	}},
	{"full binding", `$mod+Return exec $term`, []item{
		{itemString, 0, "$mod+Return"},
		{itemString, 0, "exec"},
		{itemString, 0, "$term"},
		tEOF,
	}},
	{"combined command", `exec dmenu_run -fn "Terminus-8" -nb "#000" -nf "#fff"`, []item{
		{itemString, 0, "exec"},
		{itemString, 0, "dmenu_run"},
		{itemString, 0, "-fn"},
		{itemString, 0, "Terminus-8"},
		{itemString, 0, "-nb"},
		{itemString, 0, "#000"},
		{itemString, 0, "-nf"},
		{itemString, 0, "#fff"},
		tEOF,
	}},
	{"escaped quote", `escaped\"`, []item{
		{itemString, 0, `escaped\"`},
		tEOF,
	}},
	{"criteria", `[class="firefox"] exec terminal`, []item{
		{itemCriteria, 0, `[class="firefox"]`},
		{itemString, 0, `exec`},
		{itemString, 0, `terminal`},
		tEOF,
	}},
	{"unterminated criteria", `[class="firefox"`, []item{
		{itemError, 0, "unterminated criteria"},
	}},
	{"escaped char", `\, \; \\ \  \	`, []item{
		{itemString, 0, `\,`},
		{itemString, 0, `\;`},
		{itemString, 0, `\\`},
		{itemString, 0, `\ `},
		{itemString, 0, "\\\t"},
		tEOF,
	}},
	{"invalid escaped char", `\x`, []item{
		{itemError, 0, fmt.Sprintf("invalid escape char: '\\%x'", "x")},
	}},

	// {"color", `#ffffff`, []item{
	// 	{itemHexColor, 0, "#ffffff"},
	// 	tEOF,
	// }},
	// {"color alpha", `#ffffff00`, []item{
	// 	{itemHexColor, 0, "#ffffff00"},
	// 	tEOF,
	// }},
	// {"color comment", "\t\t#ffffff # this is a color", []item{
	// 	{itemHexColor, 0, "#ffffff"},
	// 	tEOF,
	// }},
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

// equal checks if two items are equal.
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
