// Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.

package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/facebookresearch/Clinical-Trial-Parser/src/common/conf"
	"github.com/facebookresearch/Clinical-Trial-Parser/src/common/param"
	"github.com/facebookresearch/Clinical-Trial-Parser/src/common/util/fio"
	"github.com/facebookresearch/Clinical-Trial-Parser/src/common/util/timer"
	"github.com/facebookresearch/Clinical-Trial-Parser/src/ct/studies"
	"github.com/facebookresearch/Clinical-Trial-Parser/src/ct/units"
	"github.com/facebookresearch/Clinical-Trial-Parser/src/ct/variables"

	"github.com/golang/glog"
)

// main runs the CFG parser on clinical study eligibility criteria.
// The output is a file that contains parsed criteria in a JSON format.
// This example serves as a reference design for the CFG parser code.
func main() {
	p := NewParser()
	if err := p.LoadParameters(); err != nil {
		glog.Fatal(err)
	}
	if err := p.Initialize(); err != nil {
		glog.Fatal(err)
	}
	if err := p.Ingest(); err != nil {
		glog.Fatal(err)
	}
	p.Parse()
	p.Close()
}

// Parser defines the struct for processing eligibility criteria.
type Parser struct {
	parameters conf.Config
	registry   studies.Studies
	clock      timer.Timer
}

// NewParser creates a new parser to parse eligibility criteria.
func NewParser() *Parser {
	return &Parser{clock: timer.New()}
}

// LoadParameters loads parameters from command line and a config file.
func (p *Parser) LoadParameters() error {
	configFname := flag.String("conf", "", "Config file")
	inputFname := flag.String("i", "", "Input file")
	outputFname := flag.String("o", "", "Output file")

	flag.Parse()
	if len(*configFname) == 0 || len(*inputFname) == 0 || len(*outputFname) == 0 {
		return fmt.Errorf("usage: %s -conf <config file> -i <input file> -o <output name>", os.Args[0])
	}

	parameters, err := conf.Load(*configFname)
	if err != nil {
		return err
	}
	parameters.Put("input_file", *inputFname)
	parameters.Put("output_file", *outputFname)
	p.parameters = parameters

	return nil
}

// Initialize initializes the parser by loading the resource data.
func (p *Parser) Initialize() error {
	fname := p.parameters.GetResourcePath("variable_file")
	variableDictionary, err := variables.Load(fname)
	if err != nil {
		return err
	}
	variables.Set(variableDictionary)

	fname = p.parameters.GetResourcePath("unit_file")
	unitDictionary, err := units.Load(fname)
	if err != nil {
		return err
	}
	units.Set(unitDictionary)

	return nil
}

// Ingest ingests eligibility criteria from a file.
func (p *Parser) Ingest() error {
	fname := p.parameters.Get("input_file")
	f, err := os.Open(fname)
	if err != nil {
		return err
	}
	defer f.Close()

	registry := studies.New()
	r := csv.NewReader(f)
	r.Comment = rune(param.Comment)

	for {
		line, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if len(line) < 5 {
			return fmt.Errorf("too few columns, at least 4 needed: %v", line)
		}
		nctID := line[0]
		title := line[1]
		// Skip line[2]: has_us_facility
		conditions := strings.Split(line[3], param.FieldSep)
		eligibilityCriteria := line[4]

		study := studies.NewStudy(nctID, title, conditions, eligibilityCriteria)
		registry.Add(study)
	}
	glog.Infof("Ingested studies: %d\n", registry.Len())
	p.registry = registry

	return nil
}

// Parse parses the ingested eligibility criteria and writes the results to a file.
func (p *Parser) Parse() {
	header := "#nct_id\teligibility_type\tvariable_type\tcriterion_index\tcriterion\tquestion\trelation\n"
	criteriaCnt := 0
	parsedCriteriaCnt := 0
	relationCnt := 0
	fname := p.parameters.Get("output_file")
	writer := fio.Writer(fname)
	defer writer.Close()
	writer.WriteString(header)
	for _, study := range p.registry {
		writer.WriteString(study.Parse().Relations())
		criteriaCnt += study.CriteriaCount()
		parsedCriteriaCnt += study.ParsedCriteriaCount()
		relationCnt += study.RelationCount()

	}
	ratio := 0.0
	if criteriaCnt > 0 {
		ratio = 100 * float64(relationCnt) / float64(criteriaCnt)
	}
	glog.Infof("Ingested studies: %d, Extracted criteria: %d, Parsed criteria: %d, Relations: %d, Relations per criteria: %.1f%%\n",
		p.registry.Len(), criteriaCnt, parsedCriteriaCnt, relationCnt, ratio)
}

// Close closes the parser.
func (p *Parser) Close() {
	glog.Info(p.clock.Elapsed())
	glog.Flush()
}
