package testup_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/devnev/testup"
)

// A standard test function (underscore added to aid `godoc`)
func _TestMyType(t *testing.T) {
	testup.Suite(t, func(t *testing.T, check testup.Register) {
		// This setup and teardown code will be executed once as a prelude, and once for every callback to check
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

		check("dir is a directory", func() {
			fi, err := os.Stat(dir)
			if err != nil {
				t.Fatalf("unexpected error from stat: %s", err)
			}
			if !fi.IsDir() {
				t.Fatal("expected FileInfo.IsDir() to return true")
			}
		})
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
			if len(items) != 1 || items[1].Name() != "testfile" {
				t.Fatalf("expected results to contain only testfile, got %+v", items)
			}
		})
	})
}

func ExampleSuite() {
	// see above
}
