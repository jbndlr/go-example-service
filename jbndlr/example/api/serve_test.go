package api_test

import (
	"fmt"
	"jbndlr/example/api"
	"os"
	"syscall"
	"testing"
	"time"
)

func TestServeGracefullyPanic(t *testing.T) {
	start := func() error { panic("don't panic") }
	stop := func() error { return nil }
	err := api.ServeGracefully(start, stop)

	if err == nil || err.Error() != "Recovered from: don't panic" {
		t.Errorf("Panic recovery failed (%v)", err)
	}
}

func TestServeGracefullyPanicError(t *testing.T) {
	start := func() error { panic(fmt.Errorf("don't panic")) }
	stop := func() error { return nil }
	err := api.ServeGracefully(start, stop)

	if err == nil || err.Error() != "don't panic" {
		t.Errorf("Error not properly unpacked from panic: %v", err)
	}
}

func TestServeGracefullyError(t *testing.T) {
	start := func() error { return fmt.Errorf("failure") }
	stop := func() error { return nil }
	err := api.ServeGracefully(start, stop)

	if err == nil || err.Error() != "failure" {
		t.Errorf("Received unexpected error")
	}
}

func TestServeGracefullySignal(t *testing.T) {
	proc, errp := os.FindProcess(os.Getpid())
	if errp != nil {
		t.Fatal(errp)
	}

	start := func() error { time.Sleep(1 * time.Second); return nil }
	stop := func() error { return fmt.Errorf("stopped") }

	errc := make(chan error)
	go func() {
		err := api.ServeGracefully(start, stop)
		errc <- err
	}()
	proc.Signal(syscall.SIGTERM)

	//t.Errorf("Error: %v", <-errc)
}
