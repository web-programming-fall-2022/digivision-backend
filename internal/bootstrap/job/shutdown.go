package job

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

type WithGracefulShutDown interface {
	ShutDown(ctx context.Context) error
}

func ShutDown(processes []WithGracefulShutDown, duration time.Duration) (shutdownError error) {
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()
	wg := &sync.WaitGroup{}
	for i := range processes {
		wg.Add(1)
		process := processes[i] // Beware of https://github.com/golang/go/wiki/CommonMistakes#using-goroutines-on-loop-iterator-variables
		go func() {
			if err := process.ShutDown(ctx); err != nil {
				curErr := fmt.Errorf("error in shutting down service: %v", err)
				if shutdownError != nil {
					logrus.Error(curErr.Error())
				} else {
					shutdownError = curErr
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()
	return
}
