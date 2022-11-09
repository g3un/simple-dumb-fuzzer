package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Sdf struct {
	target string

	basePath  string
	seedPath  string
	casePath  string
	crashPath string

	sample string

	iterCount  int
	crashCount int

	// file i/o
	fin     *os.File
	finBuf  *bufio.Reader
	fout    *os.File
	foutBuf *bufio.Writer

	dbg        *debugger
	cmd        string
	statsDAddr string
}

var _ fuzzer = (*Sdf)(nil)

func NewSdf(target string) *Sdf {
	_, err := os.Stat(target)
	if os.IsNotExist(err) {
		panic(err)
	}

	basePath, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	if err = mkdir("seeds"); err != nil {
		panic(err)
	}
	if err = mkdir("cases"); err != nil {
		panic(err)
	}
	if err = mkdir("crashes"); err != nil {
		panic(err)
	}

	return &Sdf{
		target: target,

		basePath:  basePath,
		seedPath:  basePath + "/seeds",
		casePath:  basePath + "/cases",
		crashPath: basePath + "/crashes",

		sample: "",

		iterCount:  0,
		crashCount: 0,

		dbg:        nil,
		cmd:        "",
		statsDAddr: "",
	}
}

func (s *Sdf) SetDebugger(dbg debugger) error {
	s.dbg = &dbg
	return nil
}

func (s *Sdf) SetCommand(cmd string) error {
	s.cmd = cmd
	return nil
}

func (s *Sdf) SetStatsD(addr string) error {
	s.statsDAddr = addr
	return nil
}

func (s *Sdf) pick() error {
	matches, err := filepath.Glob(s.seedPath + "/*")
	if err != nil {
		return err
	}
	if len(matches) == 0 {
		return fmt.Errorf("Need seed file")
	}

	sample := randomChoice(matches)
	s.sample = sample

	fin, err := os.Open(sample)
	if err != nil {
		return err
	}

	// TODO
	// need to check
	fileName := getFileName(sample)
	caseName := fmt.Sprintf("case-%d-%s", s.iterCount, fileName)

	fout, err := os.Create(s.casePath + caseName)
	if err != nil {
		return err
	}

	s.fin = fin
	s.finBuf = bufio.NewReader(fin)
	s.fout = fout
	s.foutBuf = bufio.NewWriter(fout)

	return nil
}

func (s Sdf) mutate() error {
	defer s.fin.Close()
	defer s.fout.Close()
	defer s.foutBuf.Flush()

	buf := make([]byte, s.finBuf.Size())
	if _, err := s.finBuf.Read(buf); err != nil {
		return err
	}

	switch randomInt(0, 2) {
	// Insert
	case 0:
		mutateOffset := randomInt(1, s.finBuf.Size())
		mutateSize := randomInt(1, 1000)
		mutateString := randomChoice([]string{"\x00", " ", "A", "1", "%s"})

		s.foutBuf.Write(buf[:mutateOffset])
		s.foutBuf.WriteString(strings.Repeat(mutateString, mutateSize))
		s.foutBuf.Write(buf[mutateOffset:])

	// Delete
	case 1:
		mutateOffset := randomInt(1, s.finBuf.Size())
		mutateSize := randomInt(1, 1000)

		s.foutBuf.Write(buf[:mutateOffset])
		s.foutBuf.Write(buf[mutateOffset+mutateSize:])
	}

	return nil
}

func (s Sdf) execute() error {
	// TODO
	return nil
}

func (s Sdf) monitor(ch chan bool) error {
	// TODO
	ch <- false
	return nil
}

func (s Sdf) report() error {
	// TODO
	return nil
}

func (s Sdf) clear() error {
	// TODO
	return nil
}

func (s Sdf) sendStatsD() error {
	// TODO
	return nil
}

func (s Sdf) Run() {
	for {
		ch := make(chan bool, 1)
		defer close(ch)

		if err := s.pick(); err != nil {
			panic(err)
		}
		if err := s.mutate(); err != nil {
			panic(err)
		}

		go s.execute()
		go s.monitor(ch)

		s.iterCount += 1

		if <-ch {
			s.report()
			s.crashCount += 1
		}

		if s.statsDAddr != "" {
			if err := s.sendStatsD(); err != nil {
				panic(err)
			}
		}

		if (s.iterCount % 10000) == 0 {
			s.clear()
		}
	}
}
