package pomodoro

import (
	"context"
	"errors"
	"fmt"
	"time"
)

const (
	CategoryPomodoro   = "Pomodoro"
	CategoryShortBreak = "ShortBreak"
	CategoryLongBreak  = "LongBreak"
)

const (
	StateNotStarted = iota
	StateRunning
	StatePaused
	StateDone
	StateCancelled
)

type Interval struct {
	ID              int64
	StartTime       time.Time
	PlannedDuration time.Duration
	ActualDuration  time.Duration
	Category        string
	State           int
}

type Repository interface {
	// Create to create/save a new Interval
	Create(i Interval) (int64, error)
	// Update to update details about an interval
	Update(i Interval) error
	// ByID to retrieve a specific Interval by its ID
	ByID(id int64) (Interval, error)
	// Last to find the last Interval
	Last() (Interval, error)
	// Breaks to retrieve a given number of Interval items
	// that matches CategoryLongBreak or CategoryShortBreak
	Breaks(n int) ([]Interval, error)
	CategorySummary(day time.Time, filter string) (time.Duration, error)
}

var (
	ErrNoIntervals        = errors.New("no intervals")
	ErrIntervalNotRunning = errors.New("interval not running")
	ErrIntervalCompleted  = errors.New("interval is completed or cancelled")
	ErrInvalidState       = errors.New("invalid state")
	ErrInvalidID          = errors.New("invalid ID")
)

type IntervalConfig struct {
	repo               Repository
	PomodoroDuration   time.Duration
	ShortBreakDuration time.Duration
	LongBreakDuration  time.Duration
}

func NewConfig(repo Repository, pomodoro, shortBreak, longBreak time.Duration) *IntervalConfig {
	c := &IntervalConfig{
		repo:               repo,
		PomodoroDuration:   25 * time.Minute,
		ShortBreakDuration: 5 * time.Minute,
		LongBreakDuration:  15 * time.Minute,
	}
	if pomodoro > 0 {
		c.PomodoroDuration = pomodoro
	}
	if shortBreak > 0 {
		c.ShortBreakDuration = shortBreak
	}
	if longBreak > 0 {
		c.LongBreakDuration = longBreak
	}
	return c
}

// takes a reference to the repository as input
// and returns the next interval category as a string or an error
func nextCategory(r Repository) (string, error) {
	last, err := r.Last()
	if err != nil {
		if err == ErrNoIntervals {
			return CategoryPomodoro, nil
		}
		return "", err
	}
	if last.Category == CategoryLongBreak || last.Category == CategoryShortBreak {
		return CategoryPomodoro, nil
	}

	lastBreaks, err := r.Breaks(3)
	if err != nil {
		return "", err
	}
	if len(lastBreaks) < 3 {
		return CategoryShortBreak, nil
	}
	for _, i := range lastBreaks {
		if i.Category == CategoryLongBreak {
			return CategoryShortBreak, nil
		}
	}
	return CategoryLongBreak, nil
}

type Callback func(interval Interval)

// control the interval timer.
// ctx: indicates a cancellation
// id: the Interval to control
// config:
// start: Callback function that execute at the start
// periodic: Callback function that execute periodically
// end: Callback function that execute at the end
func tick(ctx context.Context, id int64, config *IntervalConfig,
	start, periodic, end Callback) error {

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	i, err := config.repo.ByID(id)
	if err != nil {
		return err
	}

	expire := time.After(i.PlannedDuration - i.ActualDuration)
	start(i)
	for {
		select {
		case <-ticker.C:
			i, err := config.repo.ByID(id)
			if err != nil {
				return err
			}
			if i.State == StatePaused {
				return nil
			}
			i.ActualDuration += time.Second
			if err := config.repo.Update(i); err != nil {
				return err
			}
			periodic(i)
		case <-expire:
			i, err := config.repo.ByID(id)
			if err != nil {
				return err
			}
			i.State = StateDone
			end(i)
			return config.repo.Update(i)
		case <-ctx.Done():
			i, err := config.repo.ByID(id)
			if err != nil {
				return err
			}
			i.State = StateCancelled
			return config.repo.Update(i)
		}
	}
}

func newInterval(config *IntervalConfig) (Interval, error) {
	i := Interval{}
	category, err := nextCategory(config.repo)
	if err != nil {
		return i, err
	}
	i.Category = category

	switch category {
	case CategoryPomodoro:
		i.PlannedDuration = config.PomodoroDuration
	case CategoryShortBreak:
		i.PlannedDuration = config.ShortBreakDuration
	case CategoryLongBreak:
		i.PlannedDuration = config.LongBreakDuration
	}

	if i.ID, err = config.repo.Create(i); err != nil {
		return i, err
	}

	return i, nil
}

func GetInterval(config *IntervalConfig) (Interval, error) {
	i := Interval{}
	var err error

	i, err = config.repo.Last()
	if err != nil && err != ErrNoIntervals {
		return i, err
	}
	if err == nil && i.State != StateCancelled && i.State != StateDone {
		return i, nil
	}
	return newInterval(config)
}

func (i Interval) Start(ctx context.Context, config *IntervalConfig,
	start, periodic, end Callback) error {

	switch i.State {
	case StateRunning:
		return nil
	case StateNotStarted:
		i.StartTime = time.Now()
		fallthrough
	case StatePaused:
		i.State = StateRunning
		if err := config.repo.Update(i); err != nil {
			return err
		}
		return tick(ctx, i.ID, config, start, periodic, end)
	case StateCancelled, StateDone:
		return fmt.Errorf("%w: cannot start", ErrIntervalCompleted)
	default:
		return fmt.Errorf("%w: %d", ErrInvalidState, i.State)
	}
}

func (i Interval) Pause(config *IntervalConfig) error {
	if i.State != StateRunning {
		return ErrIntervalNotRunning
	}
	i.State = StatePaused
	return config.repo.Update(i)
}
