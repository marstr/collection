package collection

import "context"

type fibonacciGenerator struct{}

// Fibonacci is an Enumerable which will dynamically generate the fibonacci sequence.
var Fibonacci Enumerable = fibonacciGenerator{}

func (gen fibonacciGenerator) Enumerate(ctx context.Context) Enumerator {
	retval := make(chan interface{})

	go func() {
		defer close(retval)
		a, b := 0, 1

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
