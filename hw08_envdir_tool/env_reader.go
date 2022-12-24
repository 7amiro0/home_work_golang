package main

import (
	"bufio"
	"os"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

func getFirstString(file *os.File) string {
	fileScanner := bufio.NewScanner(file)
	fileScanner.Split(bufio.ScanLines)

	fileScanner.Scan()
	firstLine := fileScanner.Text()

	editedLine := strings.TrimRight(firstLine, " 	")
	editedLine = strings.ReplaceAll(editedLine, "\x00", "\n")

	return editedLine
}

func ReadDir(dir string) (Environment, error) {
	allFiles, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	data := make(Environment)
	for _, file := range allFiles {
		if file.IsDir() || strings.Contains(file.Name(), "=") {
			continue
		}

		openFile, err := os.Open(dir + "/" + file.Name())
		if err != nil {
			return nil, err
		}

		needRemove := false

		firstLine := getFirstString(openFile)

		if info, _ := file.Info(); info.Size() == 0 {
			needRemove = true
		}

		data[file.Name()] = EnvValue{firstLine, needRemove}
	}

	return data, nil
}
