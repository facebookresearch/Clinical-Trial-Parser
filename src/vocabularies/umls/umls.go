// Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.

package umls

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/facebookresearch/Clinical-Trial-Parser/src/common/col/set"
	"github.com/facebookresearch/Clinical-Trial-Parser/src/vocabularies/taxonomy"

	"github.com/golang/glog"
)

// Load loads a UMLS vocabulary from MRCONSO.RRF.
func Load(fname string) *taxonomy.Taxonomy {
	file, err := os.Open(fname)
	if err != nil {
		glog.Fatal(err)
	}
	defer file.Close()

	ids := set.New()
	root := taxonomy.NewNode("root")
	var de *taxonomy.Node

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}

		values := strings.Split(line, "|")
		if len(values) < 17 {
			fmt.Printf("Wrong number of columns; expected %d: %s\n", 17, line)
			continue
		}
		lang := strings.TrimSpace(values[1])
		vocabularly := strings.TrimSpace(values[11])
		if lang != "ENG" || vocabularly != "SNOMEDCT_US" && vocabularly != "MSH" {
			continue
		}

		id := strings.TrimSpace(values[0])
		name := strings.TrimSpace(values[14])

		if ids[id] {
			de.AddSynonym(name)
		} else {
			ids.Add(id)
			de = taxonomy.NewNode(name)
			de.AddTreeNumber(id)
			de.AddSynonym(name)
			root.AddChild(de)
		}
	}

	t := taxonomy.New(root)
	t.SetBaseIndex()

	return t
}
