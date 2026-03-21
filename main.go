package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

func main() {
	var dir string
	for i, v := range os.Args {
		if i == 1 {
			dir = v
		}
	}

	absDir, err := filepath.Abs(dir)
	if err != nil {
		fmt.Println(err)
		return
	}

	runPath := filepath.Join(absDir, "run.bat")
	fmt.Println(runPath)

	for {
		cmd := exec.Command(runPath)
		cmd.Dir = dir

		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Start()
		if err != nil {
			fmt.Println(err)
			return
		}

		err = cmd.Wait()
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("end child process")

		time.Sleep(1 * time.Second)

	}

}
