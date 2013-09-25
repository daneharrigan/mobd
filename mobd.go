package main

import (
	"bufio"
	"errors"
	"github.com/kr/logfmt"
	"log"
	"os"
)

var (
	payload []byte
	size    chan bool
	fwd     chan bool
	max     int
)

func init() {
	size = make(chan bool)
	fwd = make(chan bool)
	max = 1000
}

type stream struct{}

func (s *stream) HandleLogfmt(k, v []byte) error {
	return errors.New("no metrics found")
}

func main() {
	log.SetPrefix("app=mobd ")
	log.SetFlags(0)

	go forward()
	go verify()

	sc := bufio.NewScanner(os.Stdin)
	st := new(stream)

	for {
		log.Println("at=scan")
		sc.Scan()
		if err := logfmt.Unmarshal(sc.Bytes(), st); err != nil {
			payload = append(payload, sc.Bytes()...)
			size <- true
		}
	}
}

func verify() {
	for {
		select {
		case <-size:
			if len(payload) >= max {
				fwd <- true
			}
		}
	}
}

func forward() {
	for _ = range fwd {
		log.Printf("%d", len(payload))
		payload = payload[:0]
	}
}
