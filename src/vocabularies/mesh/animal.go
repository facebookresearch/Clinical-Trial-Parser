// Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.

package mesh

import (
	"strings"
)

var animals = []string{
	"Canine",
	"Canid",
	"Bovine",
	"Equid",
	"Feline",
	"Duck",
	"Gallid",
	"Woodchuck",
	"Cercopithecine",
	"Simian",
	"Suid",
}

func isAnimalConcept(s string) bool {
	for _, a := range animals {
		if strings.Contains(s, a) {
			return true
		}
	}
	return false
}
