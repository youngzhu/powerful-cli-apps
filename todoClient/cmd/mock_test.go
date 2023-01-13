package cmd

import (
	"net/http"
	"net/http/httptest"
)

var testResp = map[string]struct {
	Status int
	Body   string
}{
	"resultMany": {
		Status: http.StatusOK,
		Body: `{
"results": [
{
"Task": "Task 1",
"Done": false,
"CreateAt": "2023-01-12T17:00:51.3695194+08:00",
"Completed": "0001-01-01T00:00:00Z"
},
{
"Task": "Task 2",
"Done": false,
"CreateAt": "2023-01-12T17:00:51.3695194+08:00",
"CompletedAt": "0001-01-01T00:00:00Z"
}
],
"date": 0,
"totalResults": 2
}`},
	"resultOne": {
		Status: http.StatusOK,
		Body: `{
"results": [
{
"Task": "Task 1",
"Done": false,
"CreateAt": "2023-01-12T17:00:51.3695194+08:00",
"CompletedAt": "0001-01-01T00:00:00Z"
}
],
"date": 1572265440,
"totalResults": 1
}`,
	},
	"noResults": {
		Status: http.StatusOK,
		Body: `{
"results": [],
"date": 0,
"totalResults": 0
}`,
	},
	"root": {
		Status: http.StatusOK,
		Body:   "There's an API here",
	},
	"notFound": {
		Status: http.StatusNotFound,
		Body:   "404 - not found",
	},
	"created": {
		Status: http.StatusCreated,
		Body:   "",
	},
	"noContent": {
		Status: http.StatusNoContent,
		Body:   "",
	},
}

func mockServer(h http.HandlerFunc) (string, func()) {
	ts := httptest.NewServer(h)
	return ts.URL, func() {
		ts.Close()
	}
}
