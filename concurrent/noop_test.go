package concurrent

import (
	"fmt"
	"testing"
	"time"
)

func Test_Noop(t *testing.T) {
	f := func() { fmt.Printf("Time %s, hello world\n", time.Now()) }
	noop := NewNoop(5, f)
	noop.Run()
}
