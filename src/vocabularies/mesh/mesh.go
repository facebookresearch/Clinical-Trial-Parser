// Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.

package mesh

import (
	"encoding/xml"
	"io/ioutil"
	"os"

	"github.com/facebookresearch/Clinical-Trial-Parser/src/vocabularies/taxonomy"

	"github.com/golang/glog"
)

// Descriptors defines the xml struct for Descriptors.
type Descriptors struct {
	XMLName     xml.Name     `xml:"DescriptorRecordSet"`
	Descriptors []Descriptor `xml:"DescriptorRecord"`
}

// Descriptor defines the xml struct for Descriptor.
type Descriptor struct {
	XMLName     xml.Name       `xml:"DescriptorRecord"`
	Name        DescriptorName `xml:"DescriptorName"`
	Concepts    Concepts       `xml:"ConceptList"`
	TreeNumbers TreeNumbers    `xml:"TreeNumberList"`
}

// Concepts defines the xml struct for Concepts.
type Concepts struct {
	XMLName  xml.Name  `xml:"ConceptList"`
	Concepts []Concept `xml:"Concept"`
}

// Concept defines the xml struct for Concept.
type Concept struct {
	XMLName xml.Name    `xml:"Concept"`
	Name    ConceptName `xml:"ConceptName"`
	Terms   Terms       `xml:"TermList"`
}

// TreeNumbers defines the xml struct for TreeNumbers.
type TreeNumbers struct {
	XMLName     xml.Name `xml:"TreeNumberList"`
	TreeNumbers []string `xml:"TreeNumber"`
}

// Terms defines the xml struct for Terms.
type Terms struct {
	XMLName xml.Name `xml:"TermList"`
	Terms   []Term   `xml:"Term"`
}

// Term defines the xml struct for Term.
type Term struct {
	XMLName   xml.Name `xml:"Term"`
	Preferred string   `xml:"ConceptPreferredTermYN,attr"`
	Name      string   `xml:"String"`
	ID        string   `xml:"TermUI"`
}

// IsPreferred indicates whether the term is a preferred term for a concept.
func (t Term) IsPreferred() bool {
	return t.Preferred == "Y"
}

// DescriptorName defines the xml struct for DescriptorName.
type DescriptorName struct {
	XMLName xml.Name `xml:"DescriptorName"`
	Value   string   `xml:"String"`
}

// ConceptName defines the xml struct for ConceptName.
type ConceptName struct {
	XMLName xml.Name `xml:"ConceptName"`
	Value   string   `xml:"String"`
}

// Load loads a MeSH taxonomy from files.
func Load(xmlFname string, customFnames ...string) *taxonomy.Taxonomy {
	t := loadTaxonomy(xmlFname)
	if len(customFnames) > 0 {
		nodes := taxonomy.LoadNodes(customFnames...)
		cnt := t.AddNodes(nodes)
		glog.Infof("%v: Nodes read: %d, New nodes: %d\n", customFnames, nodes.Len(), cnt)
	}

	t.SetBaseIndex()

	return t
}

// loadTaxonomy loads a MeSH vocabulary from an xml dump.
func loadTaxonomy(fname string) *taxonomy.Taxonomy {
	file, err := os.Open(fname)
	if err != nil {
		glog.Fatal(err)
	}
	defer file.Close()

	byteValue, _ := ioutil.ReadAll(file)

	var descriptors Descriptors
	if err := xml.Unmarshal(byteValue, &descriptors); err != nil {
		glog.Fatal(err)
	}

	root := taxonomy.NewNode("root")

	for _, d := range descriptors.Descriptors {
		treeNumbers := d.TreeNumbers.TreeNumbers
		if HasAnimalCode(treeNumbers) {
			continue
		}
		treeNumbers = Trim(treeNumbers)
		if len(treeNumbers) == 0 {
			continue
		}
		de := taxonomy.NewNode(d.Name.Value)
		for _, c := range d.Concepts.Concepts {
			if !isAnimalConcept(c.Name.Value) {
				ce := taxonomy.NewNode(c.Name.Value)
				for _, t := range c.Terms.Terms {
					ce.AddSynonym(t.Name)
				}
				ce.AddSynonym(c.Name.Value)
				ce.AddTreeNumber(treeNumbers...)
				de.AddChild(ce)
			}
		}
		root.AddChild(de)
	}

	return taxonomy.New(root)
}
