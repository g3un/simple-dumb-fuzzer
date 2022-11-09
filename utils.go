package main

import (
	"math/rand"
	"os"
	"regexp"
	"time"
)

func getFileName(path string) string {
	re := regexp.MustCompile(`^(.*/)?(?:$|(.+?)(?:(\.[^.]*$)|$))`)
	return re.FindStringSubmatch(path)[2]
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
