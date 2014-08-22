package ui

import (
	"fmt"
	"io"
)

type UI interface {
	Sayln(string)
	Error(string)
}

type defaultUI struct {
	stdOut io.Writer
	stdErr io.Writer
}

func (dui *defaultUI) Sayln(message string) {
	dui.stdOut.Write([]byte(fmt.Sprintln(message)))
}

func (dui *defaultUI) Error(message string) {
	dui.stdErr.Write([]byte(fmt.Sprintln(message)))
}

func NewDefaultUI(stdOut, stdErr io.Writer) UI {
	return &defaultUI{
		stdOut: stdOut,
		stdErr: stdErr,
	}
}
