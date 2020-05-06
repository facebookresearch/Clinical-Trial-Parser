// Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.

package vocabularies

import (
	"strings"
)

// Source defines the vocabulary source.
type Source int

const (
	// Unknown vocabulary source
	Unknown Source = iota
	// MeSH vocabulary source
	MESH
	// UMLS vocabulary source
	UMLS
)

// ParseSource converts the string to the vocabulary source.
func ParseSource(s string) Source {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "mesh":
		return MESH
	case "umls":
		return UMLS
	default:
		return Unknown
	}
}

// String returns the corresponding string representation of source.
func (s Source) String() string {
	switch s {
	case MESH:
		return "mesh"
	case UMLS:
		return "umls"
	default:
		return "unknown"
	}
}
