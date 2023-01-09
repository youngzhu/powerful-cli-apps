package main

import "net/http"

func rootHandler(writer http.ResponseWriter, request *http.Request) {
	if request.URL.Path != "/" {
		http.NotFound(writer, request)
		return
	}

	content := "There's an API here"
	replyTextContent(writer, request, http.StatusOK, content)
}
