// Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.

package units

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/facebookresearch/Clinical-Trial-Parser/src/common/param"
	"github.com/facebookresearch/Clinical-Trial-Parser/src/common/trie"
	"github.com/facebookresearch/Clinical-Trial-Parser/src/common/util/text"

	"github.com/golang/glog"
)

var catalog *Units

func init() {
	catalog = DefaultCatalog()
}

func Set(d *Units) {
	catalog = d
}

func Get() *Units {
	return catalog
}

type Units struct {
	ids        map[string]ID // map from unit name to unit id.
	units      map[ID]*Unit
	variables  map[string]string
	dictionary *trie.Trie
}

func New() *Units {
	return &Units{
		ids:        make(map[string]ID),
		units:      make(map[ID]*Unit),
		variables:  make(map[string]string),
		dictionary: trie.New(),
	}
}

func (us *Units) Size() int {
	return len(us.units)
}

// Variable returns the variable associated with the variable id.
func (us *Units) Unit(id ID) *Unit {
	return us.units[id]
}

//ID returns the ID of the variable name.
func (us *Units) ID(name string) (ID, bool) {
	id, ok := us.ids[name]
	return id, ok
}

func (us *Units) Variable(name string) (string, bool) {
	id, ok := us.variables[name]
	return id, ok
}

// Match returns true if the candidate is in the unit catalog.
func (us *Units) Match(candidate string) bool {
	return us.dictionary.Match(candidate)
}

// Match returns true if the candidate is in the unit catalog.
func (us *Units) Get(candidate string) (string, bool) {
	if v, ok := us.dictionary.Get(candidate); ok {
		return v.Name(), true
	}
	return "", false
}

func (us *Units) Add(id ID, name string, display string, aliases []string, vname string) error {
	if _, ok := us.units[id]; ok {
		return fmt.Errorf("duplicate unit id: %s (name: %s)", id, name)
	}
	if _, ok := us.ids[name]; ok {
		return fmt.Errorf("duplicate unit name: %s (id: %s)", name, id)
	}
	u := NewUnit(id, name, display, vname)
	us.ids[name] = id
	us.units[id] = u
	if len(vname) > 0 {
		us.variables[name] = vname
	}
	for _, a := range aliases {
		vals := text.CustomizeSlash(a)
		us.dictionary.Put(name, vals...)
	}
	return nil
}

// Load loads units from a file.​
func Load(fname string) (*Units, error) {
	f, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	units := New()
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
		if len(line) < 5 {
			return nil, fmt.Errorf("too few columns, at least 5 needed: %v", line)
		}
		id := ID(line[0])
		name := line[1]
		display := line[2]
		aliases := strings.Split(line[3], param.FieldSep)
		vname := line[4]
		if err := units.Add(id, name, display, aliases, vname); err != nil {
			return nil, fmt.Errorf("%s: %v", fname, err)
		}
	}
	glog.Infof("Number of units loaded: %d\n", units.Size())

	return units, nil
}

// DefaultCatalog defines variable units and their aliases.
func DefaultCatalog() *Units {
	catalog := New()
	var aliases []string

	aliases = []string{"%"}
	catalog.Add("100", "%", "%", aliases, "")

	aliases = []string{"kg", "kilograms"}
	catalog.Add("200", "kg", "kg", aliases, "weight")

	aliases = []string{"g", "grams"}
	catalog.Add("201", "g", "g", aliases, "")

	aliases = []string{"mg"}
	catalog.Add("202", "mg", "mg", aliases, "")

	aliases = []string{"lb", "lbs", "pound", "pounds"}
	catalog.Add("203", "lb", "pound", aliases, "weight")

	aliases = []string{"day*"}
	catalog.Add("303", "day", "day", aliases, "")

	aliases = []string{"week*"}
	catalog.Add("304", "week", "week", aliases, "")

	aliases = []string{"month"}
	catalog.Add("305", "month", "month*", aliases, "")

	aliases = []string{"year*"}
	catalog.Add("306", "year", "year", aliases, "")

	aliases = []string{"ml/min"}
	catalog.Add("400", "ml/min", "ml/min", aliases, "")

	aliases = []string{"g/day"}
	catalog.Add("401", "g/day", "g/day", aliases, "")

	aliases = []string{"g/dl"}
	catalog.Add("403", "g/dl", "g/dl", aliases, "")

	aliases = []string{"ng/dl"}
	catalog.Add("404", "ng/dl", "ng/dl", aliases, "")

	aliases = []string{"ng/ml"}
	catalog.Add("405", "ng/ml", "ng/ml", aliases, "")

	aliases = []string{"mg/dl"}
	catalog.Add("407", "mg/dl", "mg/dl", aliases, "")

	aliases = []string{"cells/ul", "/ul", "mm3"}
	catalog.Add("410", "cells/ul", "cells/ul", aliases, "")

	aliases = []string{"ml/min/1"}
	catalog.Add("414", "mL/min/1.73_m2", "mL/min/1.73 m2", aliases, "")

	aliases = []string{"cells/l", "/l"}
	catalog.Add("416", "cells/l", "cells/L", aliases, "")

	aliases = []string{"cm"}
	catalog.Add("501", "cm", "cm", aliases, "")

	aliases = []string{"m"}
	catalog.Add("502", "m", "m", aliases, "")

	aliases = []string{"mmhg"}
	catalog.Add("600", "mmhg", "", aliases, "")

	aliases = []string{"kg/m2", "kg/m^2", "kg/m²", "kilogram per meter square"}
	catalog.Add("602", "kg/m2", "kg/m2", aliases, "bmi")

	aliases = []string{"uln", "upper limit of normal", "upper limits of normal", "laboratory normal"}
	catalog.Add("603", "uln", "uln", aliases, "")

	aliases = []string{"lln", "lower limit of normal", "lower limits of normal"}
	catalog.Add("604", "lln", "lln", aliases, "")

	return catalog
}
