// Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.

package parser

import (
	"testing"

	"github.com/facebookresearch/Clinical-Trial-Parser/src/ct/parser/production"
	"github.com/facebookresearch/Clinical-Trial-Parser/src/ct/relation"

	"github.com/stretchr/testify/assert"
)

func TestCYKParsingAlgoritm(t *testing.T) {
	a := assert.New(t)
	g := NewCFGrammar(production.TestRules)

	x := "or"
	y := "and"
	// input pattern: "y x x x"
	input := Items{
		NewItem(ItemType(y), y),
		NewItem(ItemType(x), x),
		NewItem(ItemType(x), x),
		NewItem(ItemType(x), x),
	}
	expected := `{"score":1.000,"tree":{"value":"S","left":{"value":"A","left":{"value":"B","left":{"value":"and"}},"right":{"value":"A","left":{"value":"or"}}},"right":{"value":"B","left":{"value":"C","left":{"value":"or"}},"right":{"value":"C","left":{"value":"or"}}}}}`
	actual := g.BuildTrees(input).String()
	a.Equal(expected, actual)
}

func TestCriteriaParsing(t *testing.T) {
	a := assert.New(t)
	g := NewCFGrammar(production.CriterionRules)

	// input pattern: "A1c > 5.7 % < 10 ECOG 0 1 or 2 or Height < 200 cm"
	input := Items{
		NewItem(ItemType("variable"), "a1c"),
		NewItem(ItemType("comparison"), ">"),
		NewItem(ItemType("number"), "5.7"),
		NewItem(ItemType("unit"), "%"),
		NewItem(ItemType("comparison"), "<"),
		NewItem(ItemType("number"), "10"),
		NewItem(ItemType("variable"), "ecog"),
		NewItem(ItemType("number"), "0"),
		NewItem(ItemType("number"), "1"),
		NewItem(ItemType("or"), "or"),
		NewItem(ItemType("number"), "2"),
		NewItem(ItemType("or"), "or"),
		NewItem(ItemType("variable"), "height"),
		NewItem(ItemType("comparison"), "<"),
		NewItem(ItemType("number"), "200"),
		NewItem(ItemType("unit"), "cm"),
	}
	expected := relation.Relations{
		relation.Parse(`{"id":"100","name":"ecog","value":["0","1","2"]}`),
		relation.Parse(`{"id":"201","name":"height","unit":"cm","upper":{"incl":false,"value":"200"}}`),
		relation.Parse(`{"id":"400","name":"a1c","unit":"%","lower":{"incl":false,"value":"5.7"},"upper":{"incl":false,"value":"10"}}`),
	}
	trees := g.BuildTrees(input)
	actualOrRels, actualAndRels := trees.Relations()
	actualOrRels.SetScore(0)

	a.Len(actualOrRels, 3)
	a.Len(actualAndRels, 0)
	a.Equal(expected, actualOrRels)
}
