package utils

import (
	"sync"
)

func ConcurrentMap[I any, O any](f func(I) (O, error), inputs []I) ([]O, error) {
	outputs := make([]O, len(inputs))
	var wg sync.WaitGroup
	var lock sync.Mutex
	wg.Add(len(inputs))
	outErr := make(chan error, len(inputs))
	for i, input := range inputs {
		go func(i int, input I) {
			defer wg.Done()
			output, err := f(input)
			if err != nil {
				outErr <- err
				return
			}
			lock.Lock()
			outputs[i] = output
			lock.Unlock()
		}(i, input)
	}
	wg.Wait()
	if len(outErr) > 0 {
		return nil, <-outErr
	}
	close(outErr)
	return outputs, nil
}
