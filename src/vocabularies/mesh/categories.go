// Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.

package mesh

import (
	"strings"

	"github.com/facebookresearch/Clinical-Trial-Parser/src/common/col/set"
	"github.com/facebookresearch/Clinical-Trial-Parser/src/common/util/text"
)

// Categories lists the top-level categories in the MeSH descriptor hierarchy.
var Categories = map[string]string{
	"A": "Anatomy",
	"B": "Organisms",
	"C": "Diseases",
	"D": "Chemicals and Drugs",
	"E": "Analytical, Diagnostic and Therapeutic Techniques and Equipment",
	"F": "Psychiatry and Psychology",
	"G": "Biological Sciences",
	"H": "Physical Sciences",
	"I": "Anthropology, Education, Sociology and Social Phenomena",
	"J": "Technology and Food and Beverages",
	"K": "Humanities",
	"L": "Information Science",
	"M": "Persons",
	"N": "Health Care",
	"V": "Publication Characteristics",
	"Z": "Geographic Locations",
}

// ClinicalCategories lists the top-level clinical categories in the MeSH descriptor hierarchy.
var ClinicalCategories = map[string]string{
	"A": "Anatomy",
	"B": "Organisms",
	"C": "Diseases",
	"D": "Chemicals and Drugs",
	"E": "Analytical, Diagnostic and Therapeutic Techniques and Equipment",
	"F": "Psychiatry and Psychology",
	"G": "Biological Sciences",
	"M": "Persons",
}

// AnimalCodes list codes for animal categories
var AnimalCodes = map[string]string{
	"C22": "Animal Diseases",
}

// GetTopCode gets the top code from a tree number.
func GetTopCode(tn string) string {
	return strings.SplitN(tn, ".", 2)[0]
}

// GetTopCodes gets the top codes from tree numbers.
func GetTopCodes(tns []string) set.Set {
	codes := set.New()
	for _, tn := range tns {
		codes.Add(GetTopCode(tn))
	}
	return codes
}

// GetCategories gets the categories from tree numbers.
func GetCategories(tns []string) set.Set {
	categories := set.New()
	for _, tn := range tns {
		categories.Add(text.LetterPrefix(tn))
	}
	return categories
}

// HasClinicalCategory returns true if any of the categories are clinical.
func HasClinicalCategory(tns []string) bool {
	categories := GetCategories(tns)
	for k := range categories {
		if _, ok := ClinicalCategories[k]; ok {
			return true
		}
	}
	return false
}

// HasAnimalCode returns true if any of the top codes are for animals.
func HasAnimalCode(tns []string) bool {
	codes := GetTopCodes(tns)
	for k := range codes {
		if _, ok := AnimalCodes[k]; ok {
			return true
		}
	}
	return false
}

// Trim removes non-clinical codes.
func Trim(tns []string) []string {
	n := 0
	for i, tn := range tns {
		if c := text.LetterPrefix(tn); ClinicalCategories[c] != "" {
			tns[n] = tns[i]
			n++
		}
	}
	return tns[:n]
}
