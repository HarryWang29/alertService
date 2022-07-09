package main

import (
	"alertService/config"
	"alertService/job"
)

func main() {
	myApp := config.New()
	go func() {
		job.Start(myApp.Jobs)
	}()
	select {}
}
