package main

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
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

	sample   string
	caseName string

	iterCount  int
	crashCount int

	// file i/o
	fin     *os.File
	finBuf  *bufio.Reader
	fout    *os.File
	foutBuf *bufio.Writer

	dbg        debugger
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
		statsDAddr: "",
	}
}

func (s *Sdf) SetDebugger(dbg debugger) error {
	s.dbg = dbg
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
	// Need to check on Windows
	fileName := getFileName(sample)
	s.caseName = fmt.Sprintf("case-%d-%s", s.iterCount, fileName)

	fout, err := os.Create(s.casePath + "/" + s.caseName)
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
		mutateSize := randomInt(0, s.finBuf.Size()-mutateOffset)

		s.foutBuf.Write(buf[:mutateOffset])
		s.foutBuf.Write(buf[mutateOffset+mutateSize:])
	}

	return nil
}

func (s Sdf) execute() (bool, error) {
	if s.dbg == nil {
		return false, fmt.Errorf("Please set a debugger before start fuzzing")
	}

	cmd := fmt.Sprintf("%s %s -o run", s.target, s.casePath+"/"+s.caseName)
	out, err := s.dbg.Run(cmd)
	if err != nil {
		return false, err
	}

	fout, err := os.Create(s.casePath + "/" + s.caseName + ".log")
	defer fout.Close()

	foutBuf := bufio.NewWriter(fout)
	defer foutBuf.Flush()

	if _, err := foutBuf.Write(out); err != nil {
		return false, err
	}

	return !bytes.Contains(out, []byte("status = 0 (0x00000000)")), nil
}

func (s Sdf) report() error {
	if err := copy(s.casePath+"/"+s.caseName, s.crashPath+"/"+s.caseName); err != nil {
		return err
	}
	if err := copy(s.casePath+"/"+s.caseName+".log", s.crashPath+"/"+s.caseName+".log"); err != nil {
		return err
	}

	return nil
}

func (s Sdf) clear() error {
	if err := os.RemoveAll(s.casePath); err != nil {
		return err
	}

	if err := mkdir(s.casePath); err != nil {
		return err
	}

	return nil
}

func (s Sdf) getStatsDData() string {
	tag := fmt.Sprintf("|#banner:%s,sdf_version:%s", s.target, "0.0.1")

	return strings.Join([]string{
		fmt.Sprintf("fuzzing.iter_count:%d|g%s", s.iterCount, tag),
		fmt.Sprintf("fuzzing.carsh_count:%d|g%s", s.crashCount, tag),
	}, "\n")
}

func (s Sdf) sendStatsD() error {
	data := s.getStatsDData()

	raddr, err := net.ResolveUDPAddr("udp", s.statsDAddr)
	if err != nil {
		return err
	}

	conn, err := net.DialUDP("udp", nil, raddr)
	if err != nil {
		return err
	}
	defer conn.Close()

	conn.Write([]byte(data))

	return nil
}

func (s Sdf) Run() {
	for {
		if err := s.pick(); err != nil {
			panic(err)
		}
		if err := s.mutate(); err != nil {
			panic(err)
		}

		isCrash, err := s.execute()
		if err != nil {
			panic(err)
		}

		s.iterCount += 1

		if isCrash {
			if err := s.report(); err != nil {
				panic(err)
			}
			s.crashCount += 1
		}

		if s.statsDAddr != "" {
			if err := s.sendStatsD(); err != nil {
				panic(err)
			}
		}

		if (s.iterCount % 10000) == 0 {
			if err := s.clear(); err != nil {
				panic(err)
			}
		}
	}
}
