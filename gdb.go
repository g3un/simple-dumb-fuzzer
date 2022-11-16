package main

import (
	"os/exec"
	"strings"
)

type Gdb struct {
	path string
	cmd  *exec.Cmd
}

var _ debugger = (*Gdb)(nil)

func NewGdb(path string) *Gdb {
	return &Gdb{
		path: path,
		cmd:  nil,
	}
}

func (g *Gdb) Run(cmd string) ([]byte, error) {
	c := strings.Split(cmd, " ")

	g.cmd = exec.Command(g.path, c...)

	return g.cmd.Output()
}
