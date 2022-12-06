package terminal

import (
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Terminal struct {
	dir string
}

func NewTerminal(dir string) *Terminal {
	return &Terminal{
		dir: dir,
	}
}

func (t *Terminal) Run(command string) (io.Reader, error) {
	args := strings.Split(command, " ")
	cmd := exec.Command(args[0], args[1:]...)

	if args[0] == "cd" {
		t.ChangeDir(args[1])
		return nil, nil
	} else {
		cmd.Dir = t.dir
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			return nil, err
		}
		if err := cmd.Start(); err != nil {
			return nil, err
		}
		return stdout, nil
	}
}

func (t *Terminal) ChangeDir(dir string) error {
	newPath := filepath.Join(t.dir, dir)

	if _, err := os.Stat(newPath); os.IsNotExist(err) {
		return err
	}
	t.dir = newPath

	return nil
}
