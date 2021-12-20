package collection

import "context"

type fibonacciGenerator struct{}

// Fibonacci is an Enumerable which will dynamically generate the fibonacci sequence.
var Fibonacci Enumerable[uint] = fibonacciGenerator{}

func (gen fibonacciGenerator) Enumerate(ctx context.Context) Enumerator[uint] {
	retval := make(chan uint)

	go func() {
		defer close(retval)
		var a, b uint = 0, 1

		for {
			select {
			case retval <- a:
				a, b = b, a+b
			case <-ctx.Done():
				return
			}
		}
	}()

	return retval
}
