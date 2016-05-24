package samidare

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"google.golang.org/grpc"
)

var TrapSignals = []os.Signal{
	syscall.SIGHUP,
	syscall.SIGINT,
	syscall.SIGTERM,
	syscall.SIGQUIT,
}

type server struct {
	spawner          func() Worker
	desiredWorkerNum int64
	currentWorkerNum int64
	counter          int
	workers          []*worker
	results          []result

	stream Samidare_WatchServer
	mu     *sync.RWMutex

	desiredCh chan int
	currentCh chan int
	watchCh   chan *WatchResponse
	startCh   chan struct{}
	resultCh  chan result
	stopAllCh chan struct{}
	wg        *sync.WaitGroup
}

func (s *server) getWatcher() Samidare_WatchServer {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.stream
}

func (s *server) setWatcher(stream Samidare_WatchServer) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.stream != nil {
		return fmt.Errorf("server is busy")
	}
	s.stream = stream
	return nil
}

func (s *server) delWatcher() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.stream = nil
}

func (s *server) Watch(stream Samidare_WatchServer) error {
	if _, err := stream.Recv(); err != nil {
		log.Println(err)
		stream.Send(nil)
		return err
	}

	if err := s.setWatcher(stream); err != nil {
		log.Println(err)
		stream.Send(nil)
		return err
	}

	if err := stream.Send(&WatchResponse{}); err != nil {
		log.Println(err)
		return err
	}

	for _, res := range []*WatchResponse{
		{Type: WatchResponse_DESIRED, Num: s.desiredWorkerNum},
		{Type: WatchResponse_CURRENT, Num: s.currentWorkerNum},
	} {
		s.watchCh <- res
	}

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			continue
		}
		if err != nil {
			log.Println(err)
			s.delWatcher()
			return err
		}

		switch req.Type {
		case WatchRequest_INCREMENT:
			s.desiredCh <- DESIRED_WORKER_INCR
		case WatchRequest_DECREMENT:
			s.desiredCh <- DESIRED_WORKER_DECR
		}
	}

	return nil
}

func (s *server) watch() {
	for res := range s.watchCh {
		if s.stream == nil {
			continue
		}

		err := s.stream.Send(res)
		if err != nil {
			log.Println(err)
		}
	}
}

const (
	DESIRED_WORKER_INCR int = iota
	DESIRED_WORKER_DECR
)

func (s *server) desired() {
	for ch := range s.desiredCh {
		switch ch {
		case DESIRED_WORKER_INCR:
			w := s.newWorker()
			s.desiredWorkerNum++
			s.workers = append(s.workers, w)
			s.runWorker(w)
		case DESIRED_WORKER_DECR:
			if 0 < s.desiredWorkerNum {
				w := s.workers[0]
				s.desiredWorkerNum--
				s.workers = s.workers[1:]
				s.stopWorker(w)
			}
		}
		s.watchCh <- &WatchResponse{
			Type: WatchResponse_DESIRED,
			Num:  s.desiredWorkerNum,
		}
	}
}

const (
	CURRENT_WORKER_INCR int = iota
	CURRENT_WORKER_DECR
)

func (s *server) current() {
	for ch := range s.currentCh {
		switch ch {
		case CURRENT_WORKER_INCR:
			s.currentWorkerNum++
		case CURRENT_WORKER_DECR:
			s.currentWorkerNum--
		}
		s.watchCh <- &WatchResponse{
			Type: WatchResponse_CURRENT,
			Num:  s.currentWorkerNum,
		}
	}
}

func (s *server) stopWorker(w *worker) {
	go func() {
		w.stopCh <- struct{}{}
	}()
}

func (s *server) runWorker(w *worker) {
	s.wg.Add(1)
	go w.run()
	go func() {
		<-w.readyCh
		s.currentCh <- CURRENT_WORKER_INCR
		<-w.finishCh
		s.currentCh <- CURRENT_WORKER_DECR
		s.resultCh <- <-w.resultCh
	}()
}

func (s *server) result() {
	for ch := range s.resultCh {
		s.results = append(s.results, ch)
		s.wg.Done()
	}
}

func (s *server) newWorker() *worker {
	s.counter++
	return &worker{
		seq:       s.counter,
		worker:    s.spawner(),
		readyCh:   make(chan struct{}),
		startCh:   s.startCh,
		finishCh:  make(chan struct{}),
		resultCh:  make(chan result),
		stopCh:    make(chan struct{}),
		stopAllCh: s.stopAllCh,
	}
}

func (s *server) server(addr string) {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	server := grpc.NewServer()
	RegisterSamidareServer(server, s)
	go server.Serve(l)
	log.Printf("listening on %s", addr)

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, TrapSignals...)
	sig := <-signalCh
	switch sig.(type) {
	case syscall.Signal:
		log.Printf("Got signa: %s(%d)", sig, sig)
	default:
		log.Printf("interrupred %s", sig)
	}

	server.Stop()

	close(s.stopAllCh)
	s.wg.Wait()
}

func Run(spawner func() Worker) []result {
	var (
		addr string
		num  int
	)
	flag.StringVar(&addr, "addr", ":19300", "addr")
	flag.IntVar(&num, "c", 1, "desired worker num")
	flag.Parse()

	s := &server{
		spawner:   spawner,
		mu:        &sync.RWMutex{},
		desiredCh: make(chan int),
		currentCh: make(chan int),
		watchCh:   make(chan *WatchResponse),
		startCh:   make(chan struct{}),
		resultCh:  make(chan result),
		stopAllCh: make(chan struct{}),
		wg:        &sync.WaitGroup{},
	}

	go s.desired()
	go s.current()
	go s.watch()
	go s.result()

	for i := 0; i < num; i++ {
		w := s.newWorker()
		s.desiredWorkerNum++
		s.workers = append(s.workers, w)
		s.wg.Add(1)
		go w.run()
		go func() {
			s.currentCh <- CURRENT_WORKER_INCR
			<-w.finishCh
			s.currentCh <- CURRENT_WORKER_DECR
			s.resultCh <- <-w.resultCh
		}()
		<-w.readyCh
	}
	close(s.startCh)

	s.server(addr)

	log.Println("shutdown")

	return s.results
}

func RunFunc(f func() int) []result {
	return Run(func() Worker { return &funcWorker{f: f} })
}

func Fprint(w io.Writer, res []result) {
	var (
		sum float64
		num int
	)
	for _, r := range res {
		sum += float64(r.Score) / r.Elapsed.Seconds()
		num++
	}
	avg := sum / float64(num)
	fmt.Fprintf(w, "%f score / sec\n", avg)
}

func Print(res []result) {
	Fprint(os.Stdout, res)
}
