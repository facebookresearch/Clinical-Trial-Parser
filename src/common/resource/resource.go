// Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.

package resource

import (
	"os"
	"path"
	"strings"

	"github.com/golang/glog"
)

var (
	// LocalResourcePath is the relative path from the base of the repo
	// to the resource directory and will be used in case this repo is
	// not installed within the system's "GOPATH".
	LocalResourcePath = "src/resources"

	// ResourcePath is the default path to the resource directory.
	ResourcePath = "/github.com/facebookresearch/Clinical-Trial-Parser/" + LocalResourcePath

	// LocalDataPath is the relative path from the base of the repo
	// to the data directory and will be used in case this repo is
	// not installed within the system's "GOPATH".
	LocalDataPath = "data"

	// DataPath is the default path to the data directory.
	DataPath = "/github.com/facebookresearch/Clinical-Trial-Parser/" + LocalDataPath
)

// GetResourcePath returns the path to the project's resource directory.
func GetResourcePath() string {
	if env := os.Getenv("RESOURCE_PATH"); len(env) != 0 {
		return env
	}

	if _, err := os.Stat(LocalResourcePath); err == nil {
		return LocalResourcePath
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

	if _, err := os.Stat(LocalDataPath); err == nil {
		return LocalDataPath
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
