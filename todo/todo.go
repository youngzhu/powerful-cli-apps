package todo

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"
)

type item struct {
	Task        string
	Done        bool
	CreateAt    time.Time
	CompletedAt time.Time
}

type List []item

// Add creates a new to-do item and appends it to the list
func (list *List) Add(task string) {
	t := item{
		Task:        task,
		Done:        false,
		CreateAt:    time.Now(),
		CompletedAt: time.Time{},
	}
	*list = append(*list, t)
}

// Complete marks a to-do item as completed by
// setting Done=true and CompletedAt to the current time
func (list *List) Complete(i int) error {
	ls := *list
	if i <= 0 || i > len(ls) {
		return fmt.Errorf("item %d does not exist", i)
	}

	ls[i-1].Done = true
	ls[i-1].CompletedAt = time.Now()

	return nil
}

// Delete remove a to-do item from the list
func (list *List) Delete(i int) error {
	ls := *list
	if i <= 0 || i > len(ls) {
		return fmt.Errorf("item %d does not exist", i)
	}

	*list = append(ls[:i-1], ls[i:]...)

	return nil
}

// Save encodes the List as JSON and saves it using
// the provided file name
func (list *List) Save(filename string) error {
	js, err := json.Marshal(list)
	if err != nil {
		return err
	}

	return os.WriteFile(filename, js, 0644)
}

// Get open the proided file name, decodes
// the JSON data and parses it into a list
func (list *List) Get(filename string) error {
	file, err := os.ReadFile(filename)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}
	if len(file) == 0 {
		return nil
	}
	return json.Unmarshal(file, list)
}

func (list *List) String() string {
	formatted := ""

	for i, v := range *list {
		prefix := "  "
		if v.Done {
			prefix = "X "
		}
		formatted += fmt.Sprintf("%s%d: %s\n", prefix, i+1, v.Task)
	}

	return formatted
}
