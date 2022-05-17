package shutdown

import (
	"io"
	"log"
	"os"
	"os/signal"
)

func Graceful(signals []os.Signal, closers ...io.Closer) {
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, signals...)

	// wait for signal
	sig := <-sigc
	log.Printf("Caught signal %s. Shutting down...", sig)

	for _, closer := range closers {
		if err := closer.Close(); err != nil {
			log.Printf("failed to close %v: %v", closer, err)
		}
	}
}