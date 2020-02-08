# Go TestUp

Shared test setup/teardown for Go tests.

## Usage Example

See [testup\_test.go](testup\_test.go) for longer example.

```go
package my_test

import (
    "testing"
    "github.com/devnev/testup"
)

func TestMyType(t *testing.T) {
    testup.Suite(t, func(check testup.Register) {
        stuff := setup()
        defer func() {
          teardown(stuff)
        }()

        check("it does the thing", func(t *testing.T) {
            // assert a thing
        })
        check("it does something else", func(t *testing.T) {
            // assert something else
        })
    })
}
```
