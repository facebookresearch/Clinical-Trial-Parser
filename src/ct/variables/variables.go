// Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.

package variables

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/facebookresearch/Clinical-Trial-Parser/src/common/param"
	"github.com/facebookresearch/Clinical-Trial-Parser/src/common/trie"
	"github.com/facebookresearch/Clinical-Trial-Parser/src/common/util/text"

	"github.com/golang/glog"
)

var catalog *Variables

func init() {
	catalog = DefaultCatalog()
}

func Set(d *Variables) {
	catalog = d
}

func Get() *Variables {
	return catalog
}

// Variables defines a collection of variables.
type Variables struct {
	ids        map[string]ID // Map from variable name to id.
	variables  map[ID]*Variable
	units      map[ID]string // default unit names
	questions  map[ID]string
	dictionary *trie.Trie
}

func New() *Variables {
	return &Variables{
		ids:        make(map[string]ID),
		variables:  make(map[ID]*Variable),
		units:      make(map[ID]string),
		questions:  make(map[ID]string),
		dictionary: trie.New(),
	}
}

func (vs *Variables) Size() int {
	return len(vs.variables)
}

// ID returns the ID of the variable name.
func (vs *Variables) ID(name string) (ID, bool) {
	id, ok := vs.ids[name]
	return id, ok
}

// Variable returns the variable associated with the variable id.
func (vs *Variables) Variable(id ID) *Variable {
	return vs.variables[id]
}

// Question returns the question associated with the variable id.
func (vs *Variables) Question(id ID) string {
	return vs.questions[id]
}

// Match returns true if the candidate is in the variable catalog.
func (vs *Variables) Match(candidate string) bool {
	return vs.dictionary.Match(candidate)
}

// Match returns corresponding name if the candidate is in the variable catalog.
func (vs *Variables) Get(candidate string) (string, bool) {
	if v, ok := vs.dictionary.Get(candidate); ok {
		return v.Name(), true
	}
	return "", false
}

func (vs *Variables) Add(id ID, kind Type, name string, display string, aliases []string, bounds []string, unitName string, question string) error {
	if _, ok := vs.variables[id]; ok {
		return fmt.Errorf("duplicate variable id: %s (name: %s)", id, name)
	}
	if _, ok := vs.ids[name]; ok {
		return fmt.Errorf("duplicate variable name: %s (id: %s)", name, id)
	}
	var numBounds []float64
	if kind == Numerical && len(bounds) == 2 {
		if low, err := strconv.ParseFloat(bounds[0], 64); err == nil {
			numBounds = append(numBounds, low)
		} else {
			return err
		}
		if high, err := strconv.ParseFloat(bounds[1], 64); err == nil {
			numBounds = append(numBounds, high)
		} else {
			return err
		}
		bounds = nil
	}
	v := NewVariable(id, kind, name, display, bounds, numBounds, unitName)
	vs.ids[name] = id
	vs.variables[id] = v
	vs.units[id] = unitName
	vs.questions[id] = question
	for _, a := range aliases {
		vals := text.CustomizeSlash(a)
		vs.dictionary.Put(name, vals...)
	}
	return nil
}

// Load loads variables from a file.​
func Load(fname string) (*Variables, error) {
	f, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	variables := New()
	r := csv.NewReader(f)
	r.Comment = rune(param.Comment)

	for {
		line, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("%s: %v", fname, err)
		}
		if len(line) < 8 {
			return nil, fmt.Errorf("%s: too few columns, at least 8 needed: %v", fname, line)
		}
		id := ID(line[0])
		kind := ParseType(line[1])
		name := line[2]
		display := line[3]
		aliases := strings.Split(line[4], param.FieldSep)
		bounds := strings.Split(line[5], param.FieldSep)
		defaultUnit := line[6]
		question := line[7]
		if err := variables.Add(id, kind, name, display, aliases, bounds, defaultUnit, question); err != nil {
			return nil, fmt.Errorf("%s: %v", fname, err)
		}
	}
	glog.Infof("Number of variables loaded: %d\n", variables.Size())

	return variables, nil
}

// DefaultCatalog defines the variable names and their aliases for unit testing.
// The display names are omitted for the convenience.
func DefaultCatalog() *Variables {
	catalog := New()
	var aliases []string

	aliases = []string{}
	catalog.Add(Zero, Numerical, "_", "", aliases, nil, "", "")

	aliases = []string{"nyha", "new york heart association"}
	catalog.Add("102", Ordinal, "nyha", "", aliases, []string{"1", "2", "3", "4"}, "", "")

	aliases = []string{"ecog", "eastern cooperative oncology group"}
	catalog.Add("100", Ordinal, "ecog", "", aliases, []string{"0", "1", "2", "3", "4"}, "", "")

	aliases = []string{"age", "ages", "aged"}
	catalog.Add("200", Numerical, "age", "", aliases, nil, "", "")

	aliases = []string{"height*"}
	catalog.Add("201", Numerical, "height", "", aliases, nil, "", "")

	aliases = []string{"weigh*", "body weigh*"}
	catalog.Add("202", Numerical, "weight", "", aliases, nil, "", "")

	aliases = []string{"bmi", "body mass index"}
	catalog.Add("203", Numerical, "bmi", "", aliases, nil, "", "")

	aliases = []string{"life expectancy"}
	catalog.Add("206", Numerical, "life_expectancy", "", aliases, nil, "", "")

	aliases = []string{"systolic blood pressure", "systolic", "sbp"}
	catalog.Add("300", Numerical, "sbp", "", aliases, nil, "", "")

	aliases = []string{"diastolic blood pressure", "diastolic", "dbp"}
	catalog.Add("301", Numerical, "dbp", "", aliases, nil, "", "")

	aliases = []string{"SBP/DBP", "blood pressure", "bp"}
	catalog.Add("302", Numerical, "sbp/dbp", "", aliases, nil, "", "")

	aliases = []string{"a1c", "hba1c", "hgba1c", "hemoglobin a1c"}
	catalog.Add("400", Numerical, "a1c", "", aliases, nil, "", "")

	aliases = []string{"hemoglobin count", "hb count"}
	catalog.Add("403", Numerical, "hb_count", "", aliases, nil, "", "")

	aliases = []string{"wbc", "white blood cell count", "white blood cell", "leukocytes", "leucocytes"}
	catalog.Add("404", Numerical, "wbc", "", aliases, nil, "", "")

	aliases = []string{"platelet count", "platelet"}
	catalog.Add("405", Numerical, "platelet_count", "", aliases, nil, "", "")

	aliases = []string{"absolute neutrophil count"}
	catalog.Add("408", Numerical, "anc", "", aliases, nil, "", "")

	aliases = []string{"aspartate aminotransferase", "ast", "sgot"}
	catalog.Add("411", Numerical, "ast", "", aliases, nil, "", "")

	aliases = []string{"alanine aminotransferase", "alt", "sgpt"}
	catalog.Add("412", Numerical, "alt", "", aliases, nil, "", "")

	aliases = []string{"ast/alt", "sgot/sgpt", "aspartate aminotransferase or alanine aminotransferase"}
	catalog.Add("413", Numerical, "ast/alt", "", aliases, nil, "", "")

	aliases = []string{"ast/alt ratio", "sgot/sgpt ratio"}
	catalog.Add("414", Numerical, "ast/alt_ratio", "", aliases, nil, "", "")

	aliases = []string{"plasma total cholesterol", "total cholesterol", "serum cholesterol", "cholesterol"}
	catalog.Add("500", Numerical, "total_cholesterol", "", aliases, nil, "", "")

	aliases = []string{"ldl", "ldl-cholesterol", "ldl cholesterol", "ldl-c", "low-density lipoprotein cholesterol"}
	catalog.Add("501", Numerical, "ldl_cholesterol", "", aliases, nil, "", "")

	aliases = []string{"fasting triglyceride level*", "fasting triglyceride*", "fasting plasma triglyceride*", "fasting serum triglyceride*"}
	catalog.Add("505", Numerical, "fasting_triglyceride_level", "", aliases, nil, "", "")

	aliases = []string{"triglyceride level*", "triglyceride*", "plasma triglyceride*", "serum triglyceride*"}
	catalog.Add("506", Numerical, "triglyceride_level", "", aliases, nil, "", "")

	aliases = []string{"karnofsky", "karnofsky performance score", "lansky", "karnofsky score", "kps"}
	catalog.Add("600", Numerical, "karnofsky_score", "", aliases, nil, "", "")

	aliases = []string{"p/f ratio", "pao2/fio2", "pao2/fio2 ratio"}
	catalog.Add("904", Numerical, "pf_ratio", "", aliases, nil, "mmhg", "")

	return catalog
}
