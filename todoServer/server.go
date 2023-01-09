package main

import "net/http"

func newMux(todoFile string) http.Handler {
	m := http.NewServeMux()
	m.HandleFunc("/", rootHandler)
	return m
}

func replyTextContent(writer http.ResponseWriter, request *http.Request, status int, content string) {
	writer.Header().Set("Content-Type", "text/plain")
	writer.WriteHeader(status)
	writer.Write([]byte(content))
}
