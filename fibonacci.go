package collection

type fibonacciGenerator struct{}

// Fibonacci is an Enumerable which will dynamically generate the fibonacci sequence.
var Fibonacci Enumerable[uint] = fibonacciGenerator{}

func (gen fibonacciGenerator) Enumerate(cancel <-chan struct{}) Enumerator[uint] {
	retval := make(chan uint)

	go func() {
		defer close(retval)
		var a, b uint = 0, 1

		for {
			select {
			case retval <- a:
				a, b = b, a+b
			case <-cancel:
				return
			}
		}
	}()

	return retval
}
