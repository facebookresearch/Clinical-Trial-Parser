// Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.

package parser

import (
	"fmt"
	"sort"

	"github.com/facebookresearch/Clinical-Trial-Parser/src/common/col/set"
	"github.com/facebookresearch/Clinical-Trial-Parser/src/ct/units"
)

// itemType defines the type of parsed items.
type itemType int

const (
	itemUnknown itemType = iota
	itemOr
	itemAnd
	itemPunctuation
	itemSlash
	itemVariable
	itemComparison
	itemRange
	itemNumber
	itemUnit
)

// ItemType converts a string to itemType.
func ItemType(s string) itemType {
	switch s {
	case "or":
		return itemOr
	case "and":
		return itemAnd
	case "punctuation":
		return itemPunctuation
	case "slash":
		return itemSlash
	case "variable":
		return itemVariable
	case "comparison":
		return itemComparison
	case "range":
		return itemRange
	case "number":
		return itemNumber
	case "unit":
		return itemUnit
	default:
		return itemUnknown
	}
}

// String converts itemType to string.
func (t itemType) String() string {
	switch t {
	case itemOr:
		return "or"
	case itemAnd:
		return "and"
	case itemPunctuation:
		return "punctuation"
	case itemSlash:
		return "slash"
	case itemVariable:
		return "variable"
	case itemComparison:
		return "comparison"
	case itemRange:
		return "range"
	case itemNumber:
		return "number"
	case itemUnit:
		return "unit"
	default:
		return "unknown"
	}
}

// Item defines the lexical item that is syntactically constructed from the lexer token.
type Item struct {
	typ itemType
	val string
}

// NewItem creates a new item.
func NewItem(typ itemType, val string) *Item {
	return &Item{typ: typ, val: val}
}

// UnknownItem creates an unknown item.
func UnknownItem() *Item {
	return NewItem(itemUnknown, "")
}

// Set sets the item fields.
func (i *Item) Set(typ itemType, val string) *Item {
	i.typ = typ
	i.val = val
	return i
}

// Copy copies the fields from the other item.
func (i *Item) Copy(j *Item) {
	i.typ = j.typ
	i.val = j.val
}

// Equal tests whether two items have the same fields.
func (i *Item) Equal(j *Item) bool {
	return i.typ == j.typ && i.val == j.val
}

// Valid tests whether the item is valid: the item type is not unknown
// and the value is not empty.
func (i *Item) Valid() bool {
	return i.typ != itemUnknown && len(i.val) > 0
}

// Negate applies 'not' operation to the comparison item.
func (i *Item) Negate() *Item {
	switch i.val {
	case "<":
		return NewItem(itemComparison, "≥")
	case "≤":
		return NewItem(itemComparison, ">")
	case ">":
		return NewItem(itemComparison, "≤")
	case "≥":
		return NewItem(itemComparison, "<")
	default:
		return UnknownItem()
	}
}

// String returns the string representation of the item.
func (i *Item) String() string {
	return fmt.Sprintf("{type:%q,value:%q}", i.typ.String(), i.val)
}

// Items defines a slice of items.
type Items []*Item

// NewItems creates a slice of items.
func NewItems() Items {
	return make(Items, 0)
}

// Len returns the length of the slice.
func (is Items) Len() int {
	return len(is)
}

// Add adds the item to the items.
func (is *Items) Add(i *Item) {
	*is = append(*is, i)
}

// Empty tests whether items has any items in it.
func (is Items) Empty() bool {
	return len(is) == 0
}

func (is Items) LastType() itemType {
	if is.Empty() {
		return itemUnknown
	}
	return is[is.Len()-1].typ
}

// Get gets the items of type 'typ'.
func (is Items) Get(typ itemType) set.Set {
	set := set.New()
	if is == nil {
		return set
	}
	for _, i := range is {
		if i.typ == typ {
			set.Add(i.val)
		}
	}
	return set
}

// FixMissingVariable adds missing variable to the items,
// if it can be inferred from the unit.
func (is *Items) FixMissingVariable() bool {
	if !is.Get(itemVariable).Empty() {
		return true
	}
	unitCatalog := units.Get()
	candidates := set.New()
	for u := range is.Get(itemUnit) {
		if v, ok := unitCatalog.Variable(u); ok {
			candidates.Add(v)
		}
	}
	if candidates.Size() == 1 {
		v, _ := candidates.Get()
		i := NewItem(itemVariable, v)
		*is = append(Items{i}, *is...)
		return true
	}
	return false
}

// TrimUnknownItems merges consecutive unknown items to one, removes
// the first item if it is an unknown item, and removes unknown items
// that precede or follow a variable or unit item.
func (is *Items) TrimUnknownItems() {
	a := *is
	if len(a) < 2 {
		return
	}
	j := 0
	for i := 1; i < len(a); i++ {
		if a[j].typ == itemUnknown && a[i].typ == itemUnknown {
			continue
		}
		j++
		a[j] = a[i]
	}
	a = a[:j+1]
	if len(a) > 0 && a[0].typ == itemUnknown {
		a = a[1:]
	}

	j = 0
	for i := 1; i < len(a); i++ {
		if (a[j].typ == itemVariable || a[j].typ == itemUnit) && a[i].typ == itemUnknown {
			continue
		}
		if a[j].typ == itemUnknown && (a[i].typ == itemVariable || a[i].typ == itemUnit) {
			a[j] = a[i]
			continue
		}
		j++
		a[j] = a[i]
	}

	*is = a[:j+1]
}

// TrimKnownItems merges consecutive variable, unit, and comparison items
// with the same name to one.
func (is *Items) TrimKnownItems() {
	a := *is
	if len(a) < 2 {
		return
	}
	j := 0
	for i := 1; i < len(a); i++ {
		if (a[j].typ == itemVariable || a[j].typ == itemUnit || a[j].typ == itemComparison) && a[i].Equal(a[j]) {
			continue
		}
		j++
		a[j] = a[i]
	}
	*is = a[:j+1]
}

// TrimRangeItems removes a range item if it precedes a comparison item.
func (is *Items) TrimRangeItems() {
	a := *is
	if len(a) < 2 {
		return
	}
	j := 0
	for i := 0; i < len(a); i++ {
		if i < len(a)-1 && a[i].typ == itemRange && a[i+1].typ == itemComparison {
			continue
		}
		a[j] = a[i]
		j++
	}
	*is = a[:j]
}

// String returns the string representation of the items.
func (is Items) String() string {
	s := ""
	for _, i := range is {
		s += i.String()
	}
	return s
}

// List defines a slice of items.
type List []Items

// NewList creates a new list.
func NewList() List {
	return make(List, 0)
}

// Sort sorts the list elements by their length in descending order.
func (l List) Sort() {
	sort.SliceStable(l, func(i, j int) bool {
		return l[i].Len() > l[j].Len()
	})
}

// FixMissingVariable adds missing variables to the list of items,
// if they can be inferred from the unit info.
func (l List) FixMissingVariable() {
	for i := 0; i < len(l); i++ {
		l[i].FixMissingVariable()
	}
}

// TrimItems trims the unknown (typ = itemUnknown) and known items in the list.
func (l List) TrimItems() {
	for i := 0; i < len(l); i++ {
		l[i].TrimUnknownItems()
		l[i].TrimKnownItems()
		l[i].TrimRangeItems()
	}
}

// String returns the string representation of the list.
func (l List) String() string {
	switch len(l) {
	case 0:
		return ""
	case 1:
		return l[0].String()
	default:
		s := l[0].String()
		for i := 1; i < len(l); i++ {
			s += "\n" + l[i].String()
		}
		return s
	}
}
