package cmd

import (
	"bytes"
	"errors"
	"fmt"
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
