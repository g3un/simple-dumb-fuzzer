package main

import (
	"os/exec"
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

func (w *Windbg) Run(cmd string) error {
	return nil
}

func (w Windbg) Terminate() error {
	return nil
}
