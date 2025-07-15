package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/apawelec/go-gnome-session-inhibit/inhibit"
)

func main() {
	i, err := inhibit.Acquire("Example app", "Test inhibit", inhibit.Idle|inhibit.Suspend)
	if err != nil {
		fmt.Printf("error occurred: %v\n", err)
		os.Exit(1)
	}

	<-exitSignal()

	if err := i.Release(); err != nil { // not strictly necessary at the end of program
		fmt.Printf("error occurred: %v\n", err)
		os.Exit(1)
	}
}

func exitSignal() chan os.Signal {
	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	return termChan
}
