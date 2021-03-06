package wsl

import (
	"github.com/robfig/cron/v3"
)

type Job struct {
	Id       *cron.EntryID
	Cron     string
	MakeFunc func(*WSL) func()
}

func (this *WSL) RegisterJob(name string, job *Job) {
	this.Jobs[name] = job
}
