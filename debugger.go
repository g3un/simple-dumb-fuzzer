package main

type debugger interface {
	Run(string) ([]byte, error)
}
