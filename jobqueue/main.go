package main

import (
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type Job interface {
	Do()
	Key() string
}

type reqJob struct {
	url  *url.URL
	resp *http.Response
	err  error
}

func (r *reqJob) Do() {
	r.resp, r.err = http.Get(r.url.String())
}

func (r *reqJob) Key() string {
	return r.url.Hostname()
}

type JobQueue struct {
	sync.Mutex

	jobs chan Job
	done chan Job

	delay       time.Duration
	keys        map[string]time.Time
	concurrency int
}

func NewJobQueue(concurrency int, delay time.Duration) *JobQueue {
	q := &JobQueue{
		jobs:  make(chan Job),
		done:  make(chan Job),
		keys:  make(map[string]time.Time),
		delay: delay,
	}

	// spin up the workers
	var wg sync.WaitGroup
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			for j := range q.jobs {
				q.Block(j.Key())
				j.Do()
				q.done <- j
			}
			wg.Done()
		}()
	}

	// close the done channel when the workers are finished
	go func() {
		wg.Wait()
		close(q.done)
	}()

	return q
}

func (q *JobQueue) Finish() {
	close(q.jobs)
}

func (q *JobQueue) Enqueue(j Job) {
	q.jobs <- j
}

func (q *JobQueue) FinishedJobs() chan Job {
	return q.done
}

func (q *JobQueue) Block(key string) {
	now := time.Now()

	q.Lock()

	// if there's nothing in the map we can
	// return straight away, tracking when
	// the job ran
	if _, ok := q.keys[key]; !ok {
		q.keys[key] = now
		q.Unlock()
		return
	}

	// if time is up we can return straight away
	t := q.keys[key]
	deadline := t.Add(q.delay)
	if now.After(deadline) {
		q.keys[key] = now
		q.Unlock()
		return
	}

	remaining := deadline.Sub(now)

	// Set the time of the operation
	q.keys[key] = now.Add(remaining)
	q.Unlock()

	// Block for the remaining time
	<-time.After(remaining)
}

func main() {
	q := NewJobQueue(20, time.Second)

	u, _ := url.Parse("https://example.com")
	r := &reqJob{url: u}

	go func() {
		q.Enqueue(r)
		q.Enqueue(r)
		q.Enqueue(r)
		q.Finish()
	}()

	for j := range q.FinishedJobs() {
		r, ok := j.(*reqJob)
		if !ok {
			continue
		}
		fmt.Printf("status: %s err: %v\n", r.resp.Status, r.err)
	}

}
