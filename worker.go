package samidare

import "time"

type funcWorker struct {
	f func() int
}

func (f *funcWorker) Setup() {}

func (f *funcWorker) Process() int { return f.f() }

func (f *funcWorker) Teardown() {}

type Worker interface {
	Setup()
	Process() int
	Teardown()
}

type worker struct {
	worker    Worker
	seq       int
	readyCh   chan struct{}
	startCh   chan struct{}
	finishCh  chan struct{}
	resultCh  chan result
	stopCh    chan struct{}
	stopAllCh chan struct{}
	start     time.Time
	end       time.Time
	score     int
	count     int
}

type result struct {
	Seq     int
	Start   time.Time
	End     time.Time
	Elapsed time.Duration
	Score   int
	Count   int
}

func (w *worker) run() {
	w.worker.Setup()

	func() {
		w.readyCh <- struct{}{}
		<-w.startCh

		w.start = time.Now()
		for {
			select {
			case <-w.stopAllCh:
				return
			case <-w.stopCh:
				return
			default:
				w.score += w.worker.Process()
				w.count++
			}
		}
	}()

	w.end = time.Now()
	w.finishCh <- struct{}{}

	w.worker.Teardown()

	w.resultCh <- result{
		Seq:     w.seq,
		Start:   w.start,
		End:     w.end,
		Elapsed: w.end.Sub(w.start),
		Score:   w.score,
		Count:   w.count,
	}
}
