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
"CreatedAt": "",
"Completed": ""
},
{
"Task": "Task 2",
"Done": false,
"CreatedAt": "",
"Completed": ""
}
],
"date": 0,
"totalResults": 2
}`},
	"resultOne": {
		Status: http.StatusOK,
		Body: `{
{
"Task": "Task 1",
"Done": false,
"CreatedAt": "",
"Completed": ""
}
],
"date": 0,
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
}

func mockServer(h http.HandlerFunc) (string, func()) {
	ts := httptest.NewServer(h)
	return ts.URL, func() {
		ts.Close()
	}
}
