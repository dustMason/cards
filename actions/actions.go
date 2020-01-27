package actions

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"
)

const DefaultEditor = "vim"

func Create(dir string) error {
	suffix := "md"
	path := filepath.Join(dir, fmt.Sprintf("%v.%v", time.Now().Format("2006-01-02-15-04-05"), suffix))
	return Run(editor(), []string{"-c startinsert", path})
}

func Edit(path string) error {
	if path == "" {
		return nil
	}
	return Run(editor(), []string{path})
}

func editor() string {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = DefaultEditor
	}
	return editor
}

func Run(cmd string, args []string) error {
	executable, execErr := exec.LookPath(cmd)
	if execErr != nil {
		return execErr
	}
	c := []string{cmd}
	c = append(c, args...)
	return syscall.Exec(executable, c, os.Environ())
}

func Archive(dir string, filename string) error {
	_ = os.Mkdir(filepath.Join(dir, ".archive"), 0755)
	return os.Rename(
		filepath.Join(dir, filename),
		filepath.Join(dir, ".archive", filename),
	)
}

func Rename(dir string, old string, new string) error {
	return os.Rename(
		filepath.Join(dir, old),
		filepath.Join(dir, new),
	)
}

func Pbcopy(file string) error {
	pbcopy, err := exec.LookPath("pbcopy")
	if err != nil {
		return err
	}
	c := exec.Command(pbcopy)
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	c.Stdin = f
	return c.Run()
}
