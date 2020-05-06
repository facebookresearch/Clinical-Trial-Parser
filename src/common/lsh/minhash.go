// Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.

package lsh

import (
	"crypto/md5"
	"encoding/hex"
	"math"
	"strconv"
	"strings"

	"github.com/facebookresearch/Clinical-Trial-Parser/src/common/col/set"
	"github.com/facebookresearch/Clinical-Trial-Parser/src/common/col/tuple"
)

const (
	shingleSize = 4
	codeLength  = 10
)

type MinHash struct {
	Rows        int
	Bands       int
	ShingleSize int
	CodeLength  int
}

func New(rows, bands int) MinHash {
	return MinHash{Rows: rows, Bands: bands, ShingleSize: shingleSize, CodeLength: codeLength}
}

func (h MinHash) HashCodes(s string) set.Set {
	hashCodes := set.New()
	shingles := h.generateShingles(s)
	for i := 0; i < h.Bands; i++ {
		hashCodes.Add(strconv.Itoa(i) + "_" + h.signature(i, shingles))
	}
	return hashCodes
}

func (h MinHash) IsSimilar(s1 string, s2 string) bool {
	if len(s1) == 0 || len(s2) == 0 {
		return false
	}
	shingles1 := h.generateShingles(s1)
	shingles2 := h.generateShingles(s2)
	for i := 0; i < h.Bands; i++ {
		if h.signature(i, shingles1) == h.signature(i, shingles2) {
			return true
		}
	}
	return false
}

// Similarity computes the Jaccard similarity between strings s1 and s2.
func (h MinHash) Similarity(s1 string, s2 string) float64 {
	if len(s1) == 0 || len(s2) == 0 {
		return 0
	}
	shingles1 := h.generateShingles(s1)
	shingles2 := h.generateShingles(s2)
	return shingles1.Jaccard(shingles2)
}

// probability converts a Jaccard similarity score to probability.
func (h MinHash) probability(score float64) float64 {
	return 1.0 - math.Pow(1.0-math.Pow(score, float64(h.Rows)), float64(h.Bands))
}

func (h MinHash) generateShingles(s string) set.Set {
	shingles := set.New()
	if len(s) == 0 {
		return shingles
	}
	if len(s) < h.ShingleSize {
		s += strings.Repeat(" ", h.ShingleSize-len(s))
	}
	shingleCount := len(s) - h.ShingleSize + 1
	for i := 0; i < shingleCount; i++ {
		shingles.Add(s[i : i+h.ShingleSize])
	}
	return shingles
}

func hexdigest(s string) string {
	bytes := []byte(s)
	digest := md5.Sum(bytes)
	return hex.EncodeToString(digest[:])
}

func minhashShingle(shinglesSeedMarked []string) string {
	hashes := make(tuple.Tuples, len(shinglesSeedMarked))
	for i, shingle := range shinglesSeedMarked {
		code := hexdigest(shingle)
		hashes[i] = tuple.New(code, shingle)
	}
	hashes.Sort()
	return hashes[0][1]
}

func (h MinHash) signature(bandId int, shingles set.Set) string {
	minhashes := make([]string, h.Rows)
	seed0 := bandId * h.Rows
	for seed := 0; seed < h.Rows; seed++ {
		seedStr := strconv.Itoa(seed + seed0)
		shinglesSeedMarked := make([]string, shingles.Size())
		i := 0
		for shingle := range shingles {
			shinglesSeedMarked[i] = shingle + seedStr
			i++
		}
		minhashes[seed] = minhashShingle(shinglesSeedMarked)
	}
	joinedSignatures := strings.Join(minhashes, "_")
	return hexdigest(joinedSignatures)[:h.CodeLength]
}
