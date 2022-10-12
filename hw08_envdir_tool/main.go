package main

import (
	"os"
)

func main() {
	data, readDirError := ReadDir(os.Args[1])
	if readDirError != nil {
		println(readDirError)
	}

	commandError := RunCmd(os.Args[2:], data)
	if commandError != 0 {
		println(commandError)
	}
}
