// Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.

package criteria

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalize(t *testing.T) {
	a := assert.New(t)

	input := "Inclusion Criteria:\ni1\ni2\nExclusion Criteria:\nDoes not meet inclusion criteria.\ne1\ne2"
	expected := "Inclusion Criteria:\ni1\ni2\nExclusion Criteria:\ne1\ne2"
	actual := Normalize(input)

	a.Equal(expected, actual)
}

func TestNormalizeWithoutExclusionCriteria(t *testing.T) {
	a := assert.New(t)

	input := "Inclusion Criteria:\ni1\ni2\nExclusion Criteria:\nDoes not meet inclusion criteria.\n"
	expected := "Inclusion Criteria:\ni1\ni2\nExclusion Criteria:\n"
	actual := Normalize(input)

	a.Equal(expected, actual)
}

func TestExtractCriteria(t *testing.T) {
	a := assert.New(t)

	input := "Inclusion:\ni1\ni2\nExclusion:\ne1\ne2"
	expectedInclusions := []string{"i1\ni2"}
	expectedExclusions := []string{"e1\ne2"}

	actualInclusions := ExtractInclusionCriteria(input)
	a.Equal(expectedInclusions, actualInclusions)

	actualExclusions := ExtractExclusionCriteria(input)
	a.Equal(expectedExclusions, actualExclusions)
}

func TestExtractCriteriaLong(t *testing.T) {
	a := assert.New(t)

	input := "Inclusion criteria for all:\ni1\ni2\nSubject exclusion criteria:\ne1\ne2"
	expectedInclusions := []string{"i1\ni2"}
	expectedExclusions := []string{"e1\ne2"}

	actualInclusions := ExtractInclusionCriteria(input)
	a.Equal(expectedInclusions, actualInclusions)

	actualExclusions := ExtractExclusionCriteria(input)
	a.Equal(expectedExclusions, actualExclusions)
}

func TestExtractCriteriaMultiple(t *testing.T) {
	a := assert.New(t)

	input := "Inclusion:\ni1\ni2\nExclusion:\ne1\ne2\nInclusion:\nj1\nj2\nExclusion:\nf1\nf2"
	expectedInclusions := []string{"i1\ni2", "j1\nj2"}
	expectedExclusions := []string{"e1\ne2", "f1\nf2"}

	actualInclusions := ExtractInclusionCriteria(input)
	a.Equal(expectedInclusions, actualInclusions)

	actualExclusions := ExtractExclusionCriteria(input)
	a.Equal(expectedExclusions, actualExclusions)
}

func TestExtractCriteriaMessy(t *testing.T) {
	a := assert.New(t)

	input := "Inclusion Criteria\ni1\ni2 Exclusion Criteria.\nKey Exclusion Criteria\nInclusion criteria. e1\ne2"
	expectedInclusions := []string{"i1\ni2 Exclusion Criteria."}
	expectedExclusions := []string{"Inclusion criteria. e1\ne2"}

	actualInclusions := ExtractInclusionCriteria(input)
	a.Equal(expectedInclusions, actualInclusions)

	actualExclusions := ExtractExclusionCriteria(input)
	a.Equal(expectedExclusions, actualExclusions)
}

func TestExtractCriteriaWithoutExclusions(t *testing.T) {
	a := assert.New(t)

	input := "Inclusion Criteria:\ni1\ni2\nExclusion Criteria:\nDoes not meet inclusion criteria.\n"
	expectedInclusions := []string{"i1\ni2"}
	expectedExclusions := empty

	actual := Normalize(input)
	actualInclusions := ExtractInclusionCriteria(actual)
	a.Equal(expectedInclusions, actualInclusions)

	actualExclusions := ExtractExclusionCriteria(actual)
	a.Equal(expectedExclusions, actualExclusions)
}

func TestExtractCriteriaWithExtraWords(t *testing.T) {
	a := assert.New(t)

	input := "i) INCLUSION/EXCLUSION CRITERIA:\nStudy eligibility: General Inclusions:\nInclusions: 1. one\nKey Exclusion Criteria:\n Exclusions: 2. two"
	expectedInclusions := []string{"Inclusions: 1. one"}
	expectedExclusions := []string{"Exclusions: 2. two"}

	actual := Normalize(input)
	actualInclusions := ExtractInclusionCriteria(actual)
	a.Equal(expectedInclusions, actualInclusions)

	actualExclusions := ExtractExclusionCriteria(actual)
	a.Equal(expectedExclusions, actualExclusions)
}

func TestExtractCriteriaWithMissingNewlines(t *testing.T) {
	a := assert.New(t)

	input := "Study eligibility: General Inclusion Criteria : Inclusions: 1. one\nKey Exclusion Criteria: Exclusions: 2. two"
	expectedInclusions := []string{"Inclusions: 1. one"}
	expectedExclusions := []string{"Exclusions: 2. two"}

	actualInclusions := ExtractInclusionCriteria(input)
	a.Equal(expectedInclusions, actualInclusions)

	actualExclusions := ExtractExclusionCriteria(input)
	a.Equal(expectedExclusions, actualExclusions)
}

func TestTrimCriterion(t *testing.T) {
	a := assert.New(t)

	input := "5. Agree to participate  "
	expected := "Agree to participate"

	actual := TrimCriterion(input)
	a.Equal(expected, actual)

	input = "- Agree to participate  "
	expected = "Agree to participate"

	actual = TrimCriterion(input)
	a.Equal(expected, actual)

	input = "  -   78 Agree to participate  "
	expected = "Agree to participate"

	actual = TrimCriterion(input)
	a.Equal(expected, actual)

	input = "prior A1c 5.7 - 6.4"
	expected = "prior A1c 5.7 - 6.4"

	actual = TrimCriterion(input)
	a.Equal(expected, actual)
}

func TestCheckFollowingLine(t *testing.T) {
	a := assert.New(t)
	input := "-  No clinically significant cardiovascular disease, including any of the following"
	foundTab := false
	header := ""
	rule, header, foundTab := checkLine(input, header, foundTab)

	a.Empty(rule)
	a.Equal(header, input)
	a.True(foundTab)
}

func TestCheckBulletLine(t *testing.T) {
	a := assert.New(t)
	input := " -  Sentinel node biopsy alone (if sentinel node is negative)"
	expectedInput := TrimCriterion(input)
	foundTab := true
	header := "-  Axilla must be staged by one of the following"

	expected := header + " " + expectedInput

	rule, newHeader, foundTab := checkLine(input, header, foundTab)

	a.Equal(rule, expected)
	a.Equal(header, newHeader)
	a.True(foundTab)
}

func TestCheckNumberLine(t *testing.T) {
	a := assert.New(t)
	input := " 3. Sentinel node biopsy alone (if sentinel node is negative)"
	expectedInput := TrimCriterion(input)
	foundTab := true
	header := "-  Axilla must be staged by one of the following"

	expected := header + " " + expectedInput

	rule, newHeader, foundTab := checkLine(input, header, foundTab)

	a.Equal(rule, expected)
	a.Equal(header, newHeader)
	a.True(foundTab)
}

func TestCheckNormalLineBefore(t *testing.T) {
	a := assert.New(t)
	input := "-  Multifocal breast cancer is allowed if the intent is to undergo resection through a single lumpectomy incision"
	foundTab := false
	header := ""

	rule, header, foundTab := checkLine(input, header, foundTab)

	a.Equal(rule, input)
	a.Empty(header)
	a.False(foundTab)
}

func TestCheckNormalLineAfter(t *testing.T) {
	a := assert.New(t)
	input := "Appropriate stage for protocol entry including no clinical evidence for distant metastases"
	foundTab := true
	header := "-  Axilla must be staged by the following"

	rule, header, foundTab := checkLine(input, header, foundTab)

	a.Equal(rule, input)
	a.Empty(header)
	a.False(foundTab)
}

func TestSplit(t *testing.T) {
	a := assert.New(t)
	input := `6. Meet the following criteria:

              - cr 1

              - cr 2

              - cr 3`
	expected := []string{
		"6. Meet the following criteria: cr 1",
		"6. Meet the following criteria: cr 2",
		"6. Meet the following criteria: cr 3",
	}

	actual := Split(input)
	a.Equal(expected, actual)
}
