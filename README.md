# defuse

A Go utility, usually for testing, that helps ensure that routines complete on
time, that context cancelations are properly handled, or that the process does
not hang forever due to a deadlock.

The concept is fairly simple. A timer is created that when reached causes a
deleterious effect, usually a panic or a fatal test failure. The routine under
test must complete its work and then defuse this bomb before it goes off.

```go
package foo

import (
	"context"
	"testing"
	"time"

	"github.com/colvin/defuse"
)

func TestFoo(t *testing.T) {
	// Run the Foo routine and then cancel it. The bomb ensures that if Foo is
	// not properly canceled the test will not hang forever.
	defuse.DoOrElse(
		func() {
			ctx, cancel := context.WithCancel(context.Background())
			go Foo(ctx)
			// Simulate some runtime
			time.Sleep(100 * time.Millisecond)
			cancel()
		},
		200 * time.Millisecond,
		defuse.FailTest(t, "Foo was not canceled"),
	)
}
```
