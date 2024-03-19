package utils

import (
	"os"
	"os/signal"
)

func WaitForSignals(signals ...os.Signal) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, signals...)
	<-ch
}
