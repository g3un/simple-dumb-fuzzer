package main

import (
	"os"
)

func main() {
	sdf := NewSdf(os.Args[1])
	lldb := NewLldb(os.Args[2])

	if err := sdf.SetDebugger(lldb); err != nil {
		panic(err)
	}
	statsdServer := os.Getenv("STATSD_SERVER")
	if statsdServer != "" {
		if err := sdf.SetStatsD(statsdServer); err != nil {
			panic(err)
		}
	}

	sdf.Run()
}
