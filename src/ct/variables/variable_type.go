// Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.

package variables

import (
	"strings"
)

// Type defines the variable type.
type Type string

const (
	// Unknown type
	Unknown Type = ""
	// Boolean type of variable
	Boolean Type = "boolean"
	// Nominal (non-boolean) type of variable
	Nominal Type = "nominal"
	// Ordinal type of variable
	Ordinal Type = "ordinal"
	// Numerical (interval) type of variable
	Numerical Type = "numerical"
)

// ParseType converts the string to the variable type.
func ParseType(s string) Type {
	s = strings.ToLower(strings.TrimSpace(s))
	switch s {
	case "boolean", "nominal", "ordinal", "numerical":
		return Type(s)
	default:
		return Unknown
	}
}

// String returns the corresponding string representation of the variable type.
func (t Type) String() string {
	return string(t)
}

// Types defines a slice of variable types
type Types []Type

// ParseTypes converts the string to the variable types.
func ParseTypes(s string) Types {
	v := strings.Split(s, ",")
	types := make(Types, 0)
	for _, a := range v {
		types = append(types, ParseType(a))
	}
	return types
}

// String returns the corresponding string representation of the variable types.
func (ts Types) String() string {
	switch len(ts) {
	case 0:
		return ""
	case 1:
		return ts[0].String()
	default:
		s := ts[0].String()
		for i := 1; i < len(ts); i++ {
			s += "," + ts[i].String()
		}
		return s
	}
}

// Set converts the slice of variable types to the set.
func (ts Types) Set() map[Type]bool {
	set := make(map[Type]bool)
	for _, t := range ts {
		set[t] = true
	}
	return set
}
