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

	"github.com/golang/glog"
)

// main extracts inclusion and exclusion criteria from eligibility criteria.
// The output is a file that contains a line per inclusion or exclusion criterion.
func main() {
	p := NewExtractor()
	if err := p.LoadParameters(); err != nil {
		glog.Fatal(err)
	}
	if err := p.Ingest(); err != nil {
		glog.Fatal(err)
	}
	if err := p.Extract(); err != nil {
		glog.Fatal(err)
	}
	p.Close()
}

// Extractor defines the struct for extracting inclusion and exclusion criteria.
type Extractor struct {
	parameters conf.Config
	registry   studies.Studies
	clock      timer.Timer
}

// NewExtractor creates a new extractor.
func NewExtractor() *Extractor {
	return &Extractor{clock: timer.New()}
}

// LoadParameters loads parameters from command line.
func (p *Extractor) LoadParameters() error {
	inputFname := flag.String("i", "", "Input file")
	outputFname := flag.String("o", "", "Output file")

	flag.Parse()
	if len(*inputFname) == 0 || len(*outputFname) == 0 {
		return fmt.Errorf("usage: %s -conf <config file> -i <input file> -o <output name>", os.Args[0])
	}

	parameters := conf.New()
	parameters.Put("input_file", *inputFname)
	parameters.Put("output_file", *outputFname)
	p.parameters = parameters

	return nil
}

// Ingest ingests eligibility criteria from a file.
func (p *Extractor) Ingest() error {
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

// Extract extracts inclusion and exclusion criteria and writes them to a file.
func (p *Extractor) Extract() error {
	header := "#nct_id\teligibility_type\tcriterion\n"
	fname := p.parameters.Get("output_file")
	writer := fio.Writer(fname)
	defer writer.Close()

	criteriaCnt := 0
	writer.WriteString(header)
	for _, study := range p.registry {
		inclusions, exclusions := study.Criteria()
		for _, criterion := range inclusions {
			if _, err := fmt.Fprintf(writer, "%s\t%s\t%s\n", study.NCT(), "inclusion", criterion); err != nil {
				return err
			}
		}
		for _, criterion := range exclusions {
			if _, err := fmt.Fprintf(writer, "%s\t%s\t%s\n", study.NCT(), "exclusion", criterion); err != nil {
				return err
			}
		}
		criteriaCnt += len(inclusions) + len(exclusions)
	}
	glog.Infof("Ingested studies: %d, Extracted criteria: %d\n", p.registry.Len(), criteriaCnt)
	return nil
}

// Close closes the parser.
func (p *Extractor) Close() {
	glog.Info(p.clock.Elapsed())
	glog.Flush()
}
