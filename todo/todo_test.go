package todo_test

import (
	"io/ioutil"
	"os"
	"testing"
	"todo"
)

func TestList_Add(t *testing.T) {
	ls := todo.List{}

	task := "New Task"
	ls.Add(task)

	if ls[0].Task != task {
		t.Errorf("Expected %q, got %q instead.", task, ls[0].Task)
	}
}

func TestList_Complete(t *testing.T) {
	ls := todo.List{}

	task := "New Task"
	ls.Add(task)
	if ls[0].Task != task {
		t.Errorf("Expected %q, got %q instead.", task, ls[0].Task)
	}

	if ls[0].Done {
		t.Error("New task should not be completed")
	}

	ls.Complete(1)
	if !ls[0].Done {
		t.Error("New task should be completed")
	}
}

func TestList_Delete(t *testing.T) {
	ls := todo.List{}

	tasks := []string{
		"New Task 1",
		"New Task 2",
		"New Task 3",
	}
	for _, v := range tasks {
		ls.Add(v)
	}

	if ls[0].Task != tasks[0] {
		t.Errorf("Expected %q, got %q instead.", tasks[0], ls[0].Task)
	}

	ls.Delete(2)
	if len(ls) != 2 {
		t.Errorf("Expected list length 2, got %d instead.", len(ls))
	}
	if ls[1].Task != tasks[2] {
		t.Errorf("Expected %q, got %q instead.", tasks[2], ls[1].Task)
	}
}

func TestList_SaveGet(t *testing.T) {
	ls1 := todo.List{}
	ls2 := todo.List{}

	task := "New Task"
	ls1.Add(task)
	if ls1[0].Task != task {
		t.Errorf("Expected %q, got %q instead.", task, ls1[0].Task)
	}

	tf, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatalf("Error creating temp file: %s", err)
	}
	defer os.Remove(tf.Name())

	if err := ls1.Save(tf.Name()); err != nil {
		t.Fatalf("Error saving list to file: %s", err)
	}
	if err := ls2.Get(tf.Name()); err != nil {
		t.Fatalf("Error getting list from file: %s", err)
	}
	if len(ls2) != len(ls1) {
		t.Errorf("Expected list length %q, got %d instead.", len(ls1), len(ls2))
	}
	if ls1[0].Task != ls2[0].Task {
		t.Errorf("Task %q should match %q task.", ls1[0].Task, ls2[0].Task)
	}
}
