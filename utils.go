package main

import (
	"bufio"
	"math/rand"
	"os"
	"regexp"
	"time"
)

func getFileName(path string) string {
	re := regexp.MustCompile(`^(.*/)?(?:$|(.+?)(?:(\.[^.]*$)|$))`)
	return re.FindStringSubmatch(path)[2]
}

func copy(src string, dest string) error {
	fin, err := os.Open(src)
	if err != nil {
		return err
	}
	defer fin.Close()
	finBuf := bufio.NewReader(fin)

	fout, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer fout.Close()

	foutBuf := bufio.NewWriter(fout)
	defer foutBuf.Flush()

	finStat, err := fin.Stat()
	if err != nil {
		return err
	}

	buf := make([]byte, finStat.Size())

	_, err = finBuf.Read(buf)
	if err != nil {
		return err
	}

	foutBuf.Write(buf)

	return nil
}

// [start, end)
func randomInt(start, end int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(end-start) + start
}

func randomChoice(list []string) string {
	rand.Seed(time.Now().Unix())
	return list[rand.Intn(len(list))]
}

func mkdir(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return os.Mkdir(dir, os.ModePerm)
	}
	return nil
}
