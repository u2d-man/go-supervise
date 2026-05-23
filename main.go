package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"sync"
	"time"
)

type State struct {
	mu        sync.Mutex
	Pid       int
	StartedAt time.Time
	Restarts  int
	Status    string // "run" or "down"
}

func main() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	quit := make(chan struct{})
	go func() {
		sig := <-ch
		fmt.Println("received:", sig)

		close(quit)
	}()

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

	runPath := filepath.Join(absDir, "run")
	fmt.Println(runPath)

	fifoPath, err := ensureControlFIFO(absDir)
	if err != nil {
		fmt.Println(err)
		return
	}

	cmdCh := make(chan byte, 16)
	go readControlLoop(fifoPath, cmdCh, quit)

	state := State{}
	state.superviseLoop(runPath, dir, quit)
}

func ensureControlFIFO(serviceDir string) (path string, err error) {
	return serviceDir, nil
}

func readControlLoop(fifoPath string, ch chan<- byte, quit <-chan struct{}) error {
	return nil
}

func (s *State) superviseLoop(runPath string, dir string, quit <-chan struct{}) error {
	for {
		cmd := exec.Command(runPath)
		cmd.Dir = dir

		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Start()
		if err != nil {
			return err
		}

		s.start(cmd.Process.Pid, time.Now())
		fmt.Println(&s)

		done := make(chan error, 1)
		go func() {
			done <- cmd.Wait()
		}()

		select {
		case err := <-done:
			s.stop()
			if err != nil {
				var exitErr *exec.ExitError
				if errors.As(err, &exitErr) {
					fmt.Printf("exit code: %d\n", exitErr.ExitCode())
				}
			}
			fmt.Println("end child process")
		case <-quit:
			cmd.Wait()
			return nil
		}
	}
}

func (s *State) start(pid int, startedAt time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Pid = pid
	s.StartedAt = startedAt
	s.Restarts++
	s.Status = "run"
}

func (s *State) stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Status = "down"
}

func (s *State) String() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return fmt.Sprintf("pid=%d restarts=%d status=%s", s.Pid, s.Restarts, s.Status)
}
