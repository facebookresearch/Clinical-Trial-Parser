// Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.

package taxonomy

// Normalizer defines a normalizer function to normalize node synonyms.
type Normalizer func(string) (string, string)

var identity Normalizer = func(s string) (string, string) { return s, s }
