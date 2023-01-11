package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
)

func newMux(todoFile string) http.Handler {
	m := http.NewServeMux()

	m.HandleFunc("/", rootHandler)

	mu := &sync.Mutex{}

	t := todoRouter(todoFile, mu)
	m.Handle("/todo", http.StripPrefix("/todo", t))
	m.Handle("/todo/", http.StripPrefix("/todo/", t))

	return m
}

func replyTextContent(writer http.ResponseWriter, request *http.Request, status int, content string) {
	writer.Header().Set("Content-Type", "text/plain")
	writer.WriteHeader(status)
	writer.Write([]byte(content))
}

func replyJSONContent(w http.ResponseWriter, r *http.Request,
	status int, resp *todoResponse) {
	body, err := json.Marshal(resp)
	if err != nil {
		replyError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(body)
}

func replyError(w http.ResponseWriter, r *http.Request,
	status int, message string) {
	log.Printf("%s %s: Error: %d %s", r.URL, r.Method, status, message)
	http.Error(w, http.StatusText(status), status)
}
