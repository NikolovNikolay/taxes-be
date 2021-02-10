package util

import "sync"

type DelayedExecutor struct {
	wg      sync.WaitGroup
	delayed chan func()
}

func NewDelayedExecutor(parallel int) *DelayedExecutor {
	de := &DelayedExecutor{
		delayed: make(chan func(), 1000),
	}

	de.wg.Add(parallel)
	for i := 0; i < parallel; i++ {
		go func() {
			defer de.wg.Done()
			for f := range de.delayed {
				f()
			}
		}()
	}
	return de
}

func (de *DelayedExecutor) Stop() {
	close(de.delayed)
	de.wg.Wait()
}

func (de *DelayedExecutor) Run(f func()) {
	de.delayed <- f
}
