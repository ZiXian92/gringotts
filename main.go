package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

// Variables to be injected using compile-time injection
var version, env, opSys string

func main() {
	// Prepare to listen to OS signals
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	// Prepare a goroutine to run GUI client
	clientTermChan := make(chan error)
	clientKillChan := make(chan bool)
	go startClient(clientTermChan, clientKillChan)

	// Listen to OS signals
	go func() {
		<-sigs
		clientKillChan <- true
	}()

	fmt.Println(<-clientTermChan)
}

// Runs the GUI client
func startClient(termChan chan<- error, intChan <-chan bool) {
	var cmd *exec.Cmd
	var ext string
	done := make(chan error, 1)

	// Set up the command to run the client
	if env != "production" {
		os.Chdir("ui")
		cmd = exec.Command("npm", "run", "client")
	} else {
		if opSys == "win" {
			ext = ".exe"
		}
		cmdStr := fmt.Sprintf("gringotts-%s%s", version, ext)
		cmd = exec.Command(cmdStr)
	}

	// Start running the client
	done <- cmd.Run()

	select {
	// Main routine orders to stop client
	// Kill the client and return status to main routine
	case <-intChan:
		termChan <- cmd.Process.Kill()

	// Client terminated by itself
	// Report status to main routine
	case err := <-done:
		termChan <- err
	}
}
