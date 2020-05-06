// Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.

package units

// Unit defines the unit schema with the relevant fields.
type Unit struct {
	ID      ID     // unit id
	Name    string // unit name
	Display string // unit display name
	VName   string // variable uniquely associated with this unit
}

// New creates a new unit.
func NewUnit(id ID, name, display, vname string) *Unit {
	return &Unit{ID: id, Name: name, Display: display, VName: vname}
}
