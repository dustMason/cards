package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"
)

func editFile(path string) error {
	if path == "" {
		return nil
	}
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = DefaultEditor
	}
	return run(editor, []string{path})
}

func run(cmd string, args []string) error {
	executable, execErr := exec.LookPath(cmd)
	if execErr != nil {
		return execErr
	}
	c := []string{cmd}
	c = append(c, args...)
	return syscall.Exec(executable, c, os.Environ())
}

func archive(dir string, filename string) error {
	return os.Rename(
		filepath.Join(dir, filename),
		filepath.Join(dir, ".archive", filename),
	)
}

func rename(dir string, old string, new string) error {
	return os.Rename(
		filepath.Join(dir, old),
		filepath.Join(dir, new),
	)
}

func pbcopy(file string) error {
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

func createFile(dir string) error {
	suffix := "md"
	path := filepath.Join(dir, fmt.Sprintf("%v.%v", time.Now().Format("2006-01-02-15-04-05"), suffix))
	return editFile(path)
}
