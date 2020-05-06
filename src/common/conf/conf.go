// Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.

package conf

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/facebookresearch/Clinical-Trial-Parser/src/common/param"
	"github.com/facebookresearch/Clinical-Trial-Parser/src/common/resource"
	"github.com/facebookresearch/Clinical-Trial-Parser/src/common/util/slice"
	"github.com/golang/glog"
)

// Config defines a configuration container for storing parameters.
type Config map[string]string

// New creates a new config container.
func New() Config {
	return make(map[string]string)
}

// Put sets the 'key' parameter to the string 'value'.
func (c Config) Put(key, value string) {
	c[key] = value
}

// Load loads the config container from the file.
func Load(fname string) (Config, error) {
	c := New()
	lineCnt := 0

	replaceVar := func(k string) string {
		if v := os.Getenv(k); len(v) > 0 {
			return v
		}
		if v, ok := c[k]; ok {
			return v
		}
		return ""
	}

	file, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lineCnt++
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if len(line) > 0 && line[0] != param.Comment {
			values := strings.SplitN(line, "=", 2)
			if len(values) != 2 {
				return nil, fmt.Errorf("bad config line %d: '%s'\n", lineCnt, line)
			}
			key := strings.TrimSpace(values[0])
			value := strings.TrimSpace(values[1])
			if len(key) == 0 {
				return nil, fmt.Errorf("bad key on line %d: '%s'\n", lineCnt, line)
			}
			value = os.Expand(value, replaceVar)
			c.Put(key, value)
		}
	}
	glog.Infof("%s: Lines: %d, Parameters: %d\n", fname, lineCnt, c.Size())
	return c, nil
}

// Size returns the number of parameters in the config container.
func (c Config) Size() int {
	return len(c)
}

// Exists returns true if the 'key' parameter is defined.
func (c Config) Exists(key string) bool {
	_, ok := c[key]
	return ok
}

// Get returns the string 'value' of the 'key' parameter.
func (c Config) Get(key string) string {
	value, ok := c[key]
	if !ok {
		glog.Fatalf("Error reading config key: %s\n", key)
	}
	return value
}

// GetResourcePath returns the path of the 'key' parameter.
func (c Config) GetResourcePath(key string) string {
	value := c.Get(key)
	if path.IsAbs(value) {
		return value
	}
	return path.Join(resource.GetResourcePath(), value)
}

// GetDataPath returns the path of the 'key' parameter.
func (c Config) GetDataPath(key string) string {
	value := c.Get(key)
	if path.IsAbs(value) {
		return value
	}
	return path.Join(resource.GetDataPath(), value)
}

// GetBool returns the boolean value of the 'key' parameter.
func (c Config) GetBool(key string) bool {
	value := c.Get(key)
	ind, err := strconv.ParseBool(value)
	if err != nil {
		glog.Fatalf("Error reading config key: %s; %v\n", key, err)
	}

	return ind
}

// GetInt returns the int value of the 'key' parameter.
func (c Config) GetInt(key string) int {
	value := c.Get(key)
	number, err := strconv.Atoi(value)
	if err != nil {
		glog.Fatalf("Error reading config key: %s; %v\n", key, err)
	}
	return number
}

// GetFloat64 returns the float64 value of the 'key' parameter.
func (c *Config) GetFloat64(key string) float64 {
	value := c.Get(key)
	number, err := strconv.ParseFloat(value, 64)
	if err != nil {
		glog.Fatalf("Error reading config key: %s; %v\n", key, err)
	}
	return number
}

// GetSlice returns the slice of values separated by sep.
func (c Config) GetSlice(key, sep string) []string {
	value := c.Get(key)
	value = strings.TrimSpace(value)
	if len(value) == 0 {
		return []string{}
	}
	v := strings.Split(value, sep)
	slice.TrimSpace(v)
	return v
}
