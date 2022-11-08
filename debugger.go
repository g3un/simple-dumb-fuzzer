package main

type debugger interface {
	Run(string) error
	Terminate() error
}
