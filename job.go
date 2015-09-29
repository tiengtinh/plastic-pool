package main

import "time"

// Job holds the attributes needed to perform unit of work.
type Job struct {
	Name  string
	Delay time.Duration
}
