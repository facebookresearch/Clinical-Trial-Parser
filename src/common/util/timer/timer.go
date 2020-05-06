// Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.

package timer

import (
	"fmt"
	"time"
)

// Timer measures the elapsed time.
type Timer struct {
	start time.Time
}

func New() Timer {
	start := time.Now()
	return Timer{start}
}

func (timer Timer) Elapsed() string {
	d := time.Since(timer.start)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	s := float64(d) / float64(time.Second)
	return fmt.Sprintf("%02d:%02d:%06.3f", h, m, s)
}

func (timer Timer) PrintElapsed() {
	fmt.Println("Elapsed time: " + timer.Elapsed())
}
