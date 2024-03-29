//go:build inmemory
// +build inmemory

package repository

import (
	"fmt"
	"pomo2/pomodoro"
	"strings"
	"sync"
	"time"
)

type memoryRepo struct {
	sync.RWMutex
	intervals []pomodoro.Interval
}

func NewMemoryRepo() *memoryRepo {
	return &memoryRepo{
		intervals: []pomodoro.Interval{},
	}
}

func (r *memoryRepo) Create(i pomodoro.Interval) (int64, error) {
	r.Lock()
	defer r.Unlock()

	i.ID = int64(len(r.intervals)) + 1
	r.intervals = append(r.intervals, i)

	return i.ID, nil
}

func (r *memoryRepo) Update(i pomodoro.Interval) error {
	r.Lock()
	defer r.Unlock()

	if i.ID == 0 {
		return fmt.Errorf("%w: %d", pomodoro.ErrInvalidID, i.ID)
	}
	r.intervals[i.ID-1] = i
	return nil
}

func (r *memoryRepo) ByID(id int64) (pomodoro.Interval, error) {
	r.RLock()
	defer r.RUnlock()

	i := pomodoro.Interval{}
	if id == 0 {
		return i, fmt.Errorf("%w: %d", pomodoro.ErrInvalidID, id)
	}

	i = r.intervals[id-1]
	return i, nil
}

func (r *memoryRepo) Last() (pomodoro.Interval, error) {
	r.RLock()
	defer r.RUnlock()

	i := pomodoro.Interval{}
	size := len(r.intervals)
	if size == 0 {
		return i, pomodoro.ErrNoIntervals
	}
	return r.intervals[size-1], nil
}

func (r *memoryRepo) Breaks(n int) ([]pomodoro.Interval, error) {
	r.RLock()
	defer r.RUnlock()

	var data []pomodoro.Interval
	for k := len(r.intervals) - 1; k >= 0; k-- {
		if r.intervals[k].Category == pomodoro.CategoryPomodoro {
			continue
		}
		data = append(data, r.intervals[k])
		if len(data) == n {
			return data, nil
		}
	}
	return data, nil
}

// return a daily summary
func (r *memoryRepo) CategorySummary(day time.Time, filter string) (time.Duration, error) {
	r.RLock()
	defer r.RUnlock()

	var d time.Duration
	filter = strings.Trim(filter, "%")
	for _, i := range r.intervals {
		if i.StartTime.Year() == day.Year() &&
			i.StartTime.YearDay() == day.YearDay() {
			if strings.Contains(i.Category, filter) {
				d += i.ActualDuration
			}
		}
	}
	return d, nil
}
