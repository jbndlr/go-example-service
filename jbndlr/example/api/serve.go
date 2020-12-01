package api

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

// ServeGracefully : Run a service, listen for signals and terminate gracefully.
func ServeGracefully(start func() error, stop func() error) error {
	errc := make(chan error)
	sigc := make(chan os.Signal)
	signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		var err error

		defer func() {
			if rec := recover(); rec != nil {
				// If panic, try to recover
				if recerr, ok := rec.(error); ok {
					errc <- recerr
					return
				}
				errc <- fmt.Errorf("Recovered from: %v", rec)
				return
			}
			// Ohterwise pass on closure error
			errc <- err
			return
		}()

		err = start()
	}()

	select {
	case err := <-errc:
		return err
	case <-sigc:
		err := stop()
		return err
	}
}
