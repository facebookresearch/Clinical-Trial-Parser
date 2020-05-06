// Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.

package criteria

import (
	"regexp"
	"strings"

	"github.com/facebookresearch/Clinical-Trial-Parser/src/common/util/text"
)

var (
	reDeleteCriterion = regexp.MustCompile(`(?i)([^\n]+meet inclusion criteria|[^\n]*inclusion/exclusion criteria)\W? *(\n|$)`)
	reMatchInclusions = regexp.MustCompile(`(?is)inclusions?(?: *:| criteria(?:[^:\n]*?:| *\n))(.*?)(?:[^\n]*\bexclusions?(?: *:| criteria(?:[^:\n]*?:| *\n))|$)`)
	reMatchExclusions = regexp.MustCompile(`(?is)exclusions?(?: *:| criteria(?:[^:\n]*?:| *\n))(.*?)(?:[^\n]*\binclusions?(?: *:| criteria(?:[^:\n]*?:| *\n))|$)`)

	reCriteriaSplitter = regexp.MustCompile(`\n\n`)
	reTrimmer          = regexp.MustCompile(`^(\s*-\s*)?(\s*\d+\.?\s*)?`)

	reMatchTabs       = regexp.MustCompile(`the following(\s+criteria)?(\s*:)?\s*\n\s*(-|\d+\.|[a-z]\s)\s*`)
	reMatchTabLine    = regexp.MustCompile(`the following`)
	reMatchBulletLine = regexp.MustCompile(`^\s*(-|\d+\.|[a-z]\s)\s*`)

	empty = []string{}
)

// Normalize normalizes eligibility criteria text. For now, non-informative
// "Does not meet inclusion criteria" like criteria are removed.
func Normalize(s string) string {
	s = reDeleteCriterion.ReplaceAllString(s, "")
	return s
}

// ExtractInclusionCriteria extracts a block of inclusion criteria from the string.
func ExtractInclusionCriteria(s string) []string {
	return extractCriteria(s, reMatchInclusions)
}

// ExtractExclusionCriteria extracts a block of exclusion criteria from the string.
func ExtractExclusionCriteria(s string) []string {
	return extractCriteria(s, reMatchExclusions)
}

func extractCriteria(s string, r *regexp.Regexp) []string {
	values := r.FindAllStringSubmatch(s, -1)
	if len(values) == 0 {
		return empty
	}
	c := []string{}
	for _, value := range values {
		if len(value) == 2 {
			if v := strings.TrimSpace(value[1]); len(v) > 0 {
				c = append(c, v)
			}
		}
	}
	return c
}

// Split splits eligibility criteria numberings into individual criteria.
func Split(s string) []string {
	rules := reCriteriaSplitter.Split(s, -1)
	numTabs, header, foundTab := initLine(s)
	if numTabs == 0 {
		return rules
	}
	var newRules []string
	for _, rule := range rules {
		if rule, header, foundTab = checkLine(rule, header, foundTab); len(rule) > 0 {
			newRules = append(newRules, rule)
		}
	}
	return newRules
}

// TrimCriterion normalizes the criterion by removing leading bullets,
// numberings, and all leading and trailing punctuation.
func TrimCriterion(s string) string {
	s = reTrimmer.ReplaceAllString(s, "")
	s = text.NormalizeWhitespace(s)
	s = strings.Trim(s, ` ,.;:/"`)
	return s
}

func checkLine(rule string, header string, foundTab bool) (string, string, bool) {
	// found a bullet for a previously seen header
	if foundTab && reMatchBulletLine.MatchString(rule) {
		rule = header + " " + TrimCriterion(rule)

		// found a header
	} else if reMatchTabLine.MatchString(rule) {
		foundTab = true
		header = rule
		rule = ""

		// normal criteria
	} else {
		foundTab = false
		header = ""
	}

	return rule, header, foundTab
}

func initLine(eligibilities string) (int, string, bool) {
	numTabs := len(reMatchTabs.FindAllStringIndex(eligibilities, -1))
	header := ""
	foundTab := false
	return numTabs, header, foundTab
}
