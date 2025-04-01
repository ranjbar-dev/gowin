package data

import "github.com/ranjbar-dev/gowin/types"

// add job to queue for client
func (d *Data) AddJob(clientID string, job types.Job) {

	d.jobsMutex.Lock()
	defer d.jobsMutex.Unlock()

	d.jobs[clientID] = append(d.jobs[clientID], job)

}

// get jobs for client and clear them
func (d *Data) PullJobs(clientID string) []types.Job {

	d.jobsMutex.Lock()
	defer d.jobsMutex.Unlock()

	jobs := d.jobs[clientID]

	d.jobs[clientID] = []types.Job{}

	return jobs
}
