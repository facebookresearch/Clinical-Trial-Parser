// Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.

package parser

import (
	"strings"

	"github.com/facebookresearch/Clinical-Trial-Parser/src/common/col/set"
	"github.com/facebookresearch/Clinical-Trial-Parser/src/common/util/slice"

	"github.com/golang/glog"
)

// Rules define the grammar production rules.
type Rules struct {
	terminalRules  map[itemType]set.Set
	unaryRules     map[Element]set.Set
	binaryRules    map[Element]set.Set
	nonTerminalSet set.Set
}

// RuleType defines the type of rules that are being loaded from the string.
type RuleType int

const (
	unknownRule RuleType = iota
	terminalRule
	nonTerminalRule
)

// LoadRules loads the grammar production rules from the string.
func LoadRules(s string) *Rules {
	terminalRules := map[itemType]set.Set{}
	unaryRules := map[Element]set.Set{}
	binaryRules := map[Element]set.Set{}
	nonTerminalSet := set.New()

	ruleType := unknownRule

	lines := strings.Split(s, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)

		switch {
		case strings.HasPrefix(line, "#terminals"):
			ruleType = terminalRule
		case strings.HasPrefix(line, "#nonterminals"):
			ruleType = nonTerminalRule
		}

		if ruleType == unknownRule || len(line) == 0 || line[0] == '#' {
			continue
		}

		values := strings.Split(line, "->")
		if len(values) != 2 {
			continue
		}
		slice.TrimSpace(values)
		A := values[0]
		list := strings.Split(values[1], "|")
		slice.TrimSpace(list)

		nonTerminalSet.Add(A)

		switch ruleType {
		case terminalRule:
			for _, s := range list {
				a := ItemType(s)
				if _, ok := terminalRules[a]; !ok {
					terminalRules[a] = set.New()
				}
				terminalRules[a].Add(A)
			}
		case nonTerminalRule:
			for _, a := range list {
				symbols := strings.Fields(a)
				slice.TrimSpace(symbols)
				switch len(symbols) {
				case 0:
					glog.Fatalf("Cannot read production rule: %s\n", line)
				case 1:
					e := NewUnary(symbols[0])
					if _, ok := unaryRules[e]; !ok {
						unaryRules[e] = set.New()
					}
					unaryRules[e].Add(A)
					nonTerminalSet.Add(e.leftNonTerminal)
				default:
					e := NewBinary(symbols[0], symbols[1])
					if _, ok := binaryRules[e]; !ok {
						binaryRules[e] = set.New()
					}
					binaryRules[e].Add(A)
					nonTerminalSet.Add(e.leftNonTerminal)
					nonTerminalSet.Add(e.rightNonTerminal)
				}
			}
		}
	}

	return &Rules{
		terminalRules:  terminalRules,
		unaryRules:     unaryRules,
		binaryRules:    binaryRules,
		nonTerminalSet: nonTerminalSet,
	}
}
