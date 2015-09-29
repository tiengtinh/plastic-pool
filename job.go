package main

import "time"

// Job holds the attributes needed to perform unit of work.
type Job struct {
	Name  string
	Delay time.Duration

	Done chan struct{}
}

//start listening for the Job to be finished processed
//func (j *Job) Start() {
//	go func() {
//		for {
//			select {
//			case <-j.Done:
//				close(j.Done)
//				return
//			}
//		}
//	}()
//}

func (j *Job) Wait() {
	<-j.Done
}
