package main

type fuzzer interface {
	SetDebugger(debugger) error
	SetCommand(string) error
	SetStatsD(string) error

	pick() error
	mutate() error
	execute() error
	monitor() error
	report() error
	clean() error

	Run()
}
