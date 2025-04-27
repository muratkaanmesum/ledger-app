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
var cronScheduler *cron.Cron

func InitScheduler() {
	start := time.Now()
	cronScheduler = cron.New()
	for _, job := range scheduledJobs {
		schedule := job.Spec
		if job.Time != "" {
			schedule = job.Time
		}
		_, err := cronScheduler.AddFunc(schedule, job.JobFunc)
		if err != nil {
			log.Printf("Failed to schedule job [%s]: %v", job.Spec, err)
		}
	}
	cronScheduler.Start()
	log.Println("All cron jobs initialized")
	log.Printf("Scheduler initialization took: %v", time.Since(start))
}

func StopScheduler() {
	if cronScheduler != nil {
		log.Println("Stopping scheduler...")
		ctx := cronScheduler.Stop()
		
		select {
		case <-ctx.Done():
			log.Println("Scheduler stopped successfully")
		case <-time.After(5 * time.Second):
			log.Println("Scheduler stop timed out after 5 seconds")
		}
	}
}

func AddSchedule(job CronJob) {
	scheduledJobs = append(scheduledJobs, job)
}
