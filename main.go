package main

import "dozenplans/job"

func main() {
	job.StartTimingJobs()
	startHttpServer()
}
