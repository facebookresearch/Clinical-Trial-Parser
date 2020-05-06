// Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.

package eligibility

import (
	"strings"
)

// Type defines the eligibility criteria type.
type Type int

const (
	// Unknown is the type of unknown criteria
	Unknown Type = iota
	// Inclusion is the type of inclusion criteria
	Inclusion
	// Exclusion is the type of exclusion criteria
	Exclusion
)

// ParseType converts a string to an eligibility type.
func ParseType(s string) Type {
	switch strings.ToLower(s) {
	case "inclusion":
		return Inclusion
	case "exclusion":
		return Exclusion
	default:
		return Unknown
	}
}

// String converts the type to a string.
func (t Type) String() string {
	switch t {
	case Inclusion:
		return "inclusion"
	case Exclusion:
		return "exclusion"
	default:
		return "unknown"
	}
}
