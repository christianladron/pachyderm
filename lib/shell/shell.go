package shell

import (
	"bytes"
	"io"
	"log"
	"os/exec"
	"strings"
)

func RunStderr(c *exec.Cmd) error {
	log.Print("RunStderr: ", strings.Join(c.Args, " "))
	stderr, err := c.StderrPipe()
	if err != nil {
		return err
	}
	err = c.Start()
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(stderr)
	if buf.Len() != 0 {
		log.Print("Command had output on stderr.\n Cmd: ", strings.Join(c.Args, " "), "\nstderr: ", buf)
	}
	return c.Wait()
}

func CallCont(c *exec.Cmd, cont func(io.Reader) error) error {
	log.Print("CallCont: ", strings.Join(c.Args, " "))
	reader, err := c.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := c.StderrPipe()
	if err != nil {
		return err
	}
	err = c.Start()
	if err != nil {
		return err
	}
	err = cont(reader)
	reader.Close()
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(stderr)
	if buf.Len() != 0 {
		log.Print("Command had output on stderr.\n Cmd: ", strings.Join(c.Args, " "), "\nstderr: ", buf)
	}

	return c.Wait()
}
