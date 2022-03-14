// Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.

package resource

import (
	"os"
	"strings"
	"testing"
)

const (
	dataPathEnvVar     = "DATA_PATH"
	resourceDir        = "src/common/resource"
	resourcePathEnvVar = "RESOURCE_PATH"
)

func TestGetResourcePath(t *testing.T) {
	tests := []struct {
		name         string
		resourcePath string
		want         string
	}{
		{
			name:         "setting RESOURCE_PATH returns its value",
			resourcePath: "some-value",
			want:         "some-value",
		},
		{
			name: "no RESOURCE_PATH attempts to use LocalResourcePath",
			want: "src/resources",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.resourcePath != "" {
				os.Setenv(resourcePathEnvVar, tt.resourcePath)
			} else {
				os.Unsetenv(resourcePathEnvVar)
			}

			wd, err := os.Getwd()
			if err != nil {
				t.Fatalf("os.Getwd: %v", err)
			}
			if strings.HasSuffix(wd, resourceDir) {
				if err := os.Chdir("../../.."); err != nil {
					t.Fatalf("os.Chdir: %v", err)
				}
			}

			if got := GetResourcePath(); got != tt.want {
				t.Errorf("GetResourcePath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetDataPath(t *testing.T) {
	tests := []struct {
		name     string
		dataPath string
		want     string
	}{
		{
			name:     "setting DATA_PATH returns its value",
			dataPath: "some-value",
			want:     "some-value",
		},
		{
			name: "no DATA_PATH attempts to use LocalDataPath",
			want: "data",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.dataPath != "" {
				os.Setenv(dataPathEnvVar, tt.dataPath)
			} else {
				os.Unsetenv(dataPathEnvVar)
			}

			wd, err := os.Getwd()
			if err != nil {
				t.Fatalf("os.Getwd: %v", err)
			}
			if strings.HasSuffix(wd, resourceDir) {
				if err := os.Chdir("../../.."); err != nil {
					t.Fatalf("os.Chdir: %v", err)
				}
			}

			if got := GetDataPath(); got != tt.want {
				t.Errorf("GetDataPath() = %v, want %v", got, tt.want)
			}
		})
	}
}
