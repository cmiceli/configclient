package main

import (
	"fmt"
	lib "github.com/cmiceli/configclient/lib"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	client := lib.NewHTTPClient("http://localhost:8000")
	fw := lib.NewFileWatcher(client)
	fw.AddFile("testing", "/tmp/testing", time.Second)

	done := make(chan bool, 1)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		done <- true
	}()
	fmt.Println("awaiting signal")
	<-done
	fmt.Println("exiting")
}
