package main

import (
	"flag"
	"fmt"
	"io"
	"log"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/mix3/samidare"
	"github.com/nsf/termbox-go"
)

func (c *client) drawLine(x, y int, str string) {
	color := termbox.ColorDefault
	backgroundColor := termbox.ColorDefault
	runes := []rune(str)
	for i := 0; i < len(runes); i += 1 {
		termbox.SetCell(x+i, y, runes[i], color, backgroundColor)
	}
}

func (c *client) draw() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	c.drawLine(0, 0, "Press Ctrl-C to exit.")
	c.drawLine(0, 1, fmt.Sprintf("DesiredWorkerNum: %d", c.desired))
	c.drawLine(0, 2, fmt.Sprintf("CurrentWorkerNum: %d", c.current))
	termbox.Flush()
}

func (c *client) pollEvent() {
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyCtrlC:
				return
			case termbox.KeyArrowRight:
				err := c.stream.Send(&samidare.WatchRequest{
					Type: samidare.WatchRequest_INCREMENT,
				})
				if err != nil {
					log.Println(err)
					return
				}
			case termbox.KeyArrowLeft:
				err := c.stream.Send(&samidare.WatchRequest{
					Type: samidare.WatchRequest_DECREMENT,
				})
				if err != nil {
					log.Println(err)
					return
				}
			}
		}
	}
}

func (c *client) watch() error {
	for {
		res, err := c.stream.Recv()
		if err == io.EOF {
			continue
		}
		if err != nil {
			c.stream = nil
			termbox.Close()
			return err
		}
		switch res.Type {
		case samidare.WatchResponse_DESIRED:
			c.desired = res.Num
		case samidare.WatchResponse_CURRENT:
			c.current = res.Num
		}
		c.draw()
	}
}

func (c *client) join(conn *grpc.ClientConn) error {
	cli := samidare.NewSamidareClient(conn)
	stream, err := cli.Watch(context.Background())
	if err != nil {
		return err
	}

	err = stream.Send(&samidare.WatchRequest{})
	if err != nil {
		return err
	}

	_, err = stream.Recv()
	if err != nil {
		return err
	}

	c.stream = stream
	return nil
}

type client struct {
	stream  samidare.Samidare_WatchClient
	desired int64
	current int64
}

func (c *client) run() error {
	if err := termbox.Init(); err != nil {
		return err
	}
	defer termbox.Close()
	go c.watch()
	c.pollEvent()
	return nil
}

func main() {
	var addr string
	c := &client{}
	flag.StringVar(&addr, "addr", "localhost:19300", "addr")
	flag.Parse()

	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	if err := c.join(conn); err != nil {
		log.Fatal(err)
	}

	if err := c.run(); err != nil {
		log.Fatal(err)
	}
}
