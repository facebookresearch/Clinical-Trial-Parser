// Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.

package criteria

import (
	"fmt"
	"reflect"

	"github.com/facebookresearch/Clinical-Trial-Parser/src/ct/relation"
)

// Criterion defines an eligibility criterion record.
type Criterion struct {
	text         string             // raw criterion string
	relations    relation.Relations // parsed criterion from text, may contain multiple sub-criteria
	score        float64
	ClusterID    int
	ClusterTopic string
}

// Criteria defines a slice of eligibility criteria.
type Criteria []*Criterion

// NewCriterion creates a new criterion.
func NewCriterion(text string, score float64, rels relation.Relations) *Criterion {
	return &Criterion{text: text, score: score, relations: rels}
}

// NewCriteria creates a new slice of criteria.
func NewCriteria() Criteria {
	return make(Criteria, 0)
}

// Names returns a concatenated string of criterion/variable names.
func (c *Criterion) Names() string {
	switch len(c.relations) {
	case 0:
		return ""
	case 1:
		return c.relations[0].Name
	default:
		names := c.relations[0].Name
		for _, r := range c.relations {
			names += " " + r.Name
		}
		return names
	}
}

// Score returns the score of the question being correct.
func (c *Criterion) Score() float64 {
	return c.score
}

// Relations returns the parsed relations for the criterion.
func (c *Criterion) Relations() relation.Relations {
	return c.relations
}

// String returns the raw criterion text.
func (c *Criterion) String() string {
	return c.text
}

// Relations returns all parsed relations for the criteria.
func (cs Criteria) Relations() relation.Relations {
	rs := relation.NewRelations()
	for _, c := range cs {
		rs = append(rs, c.relations...)
	}
	return rs
}

// String returns the string of criteria.
func (cs Criteria) String() string {
	str := ""
	for _, c := range cs {
		str += fmt.Sprintf("%s\n", c)
	}
	return str
}

// JSON returns the json string of the relations.
func (c *Criterion) JSON() string {
	return c.relations.JSON()
}

// Contains returns true if cs contains c.
func (c *Criterion) Contains(cs Criteria) bool {
	for i := range cs {
		if reflect.DeepEqual(c, cs[i]) {
			return true
		}
	}
	return false
}
