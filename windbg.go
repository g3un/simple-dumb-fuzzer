package main

import (
	"os/exec"
	"strings"
)

type Windbg struct {
	path string
	cmd  *exec.Cmd
}

var _ debugger = (*Windbg)(nil)

func NewWindbg(path string) *Windbg {
	return &Windbg{
		path: path,
		cmd:  nil,
	}
}

func (w *Windbg) Run(cmd string) ([]byte, error) {
	c := strings.Split(cmd, " ")

	w.cmd = exec.Command(w.path, c...)

	return w.cmd.Output()
}
