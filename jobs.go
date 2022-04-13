package wsl

import (
	"github.com/robfig/cron/v3"
)

type Job struct {
	ID       *cron.EntryID
	Cron     string
	MakeFunc func(*App) func()
}

func (this *App) RegisterJob(name string, job *Job) {
	this.Jobs[name] = job
}
