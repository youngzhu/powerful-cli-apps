package scan_test

import (
	"errors"
	"io/ioutil"
	"os"
	"portscan/scan"
	"testing"
)

func TestHostsList_Add(t *testing.T) {
	testCases := []struct {
		name      string
		host      string
		expectLen int
		expectErr error
	}{
		{"AddNew", "host2", 2, nil},
		{"AddExisted", "host1", 1, scan.ErrExists},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			hl := &scan.HostsList{}
			// initialize list
			if err := hl.Add("host1"); err != nil {
				t.Fatal(err)
			}
			err := hl.Add(tc.host)

			if tc.expectErr != nil {
				if err == nil {
					t.Fatalf("Expected error, got nil instead\n")
				}
				if !errors.Is(err, tc.expectErr) {
					t.Errorf("Expected error %q, got %q instead\n",
						tc.expectErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("Expected no error, got %q instead\n", err)
			}

			if len(hl.Hosts) != tc.expectLen {
				t.Errorf("Expected list length %d, got %d instead\n",
					tc.expectLen, len(hl.Hosts))
			}
			if hl.Hosts[1] != tc.host {
				t.Errorf("Expected host name %q as index 1, got %q instead\n",
					tc.host, hl.Hosts[1])
			}
		})
	}
}

func TestHostsList_Remove(t *testing.T) {
	testCases := []struct {
		name      string
		host      string
		expectLen int
		expectErr error
	}{
		{"RemoveExisted", "host1", 1, nil},
		{"RemoveNotExisted", "host3", 1, scan.ErrNotExists},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			hl := &scan.HostsList{}
			// initialize list
			for _, h := range []string{"host1", "host2"} {
				if err := hl.Add(h); err != nil {
					t.Fatal(err)
				}
			}

			err := hl.Remove(tc.host)

			if tc.expectErr != nil {
				if err == nil {
					t.Fatalf("Expected error, got nil instead\n")
				}
				if !errors.Is(err, tc.expectErr) {
					t.Errorf("Expected error %q, got %q instead\n",
						tc.expectErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("Expected no error, got %q instead\n", err)
			}

			if len(hl.Hosts) != tc.expectLen {
				t.Errorf("Expected list length %d, got %d instead\n",
					tc.expectLen, len(hl.Hosts))
			}
			if hl.Hosts[0] == tc.host {
				t.Errorf("Host name %q should not be in the list\n",
					tc.host)
			}
		})
	}
}

func TestHostsList_SaveLoad(t *testing.T) {
	hl1 := scan.HostsList{}
	hl2 := scan.HostsList{}

	hostname := "host1"
	hl1.Add(hostname)

	tf, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatalf("Error creating temp file: %s", err)
	}
	defer os.Remove(tf.Name())

	if err := hl1.Save(tf.Name()); err != nil {
		t.Fatalf("Error saving list to file: %s", err)
	}
	if err := hl2.Load(tf.Name()); err != nil {
		t.Fatalf("Error getting list from file: %s", err)
	}
	if hl1.Hosts[0] != hl2.Hosts[0] {
		t.Errorf("host %q should match %q host.", hl1.Hosts[0], hl2.Hosts[0])
	}
}

func TestHostsList_Load_noFile(t *testing.T) {
	filename := "notExists"

	/*
		tf, err := ioutil.TempFile("", "")
		if err != nil {
			t.Fatalf("Error creating temp file: %s", err)
		}
		//The process cannot access the file be
		//cause it is being used by another process.
		if err := os.Remove(tf.Name()); err != nil {
			t.Fatalf("Error deleting temp file: %s", err)
		}
		filename=tf.Name()
	*/

	hl := &scan.HostsList{}
	if err := hl.Load(filename); err != nil {
		t.Errorf("Expected no error, got %q instead\n", err)
	}
}