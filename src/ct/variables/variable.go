// Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.

package variables

// Variable defines the variable schema with the relevant fields.
type Variable struct {
	ID       ID        // variable id
	Kind     Type      // variable type
	Name     string    // variable name
	Display  string    // variable display name
	Range    []string  // value range for nominal and ordinal variables
	NumRange []float64 // value range for numerical variables
	UnitName string    // variable default unit
}

// New creates a new variable.
func NewVariable(id ID, kind Type, name, display string, bounds []string, numBounds []float64, unitName string) *Variable {
	return &Variable{ID: id, Kind: kind, Name: name, Display: display, Range: bounds, NumRange: numBounds, UnitName: unitName}
}

// InRange return true if val is in the valid range or the range is not specified.
func (v *Variable) InRange(val float64) bool {
	if len(v.NumRange) == 0 || v.Kind != Numerical {
		return true
	}
	return v.NumRange[0] <= val && val <= v.NumRange[1]
}
