package testup_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/devnev/testup"
)

// testOsStat emulates what a Test function for os.Stat might look like when using testup.Suite.
func testOsStat(t *testing.T) {
	// The entire callback to testup.Suite is run three times: once as a prelude, and once for every
	// callback to check. Any setup and teardown that should be executed once for the entire suite
	// should be placed before the call to testup.Suite.
	testup.Suite(t, func(t *testing.T, check testup.Register) {
		// This setup and teardown code is re-executed in each of the three runs.
		dir, err := ioutil.TempDir(".", "test")
		if err != nil {
			t.Fatalf("unable to create test directory: %s", err)
		}
		defer func() {
			err = os.RemoveAll(dir)
			if err != nil {
				t.Fatalf("unable to remove test directory: %s", err)
			}
		}()

		// In the first "prelude" run, none of the check callbacks below are executed.

		// In the second run, only this callback is executed
		check("dir is a directory", func() {
			fi, err := os.Stat(dir)
			if err != nil {
				t.Fatalf("unexpected error from stat: %s", err)
			}
			if !fi.IsDir() {
				t.Fatal("expected FileInfo.IsDir() to return true")
			}
		})
		// In the third run, only this callback is executed
		check("dir contains created file", func() {
			// We can safely mutate resource used in other check cases as each case is started with a fresh setup
			err := ioutil.WriteFile(filepath.Join(dir, "testfile"), []byte("data"), 0644)
			if err != nil {
				t.Fatalf("unexpected error from writefile: %s", err)
			}
			items, err := ioutil.ReadDir(dir)
			if err != nil {
				t.Fatalf("unexpected error from readdir: %s", err)
			}
			if len(items) != 1 || items[0].Name() != "testfile" {
				t.Fatalf("expected results to contain only testfile, got %+v", items)
			}
		})
	})
}

func Example() {
	// See testOsStat
}
