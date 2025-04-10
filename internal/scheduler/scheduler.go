package scheduler

import (
	"github.com/robfig/cron/v3"
	"log"
	"time"
)

type CronJob struct {
	Spec    string
	JobFunc func()
	Time    string
}

var scheduledJobs []CronJob

func InitScheduler() {
	start := time.Now()
	c := cron.New()
	for _, job := range scheduledJobs {
		schedule := job.Spec
		if job.Time != "" {
			schedule = job.Time
		}
		_, err := c.AddFunc(schedule, job.JobFunc)
		if err != nil {
			log.Printf("Failed to schedule job [%s]: %v", job.Spec, err)
		}
	}
	c.Start()
	log.Println("All cron jobs initialized")
	log.Printf("Scheduler initialization took: %v", time.Since(start))
}

func AddSchedule(job CronJob) {
	scheduledJobs = append(scheduledJobs, job)
}
