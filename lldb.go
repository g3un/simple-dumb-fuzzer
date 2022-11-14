package main

import (
	"os/exec"
	"strings"
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

func (l *Lldb) Run(cmd string) ([]byte, error) {
	c := strings.Split(cmd, " ")

	l.cmd = exec.Command(l.path, c...)

	return l.cmd.Output()
}
