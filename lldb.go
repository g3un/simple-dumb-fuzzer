package main

import (
	"os/exec"
)

type Lldb struct {
	path string
	cmd  *exec.Cmd
}

var _ debugger = (*Lldb)(nil)

func NewLldb(path string) *Lldb {
	return &Lldb{
		path: path,
		cmd:  nil,
	}
}

func (l *Lldb) Run(cmd string) error {
	return nil
}

func (l Lldb) Terminate() error {
	return nil
}
