package config

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/mikkeloscar/go-wlc"
	"github.com/mikkeloscar/go-xkbcommon"
)

var modifiers = map[string]wlc.ModifierBit{
	xkb.ModNameShift: wlc.BitModShift,
	xkb.ModNameCaps:  wlc.BitModCaps,
	xkb.ModNameCtrl:  wlc.BitModCtrl,
	"Ctrl":           wlc.BitModCtrl,
	xkb.ModNameAlt:   wlc.BitModAlt,
	"Alt":            wlc.BitModAlt,
	xkb.ModNameNum:   wlc.BitModMod2,
	"Mod3":           wlc.BitModMod3,
	xkb.ModNameLogo:  wlc.BitModLogo,
	"Mod5":           wlc.BitModMod5,
}

type parser struct {
	input    *string
	config   *Config
	cmdStart *item
	lexer    *lexer
}

// func ParseCmd(config *Config, input string) error {

// }

func ParseConfig(config *Config, input string) error {
	// TODO check if in parsing config mode or parse single commands mode
	var err error

	p := parser{
		input:  &input,
		config: config,
		lexer:  lex(input),
	}

Loop:
	for {
		token := p.lexer.nextItem()
		p.cmdStart = &token
		switch token.typ {
		case itemEOF:
			break Loop
		case itemNewline:
			break
		case itemSet:
			err = p.parseSet()
			if err != nil {
				// TODO error log
				fmt.Println(err)
			}
		case itemBindsym:
			err = p.parseBindsym()
			if err != nil {
				// TODO error log
				fmt.Println(err)
			}
		case itemError:
			// TODO log error
		default:
			// cmd, next, err = p.parseCmds(next, nil)
			// ignore
		}
	}

	return nil
}

func (p *parser) parseSet() error {
	variable := p.lexer.nextItem()

	if variable.typ != itemText {
		return p.errorf(variable, "Unexpected token")
	}

	if strings.HasPrefix(variable.val, "$") {
		for _, c := range variable.val[1:] {
			if isAlphaNumericUnderscore(c) {
				continue
			} else {
				return p.errorf(variable, "Invalid variable")
			}
		}
	} else {
		return p.errorf(variable, "Invalid variable, has to start with $")
	}

	values := make([]string, 0, 1)

	for {
		next := p.lexer.nextItem()
		switch next.typ {
		case itemText, itemString:
			values = append(values, next.val)
		default:
			return p.errorf(next, "Invalid value")
		}
	}

	if len(values) == 0 {
		return p.errorf(variable, "No value assigned to variable %s", variable.val)
	}

	combinedValues := strings.Join(values, " ")
	p.config.vars[variable.val] = combinedValues

	return nil
}

func (p *parser) parseBindsym() error {
	next := p.lexer.nextItem()
	if next.typ != itemText {
		return p.errorf(next, "Unexpected token, expected key combination")
	}

	keyStr := p.varReplacement(next.val)
	keys := strings.Split(keyStr, "+")

	binding := new(Binding)

	for _, key := range keys {
		if mod, ok := modifiers[key]; ok {
			binding.Modifiers |= uint32(mod)
			continue
		}

		sym := xkb.KeySymFromName(key, xkb.KeySymCaseInsensitive)
		if sym == xkb.KeyNoSymbol {
			return p.errorf(next, "Unknown key '%s'", key)
		}

		binding.Keys = append(binding.Keys, sym)
	}

	next = p.lexer.nextItem()

	// TODO: handle criteria
	var criteria *Criteria
	if next.typ == itemCriteria {
		criteria = nil
		next = p.lexer.nextItem()
	}

	// handle command
	if next.typ != itemNewline && next.typ != itemEOF {
		cmd, err := p.parseCmds(next, criteria)
		if err != nil {
			return err
		}
		binding.Command = cmd
	} else {
		return p.errorf(next, "Missing command definition")
	}

	replaced := p.config.AddBinding(binding)
	if replaced {
		// log. already exists, overwriting
	}

	return nil
}

func (p *parser) parseCmds(cmd item, criteria *Criteria) (*Command, error) {
	var err error
	var start *Command
	var curr *Command

	next := cmd

Loop:
	for {
		cmd := &Command{Criteria: criteria}
		switch next.typ {
		case itemFullscreen, itemFloating:
			cmd.Exec = Fullscreen
			if next.typ == itemFloating {
				cmd.Exec = Float
			}

			next = p.lexer.nextItem()
			var args []string
			switch next.typ {
			case itemEnable, itemDisable, itemToggle:
				args = []string{next.val}
			default:
				return nil, p.errorf(next, "Invalid command")
			}
			cmd.Args = args
			next = p.lexer.nextItem()
		case itemExec:
			var item *item
			cmd, item, err = p.parseExec(criteria)
			if err != nil {
				return nil, err
			}
			next = *item
		case itemKill:
			cmd.Exec = Kill
			next = p.lexer.nextItem()
		// case itemRestart:
		case itemReload:
			cmd.Exec = Reload
			next = p.lexer.nextItem()
		case itemExit:
			cmd.Exec = Exit
			next = p.lexer.nextItem()
		case itemWorkspace:
			next = p.lexer.nextItem()
			if next.typ != itemText {
				return nil, p.errorf(next, "Invalid token, expected workspace name/number")
			}

			cmd.Args = []string{p.varReplacement(next.val)}
			cmd.Exec = WorkspaceFn

			next = p.lexer.nextItem()

			if next.typ == itemOutput {
				next = p.lexer.nextItem()
				if next.typ == itemText {
					cmd.Args = append(cmd.Args, "output")
					cmd.Args = append(cmd.Args, p.varReplacement(next.val))
					next = p.lexer.nextItem()
					break
				}
				return nil, p.errorf(next, "Invalid output argument")
			}
		case itemMove:
			var item *item
			cmd, item, err = p.parseMove(criteria)
			if err != nil {
				return nil, err
			}
			next = *item
		case itemSplit:
			var split string
			next = p.lexer.nextItem()
			switch next.typ {
			case itemVertical, itemHorizontal, itemToggle:
				split = next.val
			case itemText:
				switch next.val {
				case "v":
					split = "vertical"
				case "h":
					split = "horizontal"
				case "t":
					split = "toggle"
				default:
					return nil, p.errorf(next, "Invalid split value")
				}
			default:
				return nil, p.errorf(next, "Unexpected token")
			}
			cmd.Args = []string{split}
			cmd.Exec = Split
			next = p.lexer.nextItem()
		case itemSplitv:
			cmd.Args = []string{"vertical"}
			cmd.Exec = Split
			next = p.lexer.nextItem()
		case itemSplith:
			cmd.Args = []string{"horizontal"}
			cmd.Exec = Split
			next = p.lexer.nextItem()
		case itemSplitt:
			cmd.Args = []string{"toggle"}
			cmd.Exec = Split
			next = p.lexer.nextItem()
		case itemLayout:
			next = p.lexer.nextItem()
			switch next.typ {
			case itemDefault, itemTabbed, itemStacking, itemSplitv, itemSplith:
				cmd.Args = []string{next.val}
				cmd.Exec = Layout
			case itemToggle:
				next = p.lexer.nextItem()
				cmd.Exec = LayoutToggle
				switch next.typ {
				case itemSplit, itemAll:
					cmd.Args = []string{next.val}
					next = p.lexer.nextItem()
				case itemNewline, itemEOF:
					break
				default:
					return nil, p.errorf(next, "Unexpected token")
				}
			default:
				return nil, p.errorf(next, "Unexpected token")
			}
			next = p.lexer.nextItem()
		case itemFocus:
			next = p.lexer.nextItem()
			switch next.typ {
			case itemLeft, itemRight, itemDown, itemUp:
				cmd.Args = []string{next.val}
				cmd.Exec = FocusDirection
			case itemParent, itemChild, itemFloating, itemTiling, itemModeToggle:
				cmd.Args = []string{next.val}
				cmd.Exec = FocusCriteria
			case itemOutput:
				cmd.Exec = FocusOutput
				next = p.lexer.nextItem()
				switch next.typ {
				case itemLeft, itemRight, itemUp, itemDown:
					cmd.Args = []string{next.val}
				case itemText:
					cmd.Args = []string{p.varReplacement(next.val)}
				default:
					return nil, p.errorf(next, "Unexpected token")
				}
			default:
				return nil, p.errorf(next, "Unexpected token")
			}
			next = p.lexer.nextItem()
		case itemSticky:
			next = p.lexer.nextItem()
			switch next.typ {
			case itemEnable, itemDisable, itemToggle:
				cmd.Args = []string{next.val}
				cmd.Exec = Sticky
			default:
				return nil, p.errorf(next, "Unexpected token")
			}
			next = p.lexer.nextItem()
		case itemRename:
			var args []string
			next = p.lexer.nextItem()
			if next.typ != itemWorkspace {
				return nil, p.errorf(next, "Unexpected token")
			}

			next = p.lexer.nextItem()
			if next.typ == itemText {
				args = append(args, p.varReplacement(next.val))
				next = p.lexer.nextItem()
			}

			if next.typ != itemTo {
				return nil, p.errorf(next, "Unexpected token")
			}

			next = p.lexer.nextItem()
			if next.typ != itemText {
				return nil, p.errorf(next, "Unexpected token")
			}

			cmd.Args = append(args, p.varReplacement(next.val))
			cmd.Exec = Rename
			next = p.lexer.nextItem()
		// case itemResize:
		// 	var item *item
		// 	cmd, item, err = p.parseResize(criteria)
		// 	if err != nil {
		// 		return nil, err
		// 	}
		// 	next = *item

		// itemMode
		// itemBorder
		// itemScratchpad
		default:
			break Loop
		}

		if start == nil {
			start = cmd
			curr = cmd
		} else {
			curr.Next = cmd
			curr = cmd
		}

		switch next.typ {
		case itemSemicolon:
			criteria = nil
			fallthrough
		case itemComma:
			break
		case itemNewline, itemEOF:
			return start, nil
		default:
			break Loop
		}

		next = p.lexer.nextItem()
	}

	return nil, p.errorf(next, "Failed to parse command, invalid token")
}

func (p *parser) parseExec(criteria *Criteria) (*Command, *item, error) {
	args := make([]string, 0)
	for {
		next := p.lexer.nextItem()
		switch next.typ {
		case itemText:
			args = append(args, p.varReplacement(next.val))
		case itemString:
			val := strings.TrimPrefix(strings.TrimSuffix(next.val, `"`), `"`)
			args = append(args, p.varReplacement(val))
		case itemNewline, itemEOF:
			cmd := &Command{
				Args:     args,
				Criteria: criteria,
				Exec:     Exec,
			}
			return cmd, &next, nil
		default:
			return nil, nil, p.errorf(next, "Invalid exec command")
		}
	}
}

func (p *parser) parseMove(criteria *Criteria) (*Command, *item, error) {
	cmd := &Command{Criteria: criteria}

	next := p.lexer.nextItem()
	switch next.typ {
	case itemLeft, itemRight, itemUp, itemDown:
		cmd.Args = []string{next.val}
		cmd.Exec = MoveDirection
	case itemWindow, itemContainer, itemWorkspace:
		args := []string{next.val}

		move := next
		next = p.lexer.nextItem()
		if next.typ != itemTo {
			return nil, nil, p.errorf(next, "Unexpected token")
		}

		to := p.lexer.nextItem()
		switch to.typ {
		case itemWorkspace:
			if move.typ == itemWorkspace {
				return nil, nil, p.errorf(to, "Can't move workspace to workspace")
			}
		case itemOutput:
		default:
			return nil, nil, p.errorf(to, "Unexpected token")
		}

		args = append(args, to.val)

		next = p.lexer.nextItem()
		switch next.typ {
		case itemNumber:
			args = append(args, next.val)
			next = p.lexer.nextItem()
			if next.typ != itemText {
				return nil, nil, p.errorf(next, "Unexpected token")
			}
			args = append(args, p.varReplacement(next.val))
		case itemLeft, itemRight, itemUp, itemDown:
			if to.typ == itemWorkspace {
				return nil, nil, p.errorf(next, "Can't move to workspace at <direction>")
			}
			args = append(args, next.val)
		default:
			return nil, nil, p.errorf(to, "Unexpected token")
		}

		cmd.Args = args
		cmd.Exec = MoveContainer
	default:
		return nil, nil, p.errorf(next, "Unexpected token")
	}

	next = p.lexer.nextItem()
	return cmd, &next, nil
}

// func (p *parser) parseResize(criteria *Criteria) (*Command, *item, error) {
// 	var args []string
// 	cmd := &Command{Criteria: criteria}

// 	next := p.lexer.nextItem()
// 	switch next.typ {
// 	case itemGrow, itemShrink:
// 		args = append(args, next.val)
// 		next = p.lexer.nextItem()
// 		switch next.typ {
// 		case itemLeft, itemRight, itemUp, itemDown, itemHeight, itemWidth:
// 			args = append(args, next.val)
// 		default:
// 			return nil, nil, p.errorf(next, "Unexpected token")
// 		}
// 		next = p.lexer.nextItem()
// 		switch next.typ {
// 		case itemText:
// 		case itemNewline:
// 			cmd.Args = args
// 			cmd.Exec = Resize
// 		}

// 	}

// 	next = p.lexer.nextItem()
// 	return cmd, &next, nil
// }

func (p *parser) errorf(curr item, format string, a ...interface{}) error {
	var msg bytes.Buffer
	msg.WriteString(fmt.Sprintf(format, a...))
	msg.WriteByte('\n')
	msg.WriteString((*p.input)[p.cmdStart.pos : int(curr.pos)+len(curr.val)])
	msg.WriteByte('\n')
	msg.WriteString(strings.Repeat(" ", int(curr.pos-p.cmdStart.pos)))
	msg.WriteString(strings.Repeat("^", len(curr.val)))
	return fmt.Errorf(msg.String())
}

func (p *parser) varReplacement(str string) string {
	return str
}

func varReplacement(varTable map[string]string, variable string) string {
	if v, ok := varTable[variable]; ok {
		return v
	}

	return variable
}
