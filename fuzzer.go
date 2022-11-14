package main

type fuzzer interface {
	SetDebugger(debugger) error
	SetStatsD(string) error

	pick() error
	mutate() error
	execute() (bool, error)
	report() error
	sendStatsD() error
	clear() error

	Run()
}
