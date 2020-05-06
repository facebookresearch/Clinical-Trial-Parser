// Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.

package relation

import (
	"encoding/json"
	"fmt"
	"math"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/facebookresearch/Clinical-Trial-Parser/src/common/col/set"
	"github.com/facebookresearch/Clinical-Trial-Parser/src/common/util/slice"
	"github.com/facebookresearch/Clinical-Trial-Parser/src/common/util/text"
	"github.com/facebookresearch/Clinical-Trial-Parser/src/ct/variables"

	"github.com/golang/glog"
)

var (
	reMissingZero = regexp.MustCompile(`,00$`)
	reRadixComma  = regexp.MustCompile(`^\d{1,2},\d$`)
	reTimes       = regexp.MustCompile(`\s*(x|×)\s*`)
)

// Limit defines a lower or upper bound of a numerical relation.
type Limit struct {
	Incl  bool   `json:"incl"`  // True if limit is inclusive
	Value string `json:"value"` // Value of limit bound
}

// Relation defines a boolean, nominal, ordinal, or numerical criterion.
type Relation struct {
	ID           variables.ID   `json:"id,omitempty"`
	Name         string         `json:"name"`            // Relation name, typically the variable name
	DisplayName  string         `json:"-"`               // Variable display name
	Unit         string         `json:"unit,omitempty"`  // Variable unit
	Value        []string       `json:"value,omitempty"` // Valid values of categorical relation
	Lower        *Limit         `json:"lower,omitempty"` // Lower bound of numerical relation condition
	Upper        *Limit         `json:"upper,omitempty"` // Upper bound of numerical relation condition
	VariableType variables.Type `json:"variableType"`    // Type of relation
	Score        float64        `json:"score"`           // Confidence estimate of the relation representation being correct
}

// Relations defines a slice of relations.
type Relations []*Relation

// New creates a new relation.
func New() *Relation {
	return &Relation{VariableType: variables.Unknown}
}

// NewCategorical creates a new relation of type boolean, nominal, or ordinal.
func NewCategorical(v *variables.Variable, answer []string, score float64) *Relation {
	return &Relation{ID: v.ID, Name: v.Name, DisplayName: v.Display, VariableType: v.Kind, Value: answer, Score: score}
}

// NewRelations creates an empty slice of relations.
func NewRelations() Relations {
	return make(Relations, 0)
}

// Parse parsers the json string to the relation.
func Parse(s string) *Relation {
	var r Relation
	if err := json.Unmarshal([]byte(s), &r); err != nil {
		glog.Fatal(err)
	}
	return &r
}

// JSON converts the relation to the json string.
func (r *Relation) JSON() string {
	if b, err := json.Marshal(r); err == nil {
		return string(b)
	}
	return ""
}

// HumanReadable converts the relation to the human readable form.
func (r *Relation) HumanReadable() string {
	if r.VariableType == variables.Numerical {
		var s string
		if r.Lower != nil {
			s = r.DisplayName
			if r.Lower.Incl {
				s += " ≥ "
			} else {
				s += " > "
			}
			s += r.Lower.Value
		}
		if r.Upper != nil {
			if r.Lower != nil {
				if r.Lower.Value < r.Upper.Value {
					s += " and "
				} else {
					s += " or "
				}
			}
			s += r.DisplayName

			if r.Upper.Incl {
				s += " ≤ "
			} else {
				s += " < "
			}
			s += r.Upper.Value
		}
		if r.Unit != "" {
			s += " " + r.Unit
		}
		return s
	}
	return text.Join(text.Titles(r.Value), ", ", " or ")
}

// SetScore sets the confidence score that the relation is parsed correctly.
func (r *Relation) SetScore(score float64) {
	r.Score = score
}

// SetVariableFields sets the variable display name and type per variable id.
func (r *Relation) SetVariableFields(v *variables.Variable) {
	r.DisplayName = v.Display
	r.VariableType = v.Kind
}

// SetUnitField sets the variable unit to the default value if the unit missing.
func (r *Relation) SetUnitField(v *variables.Variable) {
	r.Unit = v.UnitName
}

// Normalize normalizes the relation by making the relation content
// consistent with the variable type.
func (r *Relation) Normalize(valueRange []string) {
	switch r.VariableType {
	case variables.Boolean, variables.Nominal:
		r.Lower = nil
		r.Upper = nil
	case variables.Ordinal:
		r.romanToArabicNumerals()
		if r.Lower != nil || r.Upper != nil {
			r.Value = r.apply(valueRange)
		}
		r.Lower = nil
		r.Upper = nil
	case variables.Numerical:
		if r.Lower != nil || r.Upper != nil {
			r.Value = nil
		}
	}
}

// romanToArabicNumerals converts values and limits from roman to arabic numerals.
func (r *Relation) romanToArabicNumerals() {
	for i, a := range r.Value {
		r.Value[i] = text.RomanToArabicNumerals(a)
	}
	if r.Lower != nil {
		r.Lower.Value = text.RomanToArabicNumerals(r.Lower.Value)
	}
	if r.Upper != nil {
		r.Upper.Value = text.RomanToArabicNumerals(r.Upper.Value)
	}
}

// Valid returns false if the relation's name is empty, the ordinal variable
// has an empty value set, or the numerical variable has no limits.
func (r *Relation) Valid() bool {
	if len(r.ID) == 0 || len(r.Name) == 0 {
		return false
	}
	switch r.VariableType {
	case variables.Boolean, variables.Nominal, variables.Ordinal:
		if len(r.Value) == 0 {
			return false
		}
	case variables.Numerical:
		if r.Lower == nil && r.Upper == nil {
			return false
		}
	}
	return true
}

// apply applies the relation to the value range. It returns the values
// that agree with the relation.
func (r *Relation) apply(valueRange []string) []string {
	v := slice.ToIntSet(valueRange)
	if r.Lower != nil {
		b, _ := strconv.Atoi(r.Lower.Value)
		for a := range v {
			if a < b {
				delete(v, a)
			}
			if a == b && !r.Lower.Incl {
				delete(v, a)
			}
		}
	}
	if r.Upper != nil {
		b, _ := strconv.Atoi(r.Upper.Value)
		for a := range v {
			if a > b {
				delete(v, a)
			}
			if a == b && !r.Upper.Incl {
				delete(v, a)
			}
		}
	}
	return slice.IntSetToStringSlice(v)
}

// Negate negates the relation.
func (r *Relation) Negate(valueRange []string) {
	r.Lower, r.Upper = r.Upper, r.Lower
	if r.Lower != nil {
		r.Lower.Incl = !r.Lower.Incl
	}
	if r.Upper != nil {
		r.Upper.Incl = !r.Upper.Incl
	}
	if valueRange != nil && len(r.Value) > 0 {
		values := set.New(valueRange...)
		for _, a := range r.Value {
			values.Remove(a)
		}
		r.Value = slice.SetToSlice(values)
	}
}

// Transform transforms criteria relations by converting parsed values to strings of valid literals.
// If a valid literal cannot be inferred, the confidence score of the relation is set to zero.
// Indifferent nominal relations are removed by setting the confidence score to zero.
func (r *Relation) Transform() {
	variableCatalog := variables.Get()
	v := variableCatalog.Variable(r.ID)
	switch r.VariableType {
	case variables.Boolean:
		if r.ID != variables.Zero && text.IsYesNo(r.Value) {
			r.Score = 0
		}
	case variables.Numerical:
		if r.Lower != nil {
			if s, err := transform(v, r.Lower.Value); err == nil {
				r.Lower.Value = s
			} else {
				r.Score = 0
			}
		}
		if r.Upper != nil {
			if s, err := transform(v, r.Upper.Value); err == nil {
				r.Upper.Value = s
			} else {
				r.Score = 0
			}
		}
	}
}

// transform replaces the radix comma by dot, adds a missing zero (e.g., 150,00 -> 150,000),
// and removes the thousand commas. If the string value cannot be converted to a float literal,
// non-nil error is returned.
func transform(v *variables.Variable, s string) (string, error) {
	if reRadixComma.MatchString(s) {
		s = strings.Replace(s, ",", ".", 1)
	} else {
		if v.Name != "wbc" { // For wbc, 100,00 may mean 10,000.
			s = reMissingZero.ReplaceAllString(s, "000")
		}
		s = strings.Replace(s, ",", "", -1)
	}
	values := reTimes.Split(s, 2)
	if len(values) == 2 {
		s = values[0] + text.NormalizeScientificMultiplier(values[1])
	}
	val, err := strconv.ParseFloat(s, 64)
	if err == nil && !v.InRange(val) {
		err = fmt.Errorf("value %q not in valid range of variable: %s", s, v.Name)
	}
	return s, err
}

// Split splits the combination relation into individual relations. Because such relations may not have
// a valid ID, the split operation needs to be done before validation.
func (r *Relation) Split() Relations {
	variableCatalog := variables.Get()
	names := strings.Split(r.Name, "/")
	if len(names) != 2 {
		return Relations{r}
	}
	slice.TrimSpace(names)
	id0, ok0 := variableCatalog.ID(names[0])
	id1, ok1 := variableCatalog.ID(names[1])
	if !(ok0 && ok1) {
		return Relations{r}
	}
	r0 := &Relation{ID: id0, Name: names[0], Unit: r.Unit, VariableType: r.VariableType, Score: r.Score}
	r1 := &Relation{ID: id1, Name: names[1], Unit: r.Unit, VariableType: r.VariableType, Score: r.Score}
	if r.Lower != nil {
		values := strings.Split(r.Lower.Value, "/")
		slice.TrimSpace(values)
		switch len(values) {
		case 1:
			r0.Lower = &Limit{Incl: r.Lower.Incl, Value: values[0]}
			r1.Lower = &Limit{Incl: r.Lower.Incl, Value: values[0]}
		case 2:
			r0.Lower = &Limit{Incl: r.Lower.Incl, Value: values[0]}
			r1.Lower = &Limit{Incl: r.Lower.Incl, Value: values[1]}
		}
	}
	if r.Upper != nil {
		values := strings.Split(r.Upper.Value, "/")
		slice.TrimSpace(values)
		switch len(values) {
		case 1:
			r0.Upper = &Limit{Incl: r.Upper.Incl, Value: values[0]}
			r1.Upper = &Limit{Incl: r.Upper.Incl, Value: values[0]}
		case 2:
			r0.Upper = &Limit{Incl: r.Upper.Incl, Value: values[0]}
			r1.Upper = &Limit{Incl: r.Upper.Incl, Value: values[1]}
		}
	}
	return Relations{r0, r1}
}

// Less compares two numerical relations by their limits.
func (r *Relation) Less(q *Relation) bool {
	if r.ID != q.ID || r.VariableType != variables.Numerical {
		return false
	}
	var rval string
	switch {
	case r.Upper != nil:
		rval = r.Upper.Value
	case r.Lower != nil:
		rval = r.Lower.Value
	}
	var qval string
	switch {
	case q.Lower != nil:
		qval = q.Lower.Value
	case q.Upper != nil:
		qval = q.Upper.Value
	}
	return rval < qval
}

// JSON converts the the relations slice to the json string.
func (rs Relations) JSON() string {
	b, err := json.Marshal(rs)
	if err == nil {
		return string(b)
	}
	glog.Warningf("Failed to marshal relation: %v\n", err)
	return ""
}

// SetScore sets the confidence score that the relations are parsed correctly.
func (rs Relations) SetScore(score float64) {
	for _, r := range rs {
		r.SetScore(score)
	}
}

// MinScore returns the minimum score of the relations.
func (rs Relations) MinScore() float64 {
	if len(rs) == 0 {
		return 0
	}
	minScore := math.MaxFloat64
	for _, r := range rs {
		minScore = math.Min(minScore, r.Score)
	}
	return minScore
}

// Sort sorts relations by their id.
func (rs Relations) Sort() {
	sort.SliceStable(rs, func(i, j int) bool {
		if rs[i].ID == rs[j].ID {
			return rs[i].Less(rs[j])
		}
		return rs[i].ID < rs[j].ID
	})
}

// Dedupe removes the duplicate relations.
func (rs *Relations) Dedupe() {
	a := *rs
	if len(a) < 2 {
		return
	}
	a.Sort()
	for i := len(a) - 1; i > 0; i-- {
		if a[i-1].Name == a[i].Name {
			a = append(a[:i], a[i+1:]...)
		}
	}
	*rs = a
}

// setRelationFields sets the relation variable and unit fields.
func (rs Relations) setRelationFields() {
	variableCatalog := variables.Get()
	for _, r := range rs {
		if v := variableCatalog.Variable(r.ID); v != nil {
			r.SetVariableFields(v)
			if r.Unit == "" {
				r.SetUnitField(v)
			}
		}
	}
}

// split splits the multi-variable relation to the individual relations.
func (rs *Relations) split() {
	rels := NewRelations()
	for _, r := range *rs {
		rels = append(rels, r.Split()...)
	}
	*rs = rels
}

// Normalize normalizes the relation.
func (rs Relations) normalize() {
	variableCatalog := variables.Get()
	for _, r := range rs {
		if v := variableCatalog.Variable(r.ID); v != nil {
			r.Normalize(v.Range)
		}
	}
}

// Validate removes the invalid relations.
func (rs *Relations) validate() {
	a := *rs
	for i := len(a) - 1; i >= 0; i-- {
		if !a[i].Valid() {
			a = append(a[:i], a[i+1:]...)
		}
	}
	*rs = a
}

// Process splits the relations if needed, sets the correct types,
// normalizes and removes invalid relations.
func (rs *Relations) Process() {
	rs.split()
	rs.setRelationFields()
	rs.normalize()
	rs.validate()
	rs.Sort()
}

// Negate negates the relations.
func (rs Relations) Negate() {
	variableCatalog := variables.Get()
	for _, r := range rs {
		valueRange := variableCatalog.Variable(r.ID).Range
		r.Negate(valueRange)
	}
}

// Transform transforms the relations.
func (rs Relations) Transform() {
	for _, r := range rs {
		r.Transform()
	}
}

// VariableIDs returns the variable ids contained in the relations.
func (rs Relations) VariableIDs() []variables.ID {
	ids := make([]variables.ID, 0, len(rs))
	for _, r := range rs {
		ids = append(ids, r.ID)
	}
	return ids
}

// Empty returns true if the relations has no relation elements.
func (rs Relations) Empty() bool {
	return len(rs) == 0
}
