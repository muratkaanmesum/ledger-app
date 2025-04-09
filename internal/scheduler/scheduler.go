package scheduler

import (
	"github.com/robfig/cron/v3"
	"log"
)

type CronJob struct {
	Spec    string
	JobFunc func()
}

var scheduledJobs []CronJob

func InitCronJobs() {
	c := cron.New()
	for _, job := range scheduledJobs {
		_, err := c.AddFunc(job.Spec, job.JobFunc)
		if err != nil {
			log.Printf("Failed to schedule job [%s]: %v", job.Spec, err)
		}
	}
	c.Start()
	log.Println("All cron jobs initialized")
}

func AddSchedule(CronJob []CronJob) {
	c := cron.New()
	for _, job := range CronJob {
		_, err := c.AddFunc(job.Spec, job.JobFunc)
		if err != nil {
			log.Printf("Failed to schedule job [%s]: %v", job.Spec, err)
		}
	}
}
