// Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.

package studies

import (
	"fmt"
	"strings"

	"github.com/facebookresearch/Clinical-Trial-Parser/src/common/col/set"
	"github.com/facebookresearch/Clinical-Trial-Parser/src/common/util/slice"
	"github.com/facebookresearch/Clinical-Trial-Parser/src/ct/criteria"
	"github.com/facebookresearch/Clinical-Trial-Parser/src/ct/parser"
	"github.com/facebookresearch/Clinical-Trial-Parser/src/ct/relation"
	"github.com/facebookresearch/Clinical-Trial-Parser/src/ct/variables"
)

// Study defines a record for a clinical study.
type Study struct {
	nct                 string   // National clinical trial identifier
	name                string   // Study name
	conditions          []string // Conditions
	eligibilityCriteria string   // Eligibility criteria

	inclusionCriteria criteria.Criteria
	exclusionCriteria criteria.Criteria
	criteriaCnt       int
}

// NewStudy creates a record for a new study.
func NewStudy(nct, name string, conditions []string, eligibilityCriteria string) *Study {
	return &Study{nct: nct, name: name, conditions: conditions, eligibilityCriteria: eligibilityCriteria}
}

// NCT returns the national clinical trial id.
func (s *Study) NCT() string {
	return s.nct
}

// Name returns the study name.
func (s *Study) Name() string {
	return s.name
}

// InclusionCriteria returns the inclusion criteria for the study.
func (s *Study) InclusionCriteria() criteria.Criteria {
	return s.inclusionCriteria
}

// ExclusionCriteria returns the exclusion criteria for the study.
func (s *Study) ExclusionCriteria() criteria.Criteria {
	return s.exclusionCriteria
}

// Parse parses eligibility criteria text to relations for the study s.
func (s *Study) Parse() *Study {
	interpreter := parser.Get()

	inclusions, exclusions := s.Criteria()
	s.criteriaCnt = len(inclusions) + len(exclusions)

	// Parse inclusion criteria:
	inclusionCriteria := criteria.NewCriteria()
	for _, inclusion := range inclusions {
		lowercase := strings.ToLower(inclusion)
		orRelations, andRelations := interpreter.Interpret(lowercase)

		orRelations.Process()
		andRelations.Process()

		if !orRelations.Empty() {
			criterion := criteria.NewCriterion(inclusion, orRelations.MinScore(), orRelations)
			inclusionCriteria = append(inclusionCriteria, criterion)
		}
		if !andRelations.Empty() {
			for _, r := range andRelations {
				rs := relation.Relations{r}
				criterion := criteria.NewCriterion(inclusion, rs.MinScore(), rs)
				inclusionCriteria = append(inclusionCriteria, criterion)
			}
		}
	}
	s.inclusionCriteria = inclusionCriteria

	// Parse exclusion criteria:
	exclusionCriteria := criteria.NewCriteria()
	for _, exclusion := range exclusions {
		lowercase := strings.ToLower(exclusion)
		orRelations, andRelations := interpreter.Interpret(lowercase)
		orRelations.Process()
		andRelations.Process()
		orRelations.Negate()
		andRelations.Negate()

		if !andRelations.Empty() {
			criterion := criteria.NewCriterion(exclusion, andRelations.MinScore(), andRelations)
			exclusionCriteria = append(exclusionCriteria, criterion)
		}
		if !orRelations.Empty() {
			for _, r := range orRelations {
				rs := relation.Relations{r}
				criterion := criteria.NewCriterion(exclusion, rs.MinScore(), rs)
				exclusionCriteria = append(exclusionCriteria, criterion)
			}
		}
	}

	s.exclusionCriteria = exclusionCriteria
	s.Transform()

	return s
}

// Criteria extracts inclusion and exclusion criteria from the eligibility criteria string.
func (s *Study) Criteria() ([]string, []string) {
	eligibilityCriteria := criteria.Normalize(s.eligibilityCriteria)

	// Parse inclusion criteria:
	var inclusions []string
	inclusionList := criteria.ExtractInclusionCriteria(eligibilityCriteria)
	for _, s := range inclusionList {
		inclusions = append(inclusions, criteria.Split(s)...)
	}

	for i, c := range inclusions {
		inclusions[i] = criteria.TrimCriterion(c)
	}
	inclusions = slice.RemoveEmpty(inclusions)

	// Parse exclusion criteria:
	var exclusions []string
	exclusionList := criteria.ExtractExclusionCriteria(eligibilityCriteria)
	for _, s := range exclusionList {
		exclusions = append(exclusions, criteria.Split(s)...)
	}

	for i, c := range exclusions {
		exclusions[i] = criteria.TrimCriterion(c)
	}
	exclusions = slice.RemoveEmpty(exclusions)

	return inclusions, exclusions
}

// Transform transforms criteria relations by converting parsed values to strings of valid literals.
// If a valid literal cannot be inferred, the confidence score of the relation is set to zero.
func (s *Study) Transform() {
	s.inclusionCriteria.Relations().Transform()
	s.exclusionCriteria.Relations().Transform()
}

// Relations returns the string representation of the parsed criteria.
// Relations that are parsed from the same criterion and are conjoined
// by 'or' have the same criterion id (cid).
func (s *Study) Relations() string {
	variableCatalog := variables.Get()
	relations := ""
	cid := 0
	for _, c := range s.inclusionCriteria {
		for _, r := range c.Relations() {
			q := variableCatalog.Question(r.ID)
			relations += fmt.Sprintf("%s\t%s\t%s\t%d\t%s\t%s\t%s\n",
				s.nct, "inclusion", r.VariableType.String(), cid, c.String(), q, r.JSON())
		}
		cid++
	}
	for _, c := range s.exclusionCriteria {
		for _, r := range c.Relations() {
			q := variableCatalog.Question(r.ID)
			relations += fmt.Sprintf("%s\t%s\t%s\t%d\t%s\t%s\t%s\n",
				s.nct, "exclusion", r.VariableType.String(), cid, c.String(), q, r.JSON())
		}
		cid++
	}
	return relations
}

// CriteriaCount returns the number of criteria.
func (s *Study) CriteriaCount() int {
	return s.criteriaCnt
}

// ParsedCriteriaCount returns the number of parsed unique criteria.
func (s *Study) ParsedCriteriaCount() int {
	parsedCriteria := set.New()
	for _, c := range s.inclusionCriteria {
		for _, r := range c.Relations() {
			if r.Valid() {
				parsedCriteria.Add(c.String())
				break
			}
		}
	}
	for _, c := range s.exclusionCriteria {
		for _, r := range c.Relations() {
			if r.Valid() {
				parsedCriteria.Add(c.String())
				break
			}
		}
	}
	return parsedCriteria.Size()
}

// RelationCount returns the number of parsed relations.
func (s *Study) RelationCount() int {
	cnt := 0
	for _, c := range s.inclusionCriteria {
		for _, r := range c.Relations() {
			if r.Valid() {
				cnt++
			}
		}
	}
	for _, c := range s.exclusionCriteria {
		for _, r := range c.Relations() {
			if r.Valid() {
				cnt++
			}
		}
	}
	return cnt
}
