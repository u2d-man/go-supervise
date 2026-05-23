package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestEnsureControlFIFO(t *testing.T) {
	dir := t.TempDir()
	path, err := ensureControlFIFO(dir)
	if err != nil {
		t.Errorf("not nil error.: %s", err)
	}

	if m, err := filepath.Match(path, dir+"/supervise/control"); m == false {
		t.Errorf("not <dir>/supervise/control %s", err)
	}

	info, err := os.Stat(path)
	if info.Mode()&os.ModeNamedPipe == 0 {
		t.Fatalf("not a FIFO mode=%v", info.Mode())
	}
}

func TestReadControlLoop_ReceivesByte(t *testing.T) {
	dir := t.TempDir()
	fifoPath, err := ensureControlFIFO(dir)
	if err != nil {
		t.Fatalf("")
	}

	cmdCh := make(chan byte, 16)
	quit := make(chan struct{})
	defer close(quit) // 後始末

	go readControlLoop(fifoPath, cmdCh, quit)

	// 1) FIFO に 'd' を書き込む
	//    ヒント: os.OpenFile(fifoPath, os.O_WRONLY, 0) → Write([]byte{'d'}) → Close
	//    ※ FIFO の Open は「もう片方が Open するまでブロック」する性質に注意

	// 2) cmdCh から受信。ただしタイムアウト付き（select + time.After(1*time.Second)）
	//    タイムアウトしたら t.Fatal("timeout: 'd' not received")
	//    'd' 以外が来たら t.Errorf
}
