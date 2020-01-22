package main

import (
	"bytes"
	"os"
	"os/exec"
)

func Search(dir string) (string, error) {
	ag, execErr := exec.LookPath("ag")
	if execErr != nil {
		return "", execErr
	}
	var allContents bytes.Buffer
	agCommand := exec.Command(ag, "--nobreak", "--nonumbers", "--noheading", ".", dir)
	agCommand.Stdout = &allContents
	agCommand.Stderr = os.Stderr
	err := agCommand.Run()
	if err != nil {
		return "", err
	}

	var selected bytes.Buffer
	fzf, execErr := exec.LookPath("fzf")
	if execErr != nil {
		return "", execErr
	}

	fzfCommand := exec.Command(fzf)
	fzfCommand.Stdin = &allContents
	fzfCommand.Stdout = &selected
	fzfCommand.Stderr = os.Stderr
	err = fzfCommand.Run()
	if err != nil {
		return "", err
	}

	line, _ := selected.ReadString(':')
	fileToEdit := line[:len(line)-1]
	return fileToEdit, nil
}
