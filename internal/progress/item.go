package progress

import (
	"time"

	timestamp "google.golang.org/protobuf/types/known/timestamppb"
)

const (
	statusError    = "ERROR"
	statusCanceled = "CANCELED"
	statusCached   = "CACHED"
	statusDone     = "FINISHED"

	headerPrefix = "[+] "
	prefix       = " |- "
	errorDelim   = "------"
	termHeight   = 6
)

type item interface {
	lines(int) []string
	height() int
	hide() bool
}

type itemCache map[string][]item

func newItemCache() itemCache {
	return make(map[string][]item)
}

func (c itemCache) clear(key string) {
	c[key] = []item{}
}

func (c itemCache) get(key string, gen func() []item) []item {
	if items, ok := c[key]; ok && len(items) > 0 {
		return items
	}

	newItems := gen()
	c[key] = newItems

	return newItems
}

type timer struct {
	start *timestamp.Timestamp
	end   *timestamp.Timestamp
}

func newTimer(start, end *timestamp.Timestamp) timer {
	return timer{start: start, end: end}
}

func (t timer) Started() (bool, timestamp.Timestamp) {
	if t.start == nil {
		return false, timestamp.Timestamp{}
	}

	start := t.start

	return true, *start
}

func (t timer) Completed() (bool, timestamp.Timestamp) {
	if t.end == nil {
		return false, timestamp.Timestamp{}
	}

	end := t.end

	return true, *end
}

func (t timer) Running() (bool, time.Duration) {
	started, startTime := t.Started()
	if !started {
		return false, 0 * time.Second
	}

	done, endTime := t.Completed()
	if !done {
		endTime = *timestamp.Now()
	}

	runTime := endTime.AsTime().Sub(startTime.AsTime())
	if runTime < 50*time.Millisecond {
		return true, 0 * time.Second
	}

	return true, runTime
}
