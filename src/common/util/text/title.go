// Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.

package text

import (
	"strings"
)

// Titles capitalizes the first letter of each word.
func Titles(s []string) []string {
	t := make([]string, len(s))
	for i, v := range s {
		t[i] = strings.Title(v)
	}
	return t
}
