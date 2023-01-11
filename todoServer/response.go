package main

import (
	"encoding/json"
	"time"
	"todo"
)

type todoResponse struct {
	Results todo.List `json:"results"`
}

func (t *todoResponse) MarshalJSON() ([]byte, error) {
	resp := struct {
		Results      todo.List `json:"results"`
		Date         int64     `json:"date"`
		TotalResults int       `json:"totalResults"`
	}{
		Results:      t.Results,
		Date:         time.Now().Unix(),
		TotalResults: len(t.Results),
	}
	return json.Marshal(resp)
}
