// Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.

package production

// TestRules defines the context-free grammar production rules
// for testing grammar algorithms.
var TestRules = `

#nonterminals:

S -> A B | B C
A -> B A
B -> C C
C -> A B

#terminals:

A -> or
B -> and
C -> or

`
