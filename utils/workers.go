package utils

import (
	"sync"
)

func SafeClose[T any](ch chan T) {
	select {
	case <-ch:
		return
	default:
		close(ch)
	}
}

func SendPayloads[T any](payloads []T, inChan chan<- T) {
	go func() {
		for _, v := range payloads {
			inChan <- v
		}
		close(inChan)
	}()
}

func StartWorkers[T any, S any](count int, workerFn func(T, chan<- S)) (chan<- T, chan S) {
	inChan := make(chan T, count)
	outChan := make(chan S, count)

	wg := sync.WaitGroup{}
	for i := 0; i < count; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for val := range inChan {
				workerFn(val, outChan)
			}
		}()
	}

	go func() {
		wg.Wait()
		close(outChan)
	}()

	return inChan, outChan
}

func GetResults[T any](outChan <-chan T) []T {
	results := []T{}
	for v := range outChan {
		results = append(results, v)
	}

	return results
}

func GetSlicedResults[T any](outChan <-chan []T) []T {
	results := []T{}
	for v := range outChan {
		results = append(results, v...)
	}

	return results
}
