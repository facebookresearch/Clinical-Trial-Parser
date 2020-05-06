// Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.

package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/facebookresearch/Clinical-Trial-Parser/src/common/col/set"
	"github.com/facebookresearch/Clinical-Trial-Parser/src/common/conf"
	"github.com/facebookresearch/Clinical-Trial-Parser/src/common/param"
	"github.com/facebookresearch/Clinical-Trial-Parser/src/common/util/fio"
	"github.com/facebookresearch/Clinical-Trial-Parser/src/common/util/slice"
	"github.com/facebookresearch/Clinical-Trial-Parser/src/common/util/timer"
	"github.com/facebookresearch/Clinical-Trial-Parser/src/vocabularies"
	"github.com/facebookresearch/Clinical-Trial-Parser/src/vocabularies/mesh"
	"github.com/facebookresearch/Clinical-Trial-Parser/src/vocabularies/taxonomy"
	"github.com/facebookresearch/Clinical-Trial-Parser/src/vocabularies/umls"

	"github.com/golang/glog"
)

var (
	reParentheses = regexp.MustCompile(`\([^)]*(\)|$)`)
	reConjunction = regexp.MustCompile(` and | or |,`)
)

// main matches (grounds) extracted (input) terms to vocabulary concepts.
// Matching results are written to a file. This version does not use clustering.
func main() {
	m := NewMatcher()
	if err := m.LoadParameters(); err != nil {
		glog.Fatal(err)
	}
	if err := m.LoadVocabulary(); err != nil {
		glog.Fatal(err)
	}
	if err := m.Match(); err != nil {
		glog.Fatal(err)
	}
	m.Close()
}

// Slot defines the extracted NER slot.
type Slot struct {
	label string  // Slot label
	term  string  // Slot term
	score float64 // NER score
}

type Slots []Slot

func NewSlot(label string, term string, score float64) Slot {
	return Slot{label: label, term: term, score: score}
}

func (s Slot) SubTerms() []string {
	v := reConjunction.Split(s.term, -1)
	slice.TrimSpace(v)
	return slice.RemoveEmpty(v)
}

func (s *Slot) Normalize(normalize taxonomy.Normalizer) {
	_, s.term = normalize(s.term)
}

func (s Slot) String() string {
	return fmt.Sprintf("%s\t%s\t%.3f", s.label, s.term, s.score)
}

func NewSlots() Slots {
	return make(Slots, 0)
}

func (ss *Slots) Add(label string, term string, score float64) {
	s := NewSlot(label, term, score)
	*ss = append(*ss, s)
}

func (ss Slots) Size() int {
	return len(ss)
}

// Matcher defines the struct that matches extracted terms to concepts
// from a vocabulary.
type Matcher struct {
	parameters conf.Config
	vocabulary *taxonomy.Taxonomy
	normalize  taxonomy.Normalizer
	clock      timer.Timer
}

// NewMatcher creates a new matcher.
func NewMatcher() *Matcher {
	clock := timer.New()
	return &Matcher{clock: clock}
}

// LoadParameters loads parameters from command line and a config file.
func (m *Matcher) LoadParameters() error {
	configFname := flag.String("conf", "", "Config file")
	inputFname := flag.String("i", "", "Input file")
	outputFname := flag.String("o", "", "Output file")

	flag.Parse()
	if len(*configFname) == 0 {
		return fmt.Errorf("usage: %s -conf <config file> -i <input file> -o <output file>", os.Args[0])
	}

	parameters, err := conf.Load(*configFname)
	if err != nil {
		return err
	}

	if len(*inputFname) > 0 {
		parameters.Put("input_file", *inputFname)
	}
	if len(*outputFname) > 0 {
		parameters.Put("output_file", *outputFname)
	}
	if !parameters.Exists("input_file") {
		return fmt.Errorf("input file not defined")
	}
	if !parameters.Exists("output_file") {
		return fmt.Errorf("output file not defined")
	}

	m.parameters = parameters

	return nil
}

func (m *Matcher) LoadVocabulary() error {
	vocabularyFname := m.parameters.Get("vocabulary_file")
	var customFnames []string
	if m.parameters.Exists("custom_vocabulary_file") {
		path := m.parameters.Get("custom_vocabulary_file")
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

	m.normalize = mesh.Normalize
	vocabulary.Normalize(m.normalize)
	vocabulary.SetHashIndex(rows, bands)
	vocabulary.Info()

	m.vocabulary = vocabulary

	return nil
}

// getNERSlots gets the extracted terms from a string.
func getNERSlots(termStr string, nerThreshold float64, validLabels set.Set) Slots {
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(termStr), &data); err != nil {
		glog.Fatal(termStr, err)
	}
	slots := NewSlots()
	for label, values := range data {
		if validLabels.Contains(label) {
			for _, fields := range values.([]interface{}) {
				var term string
				var score float64
				for _, f := range fields.([]interface{}) {
					switch f.(type) {
					case string:
						term = f.(string)
					case float64:
						score = f.(float64)
					default:
						glog.Fatalf("unknown type: %v", f)
					}
				}
				norm := reParentheses.ReplaceAllString(term, " ")
				norm = strings.TrimSpace(norm)
				if len(norm) > 0 {
					term = norm
				}
				if score > nerThreshold && len(term) > 0 {
					slots.Add(label, term, score)
				}
			}
		}
	}
	return slots
}

func (m *Matcher) Match() error {
	nerThreshold := m.parameters.GetFloat64("ner_threshold")
	validLabels := set.New(m.parameters.GetSlice("valid_labels", ",")...)

	matchThreshold := m.parameters.GetFloat64("match_threshold")
	matchMargin := m.parameters.GetFloat64("match_margin")

	matchedSlots := make(map[string]taxonomy.Terms)
	conceptSet := set.New()
	slotCnt := 0
	matchedSlotCnt := 0

	defaultCategories := set.New()
	cancerCategories := set.New("C")

	fname := m.parameters.Get("input_file")
	file, err := os.Open(fname)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineCnt := 0

	outputFname := m.parameters.Get("output_file")
	writer := fio.Writer(outputFname)
	defer writer.Close()

	header := "#nct_id\teligibility_type\tcriterion\tlabel\tterm\tner_score\tconcepts\ttree_numbers\tnel_score\n"
	writer.WriteString(header)

	glog.Infof("Matching NER terms ...")

	for scanner.Scan() {
		lineCnt++
		line := scanner.Text()
		if len(line) == 0 || line[0] == param.Comment {
			continue
		}

		// Extract NER terms
		values := strings.Split(line, "\t")
		nctID := values[0]
		eligibilityType := values[1]
		criterion := values[2]
		termStr := values[3]
		slots := getNERSlots(termStr, nerThreshold, validLabels)
		slotCnt += slots.Size()

		// Match NER terms to concepts
		for _, slot := range slots {
			subterms := slot.SubTerms()
			for _, subterm := range subterms {
				if _, ok := matchedSlots[subterm]; !ok {
					validCategories := defaultCategories
					if slot.label == "word_scores:cancer" {
						validCategories = cancerCategories
					}
					matchedSlots[subterm] = m.vocabulary.Match(subterm, matchMargin, validCategories)
				}
			}

			slot.Normalize(m.normalize)
			hasMatch := false

			for _, subterm := range subterms {
				matchedConcepts := matchedSlots[subterm]
				if matchedConcepts.MaxValue() >= matchThreshold {
					hasMatch = true
					conceptSet.Add(matchedConcepts.Keys()...)
					concepts := strings.Join(matchedConcepts.Keys(), "|")
					nelScore := matchedConcepts.MaxValue()
					treeNumbers := strings.Join(matchedConcepts.TreeNumbers(), "|")
					if _, err := fmt.Fprintf(writer, "%s\t%s\t%s\t%s\t%s\t%s\t%.3f\n", nctID, eligibilityType, criterion, slot.String(), concepts, treeNumbers, nelScore); err != nil {
						return err
					}
				}
			}

			if hasMatch {
				matchedSlotCnt++
			} else {
				if _, err := fmt.Fprintf(writer, "%s\t%s\t%s\t%s\n", nctID, eligibilityType, criterion, slot.String()); err != nil {
					return err
				}
			}
		}
	}

	glog.Infof("Lines read: %d, Slots: %d, Unique slots: %d\n", lineCnt, slotCnt, len(matchedSlots))
	glog.Infof("%d slots matched to %d concepts\n", matchedSlotCnt, conceptSet.Size())
	glog.Infof("%d slots not matched\n", slotCnt-matchedSlotCnt)

	return nil
}

// Close closes the matcher.
func (m *Matcher) Close() {
	glog.Info(m.clock.Elapsed())
	glog.Flush()
}
