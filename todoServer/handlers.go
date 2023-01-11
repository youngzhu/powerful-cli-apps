package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"todo"
)

func rootHandler(writer http.ResponseWriter, request *http.Request) {
	if request.URL.Path != "/" {
		//http.NotFound(writer, request)
		replyError(writer, request, http.StatusNotFound, "")
		return
	}

	content := "There's an API here"
	replyTextContent(writer, request, http.StatusOK, content)
}

var (
	ErrNotFound    = errors.New("not found")
	ErrInvalidData = errors.New("invalid data")
)

func todoRouter(todoFile string, lck sync.Locker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		list := &todo.List{}

		lck.Lock()
		defer lck.Unlock()
		if err := list.Get(todoFile); err != nil {
			replyError(w, r, http.StatusInternalServerError, err.Error())
			return
		}

		// handle the /todo root path
		if r.URL.Path == "" {
			switch r.Method {
			case http.MethodGet:
				getAllHandler(w, r, list)
			case http.MethodPost:
				addHandler(w, r, list, todoFile)
			default:
				message := "Method not supported"
				replyError(w, r, http.StatusMethodNotAllowed, message)
			}
			return
		}

		id, err := validateID(r.URL.Path, list)
		if err != nil {
			if errors.Is(err, ErrNotFound) {
				replyError(w, r, http.StatusNotFound, err.Error())
				return
			}
			replyError(w, r, http.StatusBadRequest, err.Error())
			return
		}
		switch r.Method {
		case http.MethodGet:
			getOneHandler(w, r, list, id)
		case http.MethodDelete:
			deleteHandler(w, r, list, id, todoFile)
		case http.MethodPatch:
			patchHandler(w, r, list, id, todoFile)
		default:
			message := "Message not supported"
			replyError(w, r, http.StatusMethodNotAllowed, message)
		}
	}
}

func validateID(path string, list *todo.List) (int, error) {
	id, err := strconv.Atoi(path)
	if err != nil {
		return 0, fmt.Errorf("%w: Invalid ID: %s", ErrInvalidData, err)
	}
	if id < 1 {
		return 0, fmt.Errorf("%w: Invalid ID: Less than one", ErrInvalidData)
	}
	if id > len(*list) {
		return 0, fmt.Errorf("%w: ID %d not found", ErrNotFound, id)
	}
	return id, nil
}

func addHandler(w http.ResponseWriter, r *http.Request,
	list *todo.List, todoFile string) {
	item := struct {
		Task string `json:"task"`
	}{}

	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		message := fmt.Sprintf("Invalid JSON: %s", err)
		replyError(w, r, http.StatusBadRequest, message)
		return
	}

	list.Add(item.Task)
	if err := list.Save(todoFile); err != nil {
		replyError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	replyTextContent(w, r, http.StatusCreated, "")
}

// complete a specific item
func patchHandler(w http.ResponseWriter, r *http.Request,
	list *todo.List, id int, todoFile string) {

	q := r.URL.Query() // look for query parameters
	if _, ok := q["complete"]; !ok {
		message := "Missing query param 'complete'"
		replyError(w, r, http.StatusBadRequest, message)
		return
	}

	list.Complete(id)
	if err := list.Save(todoFile); err != nil {
		replyError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	replyTextContent(w, r, http.StatusNoContent, "")
}

func deleteHandler(w http.ResponseWriter, r *http.Request,
	list *todo.List, id int, todoFile string) {
	list.Delete(id)
	if err := list.Save(todoFile); err != nil {
		replyError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	replyTextContent(w, r, http.StatusNoContent, "")
}

func getOneHandler(w http.ResponseWriter, r *http.Request,
	list *todo.List, id int) {
	resp := &todoResponse{
		Results: (*list)[id-1 : id],
	}
	replyJSONContent(w, r, http.StatusOK, resp)
}

func getAllHandler(w http.ResponseWriter, r *http.Request, list *todo.List) {
	resp := &todoResponse{
		Results: *list,
	}
	replyJSONContent(w, r, http.StatusOK, resp)
}
