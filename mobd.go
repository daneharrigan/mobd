package main

import (
	"bufio"
	"github.com/kr/logfmt"
	"log"
	"os"
)

type stream struct {
	Forward bool
}

func (s *stream) HandleLogfmt(k, v []byte) error {
	switch string(k) {
	case "sample#memory_total":
		s.Forward = false
	case "sample#cpu_load_avg1m":
		s.Forward = false
	}

	return nil
}

func main() {
	log.SetPrefix("app=mobd ")
	log.SetFlags(0)
	log.Println("at=start")

	sc := bufio.NewScanner(os.Stdin)
	st := new(stream)

	for {
		sc.Scan()
		st.Forward = true

		if err := logfmt.Unmarshal(sc.Bytes(), st); err != nil {
			log.Printf("at=scan error=%q", err)
		}

		if st.Forward {
			os.Stdout.Write(append(sc.Bytes(), '\n'))
		}
	}

	log.Println("at=finish")
}
