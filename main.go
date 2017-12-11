package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"
)

// Variables to be injected using compile-time injection
var version, env, opSys, host, port string

func main() {
	// Prepare to listen to OS signals
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	// Prepare a goroutine to run GUI client
	log.Println("Starting UI client")
	clientTermChan := make(chan error, 1)
	clientKillChan := make(chan bool)
	go startClient(clientTermChan, clientKillChan)

	// Set up API server for UI client to call
	log.Println("Starting core server")
	svrTermChan := make(chan error, 1)
	svr := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", host, port),
		Handler: initCoreRouter(),
	}
	go startServer(svr, svrTermChan)

	select {
	// Received user request to terminate
	case <-sigs:

	// UI client failed
	case err := <-clientTermChan:
		log.Printf("Client stopped with error: %v", err)
		close(clientKillChan)
		clientKillChan = nil

	// Core API server failed
	case err := <-svrTermChan:
		log.Printf("Core server ended with error: %v", err)
		svr = nil
	}

	// Cleanup code
	// Wrap up UI client if it did not fail
	if clientKillChan != nil {
		log.Printf("Stopping UI client")
		clientKillChan <- true
		<-clientTermChan
	}

	// Stop the core API server
	if svr != nil {
		log.Println("Stopping core server")
		svr.Shutdown(nil)
		<-svrTermChan
	}
}

// Sets up and returns a router to expose core functionalities to
// the UI client
func initCoreRouter() http.Handler {
	r := mux.NewRouter()

	// TODO: Bind the routes

	return r
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
	// Run in another goroutine
	go func() {
		done <- cmd.Run()
	}()

	// Listen for any interruptions
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

	// Close the channel so the main routine will not block
	// when waiting for everyone to finish
	close(termChan)
}

// Starts the given HTTP server
// This function is meant to be used in a separate goroutine
func startServer(svr *http.Server, errChan chan<- error) {
	if svr != nil {
		errChan <- svr.ListenAndServe()
	}
	close(errChan)
}
