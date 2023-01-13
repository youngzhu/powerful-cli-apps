package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestListAction(t *testing.T) {
	testCases := []struct {
		name   string
		expErr error
		expOut string
		resp   struct {
			Status int
			Body   string
		}
		closeServer bool
	}{
		{
			name:   "Results",
			expErr: nil,
			expOut: "-  1  Task 1\n-  2  Task 2\n",
			resp:   testResp["resultMany"],
		},
		{
			name:   "NoResults",
			expErr: ErrNotFound,
			resp:   testResp["noResults"],
		},
		{
			name:        "InvalidURL",
			expErr:      ErrConnection,
			resp:        testResp["noResults"],
			closeServer: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url, cleanup := mockServer(
				func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(tc.resp.Status)
					fmt.Fprintln(w, tc.resp.Body)
				})
			defer cleanup()

			if tc.closeServer {
				cleanup()
			}

			var out bytes.Buffer
			err := listAtion(&out, url)

			if tc.expErr != nil {
				if err == nil {
					t.Fatalf("Expected error %q, got no error.", tc.expErr)
				}
				if !errors.Is(err, tc.expErr) {
					t.Errorf("Expected error %q, got %q.", tc.expErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("Expected no error, got %q.", err)
			}
			if tc.expOut != out.String() {
				t.Errorf("Expected output %q, got %q", tc.expOut, out.String())
			}
		})
	}
}

func TestViewAction(t *testing.T) {
	testCases := []struct {
		name   string
		expErr error
		expOut string
		resp   struct {
			Status int
			Body   string
		}
		id string
	}{
		{
			name:   "ResultOne",
			expErr: nil,
			expOut: "Task:         Task 1\nCreated at:   Jan/12 @17:00\nCompleted:    No\n",
			resp:   testResp["resultOne"],
			id:     "1",
		},
		{
			name:   "NotFound",
			expErr: ErrNotFound,
			resp:   testResp["notFound"],
			id:     "1",
		},
		{
			name:   "InvalidID",
			expErr: ErrNotNumber,
			resp:   testResp["noResults"],
			id:     "a",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url, cleanup := mockServer(
				func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(tc.resp.Status)
					fmt.Fprintln(w, tc.resp.Body)
				})
			defer cleanup()

			var out bytes.Buffer
			err := viewAction(&out, url, tc.id)

			if tc.expErr != nil {
				if err == nil {
					t.Fatalf("Expected error %q, got no error.", tc.expErr)
				}
				if !errors.Is(err, tc.expErr) {
					t.Errorf("Expected error %q, got %q.", tc.expErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("Expected no error, got %q.", err)
			}
			if tc.expOut != out.String() {
				t.Errorf("Expected output %q, got %q", tc.expOut, out.String())
			}
		})
	}
}

func TestAddAction(t *testing.T) {
	expURLPath := "/todo"
	expMethod := http.MethodPost
	expBody := "{\"task\":\"Task 1\"}\n"
	expContentType := "application/json"
	expOut := "Added task \"Task 1\" to the list.\n"
	args := []string{"Task", "1"}

	url, cleanup := mockServer(
		func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != expURLPath {
				t.Errorf("Expected path %q, got %q", expURLPath, r.URL.Path)
			}
			if r.Method != expMethod {
				t.Errorf("Expected method %q, got %q", expMethod, r.Method)
			}

			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				t.Fatal(err)
			}
			r.Body.Close()

			if string(body) != expBody {
				t.Errorf("Expected body %q, got %q", expBody, string(body))
			}

			contentType := r.Header.Get("Content-Type")
			if contentType != expContentType {
				t.Errorf("Expected Content-Type %q, got %q",
					expContentType, contentType)
			}

			w.WriteHeader(testResp["created"].Status)
			fmt.Fprintln(w, testResp["created"].Body)
		})
	defer cleanup()

	var out bytes.Buffer
	if err := addAction(&out, url, args); err != nil {
		t.Fatalf("Expected no err, got %q.", err)
	}
	if expOut != out.String() {
		t.Errorf("Expected output %q, got %q", expOut, out.String())
	}
}

func TestCompleteAction(t *testing.T) {
	expURLPath := "/todo/1"
	expMethod := http.MethodPatch
	expQuery := "complete"
	expOut := "Item number 1 marked as completed.\n"
	arg := "1"

	url, cleanup := mockServer(
		func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != expURLPath {
				t.Errorf("Expected path %q, got %q", expURLPath, r.URL.Path)
			}
			if r.Method != expMethod {
				t.Errorf("Expected method %q, got %q", expMethod, r.Method)
			}
			if _, ok := r.URL.Query()[expQuery]; !ok {
				t.Errorf("Expected query %q not found in URL", expQuery)
			}

			w.WriteHeader(testResp["noContent"].Status)
			fmt.Fprintln(w, testResp["noContent"].Body)
		})
	defer cleanup()

	var out bytes.Buffer
	if err := completeAction(&out, url, arg); err != nil {
		t.Fatalf("Expected no error, got %q.", err)
	}
	if out.String() != expOut {
		t.Errorf("Expected output %q, got %q", expOut, out.String())
	}
}

func TestDelAction(t *testing.T) {
	expURLPath := "/todo/1"
	expMethod := http.MethodDelete
	expOut := "Item number 1 deleted.\n"
	arg := "1"

	url, cleanup := mockServer(
		func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != expURLPath {
				t.Errorf("Expected path %q, got %q", expURLPath, r.URL.Path)
			}
			if r.Method != expMethod {
				t.Errorf("Expected method %q, got %q", expMethod, r.Method)
			}

			w.WriteHeader(testResp["noContent"].Status)
			fmt.Fprintln(w, testResp["noContent"].Body)
		})
	defer cleanup()

	var out bytes.Buffer
	if err := delAction(&out, url, arg); err != nil {
		t.Fatalf("Expected no error, got %q.", err)
	}
	if out.String() != expOut {
		t.Errorf("Expected output %q, got %q", expOut, out.String())
	}
}
