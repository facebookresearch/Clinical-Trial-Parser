// Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.

package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/facebookresearch/Clinical-Trial-Parser/src/common/col/set"
	"github.com/facebookresearch/Clinical-Trial-Parser/src/common/conf"
	"github.com/facebookresearch/Clinical-Trial-Parser/src/common/util/fio"
	"github.com/facebookresearch/Clinical-Trial-Parser/src/vocabularies"
	"github.com/facebookresearch/Clinical-Trial-Parser/src/vocabularies/mesh"
	"github.com/facebookresearch/Clinical-Trial-Parser/src/vocabularies/taxonomy"
	"github.com/facebookresearch/Clinical-Trial-Parser/src/vocabularies/umls"

	"github.com/golang/glog"
)

// main searches for matching concepts from a vocabulary.
// Search strings are entered from console or a file.
func main() {
	m := NewMatcher()
	if err := m.LoadParameters(); err != nil {
		glog.Fatal(err)
	}
	if err := m.LoadVocabulary(); err != nil {
		glog.Fatal(err)
	}
	m.Search()
}

// Matcher defines the struct that matches strings to concepts from a vocabulary.
type Matcher struct {
	parameters conf.Config
	vocabulary *taxonomy.Taxonomy
}

// NewMatcher creates a new matcher.
func NewMatcher() *Matcher {
	return &Matcher{}
}

// LoadParameters loads parameters from command line and a config file.
func (m *Matcher) LoadParameters() error {
	configFname := flag.String("conf", "", "Config file")
	inputFname := flag.String("i", "", "Input file")

	flag.Parse()
	if len(*configFname) == 0 {
		return fmt.Errorf("usage: %s -conf <config file>", os.Args[0])
	}

	parameters, err := conf.Load(*configFname)
	if err != nil {
		return err
	}
	if len(*inputFname) > 0 {
		parameters.Put("input_file", *inputFname)
	}
	m.parameters = parameters
	return nil
}

// LoadVocabulary loads a vocabulary from a file.
func (m *Matcher) LoadVocabulary() error {
	vocabularyFname := m.parameters.GetDataPath("vocabulary_file")
	var customFnames []string
	if m.parameters.Exists("custom_vocabulary_file") {
		path := m.parameters.GetDataPath("custom_vocabulary_file")
		customFnames = fio.ReadFnames(path)
	}

	source := vocabularies.ParseSource(m.parameters.Get("vocabulary_source"))
	var vocabulary *taxonomy.Taxonomy
	switch source {
	case vocabularies.MESH:
		glog.Info("Loading MeSH ...")
		vocabulary = mesh.Load(vocabularyFname, customFnames...)
	case vocabularies.UMLS:
		glog.Info("Loading UMLS ...")
		vocabulary = umls.Load(vocabularyFname)
	default:
		return fmt.Errorf("unknown vocabulary source")
	}

	rows := m.parameters.GetInt("lsh_rows")
	bands := m.parameters.GetInt("lsh_rows")

	vocabulary.Normalize(mesh.Normalize)
	vocabulary.SetHashIndex(rows, bands)
	vocabulary.Info()

	m.vocabulary = vocabulary

	return nil
}

// Search searches matching concepts.
func (m *Matcher) Search() {
	if m.parameters.Exists("input_file") {
		m.batchSearch()
	} else {
		m.consoleSearch()
	}
}

// consoleSearch matches terms from stdin to concepts.
func (m *Matcher) consoleSearch() {
	emptyCategories := set.New()
	matchMargin := 1.0

	reader := bufio.NewReader(os.Stdin)
	getSearchStr := func() string {
		answer, err := reader.ReadString('\n')
		if err != nil {
			glog.Warning(err)
			return ""
		}
		return strings.TrimSuffix(answer, "\n")
	}

	fmt.Println("Enter search string ('q' to quit):")

	for {
		switch s := getSearchStr(); s {
		case "q":
			return
		case "":
		// skip
		default:
			matches := m.vocabulary.Match(s, matchMargin, emptyCategories)
			fmt.Println(matches.String())
			fmt.Println()
		}
	}
}

// batchSearch matches terms from a file to concepts.
func (m *Matcher) batchSearch() {
	emptyCategories := set.New()
	matchMargin := 1.0

	fname := m.parameters.Get("input_file")
	nodes := taxonomy.LoadNodes(fname)
	for _, node := range nodes {
		matches := m.vocabulary.MatchNode(node, matchMargin, emptyCategories)
		fmt.Printf("%s:\n%s\n\n", node.Name(), matches.String())
	}
}
