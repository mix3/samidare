package samidare

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWorker(t *testing.T) {
	assert := assert.New(t)

	startCh := make(chan struct{})
	stopAllCh := make(chan struct{})
	workers := []*worker{
		&worker{
			worker:    &funcWorker{f: func() int { return 1 }},
			readyCh:   make(chan struct{}),
			startCh:   startCh,
			finishCh:  make(chan struct{}),
			resultCh:  make(chan result),
			stopCh:    make(chan struct{}),
			stopAllCh: stopAllCh,
		},
		&worker{
			worker:    &funcWorker{f: func() int { return 1 }},
			readyCh:   make(chan struct{}),
			startCh:   startCh,
			finishCh:  make(chan struct{}),
			resultCh:  make(chan result),
			stopCh:    make(chan struct{}),
			stopAllCh: stopAllCh,
		},
		&worker{
			worker:    &funcWorker{f: func() int { return 1 }},
			readyCh:   make(chan struct{}),
			startCh:   startCh,
			finishCh:  make(chan struct{}),
			resultCh:  make(chan result),
			stopCh:    make(chan struct{}),
			stopAllCh: stopAllCh,
		},
	}
	for _, w := range workers {
		go w.run()
		<-w.readyCh
	}
	close(startCh)

	time.Sleep(1 * time.Second)

	w := workers[0]
	w.stopCh <- struct{}{}
	<-w.finishCh
	assert.IsType(result{}, <-w.resultCh)

	workers = workers[1:]

	close(stopAllCh)
	for _, w := range workers {
		<-w.finishCh
		assert.IsType(result{}, <-w.resultCh)
	}
}
