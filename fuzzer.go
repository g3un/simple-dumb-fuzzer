package main

type fuzzer interface {
	SetDebugger(debugger) error
	SetCommand(string) error
	SetStatsD(string) error

	pick() error
	mutate() error
	execute() error
	monitor(chan bool) error
	report() error
	sendStatsD() error
	clear() error

	Run()
}
