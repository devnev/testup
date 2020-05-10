# Go TestUp

Shared test setup/teardown for Go tests. Allows Suite-style tests without using
reflection to detect methods. This makes the execution easier to follow and
avoids errors caused by typos in the names of reflected methods.

[![Build Status](https://travis-ci.org/devnev/testup.svg?branch=master)](https://travis-ci.org/devnev/testup)

## Usage Example

See [testup\_test.go](testup\_test.go) for longer example.

```go
package my_test

import (
    "testing"
    "github.com/devnev/testup"
)

func TestMyType(t *testing.T) {
    // Suite setup goes here (equivalent to SetupSuite/TeardownSuite functions in suite frameworks)
    suiteStuff := setupSuite(t)
    defer suiteStuff.Teardown()

    testup.Suite(t, func(t *testing.T, check testup.Register) {
        // Test setup goes here (equivalent to SetupTest/TeardownTest functions in suite frameworks)
        stuff := setup(t)
        defer func() {
          teardown(stuff)
        }()

        // Individual test cases. The names must be static and are used as the sub-test name to `t.Run`.
        check("it does the thing", func() {
            // assert a thing
        })
        check("it does something else", func() {
            // assert something else
        })
        check("with a particular setup", func() {
          // Can have checks calls within callbacks. All setup and teardown is re-run for every check.
          check("it does another thing", func() {
            // more asserts
          })
        })
    })
}
```
