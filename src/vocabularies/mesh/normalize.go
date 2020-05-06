// Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.

package mesh

import (
	"regexp"
	"sort"
	"strings"

	"github.com/facebookresearch/Clinical-Trial-Parser/src/common/col/set"
	"github.com/facebookresearch/Clinical-Trial-Parser/src/common/util/slice"
	"github.com/facebookresearch/Clinical-Trial-Parser/src/common/util/text"
	"github.com/facebookresearch/Clinical-Trial-Parser/src/vocabularies/taxonomy"
)

var (
	reParenthesis = regexp.MustCompile(`(^| )\(.*?\)|\[.*?\]( |$)`)
	rePunct       = regexp.MustCompile(`[,.;:()\[\]"']`)
	reEG          = regexp.MustCompile(`\be g\b`)
	reE           = regexp.MustCompile(`\b e$`)
	re1           = regexp.MustCompile(`\bi\b`)
	re2           = regexp.MustCompile(`\bii\b`)
	reHBV         = regexp.MustCompile(`\bhbv\b`)
	reHCV         = regexp.MustCompile(`\bhcv\b`)
	reCNS         = regexp.MustCompile(`\bcns\b`)
	reAML         = regexp.MustCompile(`\baml\b`)
	reNSCLC       = regexp.MustCompile(`\bnsclc\b`)
	reCLL         = regexp.MustCompile(`\bcll\b`)
	reHCC         = regexp.MustCompile(`\bhcc\b`)
	reMM          = regexp.MustCompile(`\bmm\b`)
	reGI          = regexp.MustCompile(`\bgi\b`)
	reMRI         = regexp.MustCompile(`\bmri\b`)
)

// Normalize defines a normalizer function for MeSH terms.
// normalizedTerm replaces the extracted NER term.
// normalizedMatch is used to match terms to concepts.
var Normalize taxonomy.Normalizer = func(str string) (string, string) {
	s := strings.ToLower(str)
	s = reEG.ReplaceAllString(s, "")
	s = reE.ReplaceAllString(s, "")
	s = strings.Replace(s, ",", " ", -1)
	s = strings.Replace(s, "b hbsag", "hbv surface antigen", -1)
	s = strings.Replace(s, "hbsag hbv", "hbv surface antigen", -1)
	s = strings.Replace(s, "her2", "her-2", -1)

	s = reCNS.ReplaceAllString(s, "central nervous system")
	s = reAML.ReplaceAllString(s, "acute myeloid leukemia")
	s = reNSCLC.ReplaceAllString(s, "non-small cell lung cancer")
	s = reCLL.ReplaceAllString(s, "chronic lymphocytic leukemia")
	s = reHCC.ReplaceAllString(s, "hepatocellular carcinoma")
	s = reMM.ReplaceAllString(s, "multiple myeloma")
	s = reGI.ReplaceAllString(s, "gastrointestinal")
	s = reMRI.ReplaceAllString(s, "magnetic resonance imaging")

	if strings.Contains(s, "diabetes") {
		s = re1.ReplaceAllString(s, "1")
		s = re2.ReplaceAllString(s, "2")
	}

	if strings.Contains(s, "hepatitis") {
		s = strings.Replace(s, "b hbv", "b", -1)
		s = strings.Replace(s, "c hcv", "c", -1)
		s = strings.Replace(s, "active ", "", -1)
		s = strings.Replace(s, " treatment", "", -1)
	} else {
		s = reHBV.ReplaceAllString(s, "b hepatitis")
		s = reHCV.ReplaceAllString(s, "c hepatitis")
		s = strings.Replace(s, "b c hep", "b c hepatitis", -1)
	}
	if len(s) == 0 {
		s = str
	}
	s = strings.Trim(s, " /.,;:-")
	normalizedTerm := text.NormalizeWhitespace(s)

	normalizedMatch := reParenthesis.ReplaceAllString(normalizedTerm, " ")
	normalizedMatch = rePunct.ReplaceAllString(normalizedMatch, " ")
	normalizedMatch = strings.TrimSpace(normalizedMatch)

	l := strings.Fields(normalizedMatch)
	l = filter(l, generalWords)
	if v := strings.Join(l, " "); len(v) > 0 {
		normalizedTerm = v
	}

	l = filter(l, labelWords)
	sort.Strings(l)
	l = slice.Dedupe(l)
	normalizedMatch = strings.Join(l, " ")
	if len(normalizedMatch) == 0 {
		normalizedMatch = normalizedTerm
	}
	return normalizedMatch, normalizedTerm
}

func filter(l []string, remove set.Set) []string {
	n := 0
	for _, a := range l {
		if !remove[a] {
			l[n] = a
			n++
		}
	}
	return l[:n]
}

var generalWords = set.New(
	"",
	"and",
	"and/or",
	"are",
	"as",
	"at",
	"by",
	"for",
	"in",
	"is",
	"its",
	"may",
	"no",
	"not",
	"of",
	"on",
	"or",
	"that",
	"the",
	"were",
	"who",
	"with",

	">",
	"≥",
	"<",
	"≤",
	"=",
	"equal",
	"greater",
	"least",
	"similar",
	"smaller",

	"another",
	"based",
	"before",
	"days",
	"defined",
	"during",
	"including",
	"total",
	"within",
	"without",

	"emoticon",
	"@number",

	"acceptable",
	"adequately",
	"allowed",
	"currently",
	"definitively",
	"demonstrated",
	"eligible",
	"evidence",
	"evidenced",
	"exception",
	"exceptions",
	"excluded",
	"excluding",
	"indicating",
	"ineligible",
	"locally",
	"management",
	"participate",
	"permitted",
	"presence",
	"presence",
	"present",
	"provide",
	"receiving",
	"resected",
	"treated",
	"undergoing",
	"unspecified",
	"urgent",

	"adult",
	"human",
	"individuals",
	"participants",
	"patient",
	"patients",
	"subjects",
	"victim",

	"stage",
	"grade",
	"3",
	"3a",
	"4",
	"ia",
	"ib",
	"ic",
	"iia",
	"iib",
	"iic",
	"iii",
	"iiia",
	"iiib",
	"iiic",
	"iv",
	"iva",
	"ivb",
	"ivc",
	"ajcc",
	"v6",
	"v7",
	"v8",
	"hbsag",
	"hepbsag",
)

var labelWords = set.New(
	"antibiotics",
	"documented",
	"dysfunction",
	"function",
	"impaired",
	"impairment",
	"infection",
	"infectious",
	"language",
	"testing",
	"therapy",
	"treatment",

	"advanced",
	"negative",
	"ongoing",
	"positive",
	"positivity",
	"serious",
	"severe",
	"symptomatic",
)
