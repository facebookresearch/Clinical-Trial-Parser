// Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.

package parser

// Grammar defines the grammar interface to parse items into a parse tree.
type Grammar interface {
	BuildTrees(Items) Trees
}
