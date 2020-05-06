// Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.

package resource

import (
	"os"
	"path"
	"strings"

	"github.com/golang/glog"
)

var (
	// ResourcePath is the default path to the resource directory.
	ResourcePath = "/github.com/facebookresearch/Clinical-Trial-Parser/src/resources"

	// DataPath is the default path to the data directory.
	DataPath = "/github.com/facebookresearch/Clinical-Trial-Parser/data"
)

// GetResourcePath returns the path to the project's resource directory.
func GetResourcePath() string {
	if env := os.Getenv("RESOURCE_PATH"); len(env) != 0 {
		return env
	}

	if env := os.Getenv("GOPATH"); env != "" {
		gopaths := strings.Split(env, ":")
		for _, gp := range gopaths {
			check := path.Join(gp, "src", ResourcePath)
			if _, err := os.Stat(check); err == nil {
				return check
			}
		}
	}

	glog.Fatalf("Cannot find resource path for %q\n", ResourcePath)
	return ""
}

// GetDataPath returns the path to the project's data directory.
func GetDataPath() string {
	if env := os.Getenv("DATA_PATH"); len(env) != 0 {
		return env
	}

	if env := os.Getenv("GOPATH"); env != "" {
		gopaths := strings.Split(env, ":")
		for _, gp := range gopaths {
			check := path.Join(gp, "src", DataPath)
			if _, err := os.Stat(check); err == nil {
				return check
			}
		}
	}

	glog.Fatalf("Cannot find data path for %q\n", DataPath)
	return ""
}
