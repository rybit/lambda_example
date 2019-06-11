package main

import (
	"bufio"
	"os"
	"os/signal"
	"syscall"
)

func readAndSend() {
	var done bool
	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
		<-signals
		done = true
		rootLogger.Info("Shutdown, hit enter if it is stuck")
	}()
	rootLogger.Info("Starting to listen for lines from the cmdline")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() && !done {
		out := &humioMsg{
			Messages: []string{scanner.Text()},
		}
		send(rootLogger, out)
	}
	rootLogger.Info("Shutdown")
	os.Exit(0)
}
