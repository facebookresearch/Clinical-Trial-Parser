// Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.

package fio

import (
	"bufio"
	"io/ioutil"
	"os"
	"strings"

	"github.com/facebookresearch/Clinical-Trial-Parser/src/common/col/set"
	"github.com/facebookresearch/Clinical-Trial-Parser/src/common/col/tuple"

	"github.com/golang/glog"
)

func Writer(fname string) *os.File {
	writer, err := os.OpenFile(fname, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
	if err != nil {
		glog.Fatalln(err)
	}
	return writer
}

func LoadList(fname string, delim string) []string {
	list := []string{}
	file, _ := os.Open(fname)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	lineCnt := 0
	for scanner.Scan() {
		lineCnt++
		line := scanner.Text()
		values := strings.Split(line, delim)
		value := strings.TrimSpace(values[0])
		if len(value) > 0 {
			list = append(list, value)
		}
	}
	glog.Infof("%s: %d, entities: %d\n", fname, lineCnt, len(list))
	return list
}

func LoadSet(fname string, delim string) set.Set {
	set := set.New()
	file, _ := os.Open(fname)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	lineCnt := 0
	for scanner.Scan() {
		lineCnt++
		line := scanner.Text()
		values := strings.Split(line, delim)
		value := strings.TrimSpace(values[0])
		if len(value) > 0 {
			set.Add(value)
		}
	}
	glog.Infof("%s: %d, entities: %d\n", fname, lineCnt, set.Size())
	return set
}

func LoadMap(fname string, delim string) map[string]string {
	mapping := make(map[string]string)
	file, _ := os.Open(fname)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	lineCnt := 0
	for scanner.Scan() {
		lineCnt++
		line := scanner.Text()
		values := strings.Split(line, delim)
		mapping[strings.TrimSpace(values[0])] = strings.TrimSpace(values[1])
	}
	glog.Infof("%s: %d\n", fname, lineCnt)
	return mapping
}

func LoadTuples(fname string, delim string) tuple.Tuples {
	var list tuple.Tuples
	file, _ := os.Open(fname)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	lineCnt := 0
	for scanner.Scan() {
		lineCnt++
		line := scanner.Text()
		values := strings.Split(line, delim)
		tuple := tuple.New(strings.TrimSpace(values[0]), strings.TrimSpace(values[1]))
		list = append(list, tuple)
	}
	glog.Infof("%s: %d\n", fname, lineCnt)
	return list
}

func Files(str string) []string {
	i := strings.LastIndex(str, "/") + 1
	path := strings.Trim(str[:i], " ")
	fnames := strings.Split(str[i:], ";")
	files := make([]string, len(fnames))
	for k, fname := range fnames {
		files[k] = path + strings.TrimSpace(fname)
	}
	return files
}

func ReadFnames(path string) []string {
	if strings.Contains(path, ";") {
		return Files(path)
	}

	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		glog.Fatalln(err)
	}
	stat, err := file.Stat()
	if err != nil {
		glog.Fatalln(err)
	}

	if stat.IsDir() {
		var fnames []string
		content, _ := ioutil.ReadDir(path)
		for _, f := range content {
			fname := f.Name()

			if (f.Mode().IsRegular() || (f.Mode()&os.ModeSymlink != 0)) && len(fname) > 0 && fname[0] != '.' {
				fpath := path + "/" + f.Name()
				fnames = append(fnames, fpath)
			}
		}
		return fnames
	} else {
		return []string{path}
	}
}
