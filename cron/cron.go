package cron

import (
	"time"
	"github.com/pengzj/swift/logger"
	"runtime"
	"fmt"
)

var (
	std *Cron
)

type Cron struct {
	entries []*Entry
	stop chan struct{}
	add chan *Entry
	location *time.Location
	running bool
}

type Job interface {
	Run()
}

type Schedule interface {
	CanTrigger(time.Time) bool
	HasFinished() bool
}

type Entry struct {
	Schedule Schedule
	Job Job
}

type FuncJob func()

func (f FuncJob) Run() {
	f()
}

func (c *Cron) AddFunc(spec string, cmd func()) error {
	return c.AddJob(spec, FuncJob(cmd))
}


func (c *Cron) AddJob(spec string, cmd Job) error {
	schedule, err := Parse(spec)
	if err != nil {
		return err
	}
	c.Schedule(schedule, cmd)
	return nil
}

func (c *Cron) Schedule(schedule Schedule, cmd Job)  {
	entry := &Entry{
		Schedule:schedule,
		Job:cmd,
	}
	if !c.running {
		c.entries = append(c.entries, entry)
	}
	c.add <- entry
}

func (c *Cron) Entries() []*Entry {
	return c.entries
}


func (c *Cron) Location() *time.Location  {
	return c.location
}

func (c *Cron) runWithRecovery(j Job)  {
	defer func() {
		if r := recover(); r != nil {
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			logger.Fatal("cron: panic running job: ", r, "\n", buf)
		}
	}()
	j.Run()
}

func (c *Cron) Run()  {
	if c.running {
		return
	}
	c.running = true
	go c.run()
}

func (c *Cron) run()  {
	ticket := time.NewTicker(time.Second * 60)
	defer ticket.Stop()
	for {
		select {
		case <-ticket.C:
			fmt.Println(time.Now(), "hello ticket")
			now := time.Now()
			for index, e := range c.entries {
				if e.Schedule.CanTrigger(now) {
					go c.runWithRecovery(e.Job)
				}

				if e.Schedule.HasFinished() {
					c.entries[index] = c.entries[len(c.entries)-1]
					c.entries = c.entries[:len(c.entries)-1]
				}
			}
		case newEntry := <-c.add:
			c.entries = append(c.entries, newEntry)
		case <- c.stop:
			c.running = false
			return
		}
	}
}

func (c *Cron) StopJob()  {
	if !c.running {
		return
	}
	c.stop <- struct {
	}{}
}

func New() *Cron  {
	return &Cron{
		entries:nil,
		add: make(chan *Entry),
		stop: make(chan struct{}),
		location:time.Now().Location(),
		running: false,
	}
}

func AddJob(spec string, cmd func()) error {
	return std.AddFunc(spec, cmd)
}

func StopJob()  {
	std.StopJob()
}

func init() {
	std = New()
	std.Run()
}





