// Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.

package production

// CriterionRules defines the context-free grammar production rules
// for parsing a clinical-trial eligibility criterion.
var CriterionRules = `

#nonterminals:

S -> C
C -> C X | R
X -> O R | R
R -> V A | A V | V
V -> V1 V2 | V1
V2 -> H V1
A -> L Y | Y Y | B W | B B | B | E
E -> E N | E Z | N
Z -> O N
B -> T L | L T
W -> O B
L -> N U | N
Y -> D L

#terminals:

O -> or | and | punctuation
V1 -> variable | unknown
T -> comparison
N -> number
U -> unit
D -> range | and
H -> slash

`
