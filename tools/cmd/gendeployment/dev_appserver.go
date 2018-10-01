package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
)

func startLocalServer(appName string, deploymentDir string, main *Service, services []*Service) {
	var args []string
	args = append(args, fmt.Sprintf("--application=%s", appName))
	args = append(args, "--storage_path", ".localdata")
	args = append(args, "--datastore_consistency_policy", "consistent")
	if len(services) > 0 {
		args = append(args, filepath.Join(deploymentDir, main.Name, "dispatch.yaml"))
	}
	args = append(args, filepath.Join(deploymentDir, main.Name, "app.yaml"))
	for _, s := range services {
		args = append(args, filepath.Join(deploymentDir, s.Name, "app.yaml"))
	}
	cmd := exec.Command("dev_appserver.py", args...)
	log.Println(strings.Join(append([]string{"$ dev_appserver.py"}, args...), " "))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	log.Println("Starting dev_appserver with", cmd.Process.Pid)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	s := <-c
	log.Println("Killing", cmd.Process.Pid, "with", s)
	cmd.Process.Signal(s)
	cmd.Process.Wait()

}
