package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"
	"plex_monitor/internal/utils"
)

func getPassword(prompt string) string {
    fmt.Print(prompt)

    // Common settings and variables for both stty calls.
    attrs := syscall.ProcAttr{
        Dir:   "",
        Env:   []string{},
        Files: []uintptr{os.Stdin.Fd(), os.Stdout.Fd(), os.Stderr.Fd()},
        Sys:   nil}
    var ws syscall.WaitStatus

    // Disable echoing.
    pid, err := syscall.ForkExec(
        "/bin/stty",
        []string{"stty", "-echo"},
        &attrs)
    if err != nil {
        panic(err)
    }

    // Wait for the stty process to complete.
    _, err = syscall.Wait4(pid, &ws, 0, nil)
    if err != nil {
        panic(err)
    }

    // Echo is disabled, now grab the data.
    reader := bufio.NewReader(os.Stdin)
    text, err := reader.ReadString('\n')
    if err != nil {
        panic(err)
    }

    // Re-enable echo.
    pid, err = syscall.ForkExec(
        "/bin/stty",
        []string{"stty", "echo"},
        &attrs)
    if err != nil {
        panic(err)
    }

    // Wait for the stty process to complete.
    _, err = syscall.Wait4(pid, &ws, 0, nil)
    if err != nil {
        panic(err)
    }

    return strings.TrimSpace(text)
}

func main() {
	password := getPassword("[plex_monitor] Please enter a password: ")
	hashBytes, _ := utils.HashString(password)
	s := string(hashBytes)
	fmt.Printf("\n[plex_monitor] ====================Hashed password====================\n%s\n", s)
}
