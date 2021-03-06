package storage

import (
	"fmt"
	"sync"

	"github.com/SeerUK/reverb/model"
)

// MemoryDriver is an in-memory data store for requests.
type MemoryDriver struct {
	sync.RWMutex
	requests []model.Request
}

// FindAll takes a destination to write results into. If there are no results, the destination will
// not be written to. If there is an error, it will be returned.
func (d *MemoryDriver) FindAll(dst *[]model.Request) error {
	d.RLock()
	defer d.RUnlock()

	*dst = d.requests

	return nil
}

// Find takes an ID and a destination to write the result into if a request is found, otherwise an
// error will be returned.
func (d *MemoryDriver) Find(id int, dst *model.Request) error {
	d.RLock()
	defer d.RUnlock()

	for _, req := range d.requests {
		if int(req.ID) == id {
			*dst = req

			return nil
		}
	}

	return fmt.Errorf("No request could be found with the ID '%d'. Maybe it's been popped out?", id)
}

// Persist takes a request and persists it in memory.
func (d *MemoryDriver) Persist(src *model.Request) error {
	d.Lock()
	defer d.Unlock()

	if len(d.requests) > 0 {
		prev := d.requests[len(d.requests)-1]
		src.ID = prev.ID + 1
	} else {
		src.ID = 1
	}

	d.requests = append(d.requests, *src)

	count := len(d.requests)
	if count > 100 {
		// Currently we're limiting the number of requests stored in memory to 100.
		d.requests = d.requests[count-100 : count]
	}

	return nil
}

// Flush completely clears the in-memory storage.
func (d *MemoryDriver) Flush() error {
	d.requests = d.requests[:0]

	return nil
}
